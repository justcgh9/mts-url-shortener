package main

import (
	"github.com/justcgh9/mts-url-shortener/internal/config"
	"github.com/justcgh9/mts-url-shortener/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Info("app started")
}