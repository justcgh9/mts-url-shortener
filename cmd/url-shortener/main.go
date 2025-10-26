package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/justcgh9/mts-url-shortener/internal/config"
	"github.com/justcgh9/mts-url-shortener/internal/db/postgres"
	"github.com/justcgh9/mts-url-shortener/internal/http/handlers"
	"github.com/justcgh9/mts-url-shortener/internal/logger"
	"github.com/justcgh9/mts-url-shortener/internal/service/url"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	_ = log

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Http.Timeout)
	defer cancel()

	storage := postgres.MustConnect(ctx, cfg.StoragePath)

	svc := url.NewService(log, storage)

	r := chi.NewRouter()

	r.Get("/{alias}", handlers.NewRedirectHandler(svc))
	r.Post("/", handlers.NewCreateHandler(svc))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Http.Port),
		Handler:      r,
		ReadTimeout:  cfg.Http.Timeout,
		WriteTimeout: cfg.Http.Timeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	log.Info("starting up server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen error", slog.String("err", err.Error()))
		}
	}()

	<-done
	log.Info("stopping the server")
}
