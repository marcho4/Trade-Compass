package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	internal "auth-service/internal/application"
	"auth-service/internal/config"
	"auth-service/internal/domain"
	"auth-service/internal/infrastructure"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Configuration loaded successfully")

	if err := infrastructure.RunMigrations(logger); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connection established")

	repo := infrastructure.NewDbRepo(pool)

	jwtConfig := domain.JWTConfig{
		SecretKey:       []byte(cfg.JWT.Secret),
		AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
	}

	jwtService, err := infrastructure.NewJWTService(jwtConfig)
	if err != nil {
		log.Fatalf("Failed to create JWT service: %v", err)
	}

	log.Println("JWT service initialized")

	service := internal.NewService(repo, jwtService, cfg.OAuth)
	handlers := internal.NewHandlers(service, cfg)
	server := internal.NewServer(handlers, cfg)
	server.MountRoutes()

	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := server.Start(cfg.Server.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
