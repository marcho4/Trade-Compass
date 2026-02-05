package main

import (
	"context"
	"financial_data/internal/application"
	"financial_data/internal/infrastructure"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	if err := infrastructure.RunMigrations(); err != nil {
		return err
	}

	pool, err := infrastructure.NewPostgresPool(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	slog.Info("Database connection pool initialized")

	ratiosRepo := infrastructure.NewRatiosRepository(pool)
	rawDataRepo := infrastructure.NewRawDataRepository(pool)
	companyRepo := infrastructure.NewCompanyRepository(pool)
	sectorRepo := infrastructure.NewSectorRepository(pool)
	dividendsRepo := infrastructure.NewDividendsRepository(pool)
	cbRateRepo := infrastructure.NewCBRateRepository(pool)
	newsRepo := infrastructure.NewNewsRepository(pool)

	priceProvider := infrastructure.NewMoexDataProvider()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(application.CORSMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	application.RegisterRatiosRoutes(r, ratiosRepo)
	application.RegisterRawDataRoutes(r, rawDataRepo)
	application.RegisterCompanyRoutes(r, companyRepo, priceProvider)
	application.RegisterSectorRoutes(r, sectorRepo)
	application.RegisterDividendsRoutes(r, dividendsRepo)
	application.RegisterMacroRoutes(r, cbRateRepo)
	application.RegisterNewsRoutes(r, newsRepo)
	application.RegisterPriceRoutes(r, priceProvider)

	srv := &http.Server{
		Addr:         ":8082",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Starting server", "address", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdown:
		slog.Info("Shutdown signal received", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			return err
		}

		slog.Info("Server stopped gracefully")
	}

	return nil
}
