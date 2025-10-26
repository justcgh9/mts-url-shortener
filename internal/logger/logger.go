package logger

import (
	"log"
	"log/slog"
	"os"
)

const (
	localEnv = "local"
	testEnv = "testing"
	prodEnv = "prod"
)

func New(env string) *slog.Logger {
	var handler slog.Handler
	switch env {
	case localEnv:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case testEnv:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	case prodEnv:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		log.Fatalf("config failed to provide valid env")
	}

	return slog.New(handler)
}