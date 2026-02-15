package main

import (
	"ai-service/internal/application"
	"ai-service/internal/handlers"
	"ai-service/internal/infrastructure"
	"ai-service/internal/infrastructure/config"
	"ai-service/internal/infrastructure/financialdata"
	"ai-service/internal/infrastructure/gemini"
	kafkaclient "ai-service/internal/infrastructure/kafka"
	authmw "ai-service/internal/infrastructure/middleware"
	"ai-service/internal/infrastructure/parser"
	"ai-service/internal/infrastructure/postgres"
	"ai-service/internal/infrastructure/s3"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()
	err := cfg.Validate()
	if err != nil {
		slog.Error("Config is not valid", slog.Any("error", err))
		os.Exit(1)
	}

	if err := infrastructure.RunMigrations(); err != nil {
		slog.Error("Failed to run migrations", slog.Any("error", err))
		os.Exit(1)
	}

	ctx := context.Background()

	geminiClient, err := gemini.NewClient(cfg.GeminiAPIKey, cfg.GeminiProxyURL)
	if err != nil {
		slog.Error("Failed to create Gemini client", slog.Any("error", err))
		os.Exit(1)
	}

	s3Client, err := s3.NewClient(cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3BucketName, cfg.S3Endpoint)
	if err != nil {
		slog.Error("Failed to create S3 client", slog.Any("error", err))
		os.Exit(1)
	}

	db, err := postgres.NewDBRepo(ctx, cfg.PostgresURL)
	if err != nil {
		slog.Error("Failed to create DB Repo", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	kafkaClient := kafkaclient.NewKafkaClient(cfg.KafkaURL, cfg.KafkaTopic)

	parserClient := parser.NewClient(cfg.ParserURL)
	fdClient := financialdata.NewClient(cfg.FinancialDataURL, cfg.FinancialDataAPIKey)

	extractorService := application.NewExtractorService(geminiClient, s3Client, parserClient, fdClient)
	geminiService := application.NewGeminiService(geminiClient, s3Client, fdClient)
	extractorHandler := handlers.NewExtractorHandler(extractorService)
	analysisHandler := handlers.NewAnalysisHandler(db)
	taskProcessor := application.NewTaskProcessor(10, geminiService, kafkaClient, db)
	taskProcessor.Start(context.Background())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Group(func(r chi.Router) {
		r.Use(authmw.APIKeyAuth(cfg.APIKey))
		r.Get("/extract", extractorHandler.HandleExtract)
		r.Get("/analysis", analysisHandler.HandleGetAnalysis)
		r.Get("/analyses", analysisHandler.HandleGetAnalysesByTicker)
	})

	addr := ":" + cfg.Port
	slog.Info("AI Service starting", slog.Any("addr", addr))

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	serverErrors := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		slog.Error("Failed to start server", slog.Any("error", err))
	case sig := <-shutdown:
		slog.Info("Received signal, shutting down gracefully...", slog.Any("signal", sig))
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		taskProcessor.Stop(ctx)
		if err := kafkaClient.Close(); err != nil {
			slog.Error("Failed to close Kafka client", slog.Any("error", err))
		}

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown server", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("Server stopped gracefully")
	}
}
