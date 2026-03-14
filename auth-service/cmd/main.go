package main

import (
	"log"
	"log/slog"
	"os"

	internal "auth-service/internal/application"
	"auth-service/internal/config"
	"auth-service/internal/infrastructure"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	slog.Info("Configuration loaded successfully")

	if err := infrastructure.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	service, err := internal.NewService(cfg)
	handlers := internal.NewHandlers(service, cfg)
	server := internal.NewServer(handlers, cfg)
	server.MountRoutes()

	if err := server.Start(cfg.Server.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
