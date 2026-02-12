package main

import (
	"ai-service/application"
	"ai-service/infrastructure/config"
	"ai-service/infrastructure/financialdata"
	"ai-service/infrastructure/gemini"
	kafkaclient "ai-service/infrastructure/kafka"
	authmw "ai-service/infrastructure/middleware"
	"ai-service/infrastructure/parser"
	"ai-service/infrastructure/s3"
	"context"
	"log"
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

	geminiClient, err := gemini.NewClient(cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	s3Client, err := s3.NewClient(cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3BucketName, cfg.S3Endpoint)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	kafkaClient := kafkaclient.NewKafkaClient(cfg.KafkaURL, cfg.KafkaTopic)

	parserClient := parser.NewClient(cfg.ParserURL)
	fdClient := financialdata.NewClient(cfg.FinancialDataURL, cfg.FinancialDataAPIKey)

	extractorService := application.NewExtractorService(geminiClient, s3Client, parserClient, fdClient)
	extractorHandler := application.NewExtractorHandler(extractorService)
	taskProcessor := application.NewTaskProcessor(10, kafkaClient)
	taskProcessor.Start(context.Background())

	if cfg.APIKey == "" {
		log.Fatal("AI_SERVICE_API_KEY is required")
	}

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
	})

	addr := ":" + cfg.Port
	log.Printf("AI Service starting on %s", addr)

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
		log.Fatalf("Failed to start server: %v", err)
	case sig := <-shutdown:
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		taskProcessor.Stop(ctx)
		if err := kafkaClient.Close(); err != nil {
			log.Printf("Failed to close Kafka client: %v", err)
		}

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown server: %v", err)
		}
		log.Println("Server stopped gracefully")
	}
}
