package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ai-service/internal/app"
	"ai-service/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		slog.Error("Invalid config", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	application, err := app.New(ctx, cfg)
	if err != nil {
		slog.Error("Failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}
	defer application.Shutdown()

	if err := application.Run(ctx); err != nil {
		slog.Error("App stopped with error", slog.Any("error", err))
		os.Exit(1)
	}
}
