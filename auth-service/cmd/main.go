package main

import (
	"log"
	"log/slog"
	"os"

	internal "auth-service/internal/application"
	"auth-service/internal/infrastructure/migrations"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := runMigrations(logger); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	server := internal.CreateServer()
	server.MountRoutes()
	server.Start()
}

func runMigrations(logger *slog.Logger) error {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	migrator, err := migrations.NewMigrator(dbURL, "file://migrations", logger)
	if err != nil {
		return err
	}
	defer migrator.Close()

	return migrator.Up()
}
