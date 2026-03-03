package main

import (
	"financial_data/internal/application"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	finData, err := application.NewFinData()
	if err != nil {
		slog.Error("Failed to create fin data service", "error", err)
		os.Exit(1)
	}

	if err := finData.Run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}
