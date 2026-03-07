package application

import (
	"context"
	"errors"
	"financial_data/internal/application/routers"
	"financial_data/internal/domain"
	"financial_data/internal/infrastructure"
	"financial_data/internal/infrastructure/kafka"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	myMiddleware "financial_data/internal/application/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type FinData struct {
	ratiosRepo     routers.RatiosRepository
	rawDataRepo    routers.RawDataRepository
	companyRepo    routers.CompanyRepository
	sectorRepo     routers.SectorRepository
	dividendsRepo  routers.DividendsRepository
	cbRateRepo     routers.MacroDataRepository
	newsRepo       routers.NewsRepository
	marketService  domain.MarketService
	eventPublisher routers.EventPublisher
	kafkaProducer  *kafka.Producer
	pool           *pgxpool.Pool
	redisClient    *redis.Client
	srv            *http.Server
}

func NewFinData() (*FinData, error) {
	ctx := context.Background()

	if err := infrastructure.RunMigrations(); err != nil {
		return nil, err
	}

	pool, err := infrastructure.NewPostgresPool(ctx)
	if err != nil {
		return nil, err
	}

	slog.Info("database connection pool initialized")

	redisClient, err := infrastructure.NewRedisClient(ctx)
	if err != nil {
		return nil, err
	}

	slog.Info("redis client initialized")

	ratiosRepo := infrastructure.NewRatiosRepository(pool)
	rawDataRepo := infrastructure.NewRawDataRepository(pool)
	companyRepo := infrastructure.NewCompanyRepository(pool, redisClient)
	sectorRepo := infrastructure.NewSectorRepository(pool)
	dividendsRepo := infrastructure.NewDividendsRepository(pool)
	cbRateRepo := infrastructure.NewCBRateRepository(pool)
	newsRepo := infrastructure.NewNewsRepository(pool)

	moexDataProvider := infrastructure.NewMoexDataProvider(redisClient)

	kafkaBrokers := []string{getEnv("KAFKA_URL", "kafka:9092")}
	parserTopic := getEnv("KAFKA_PARSER_TOPIC", "parser.parse_ticker")
	kafkaProducer := kafka.NewProducer(kafkaBrokers, parserTopic)
	eventPublisher := kafka.NewKafkaEventPublisher(kafkaProducer)

	slog.Info("Kafka producer initialized", "topic", parserTopic)

	return &FinData{
		ratiosRepo:     ratiosRepo,
		rawDataRepo:    rawDataRepo,
		companyRepo:    companyRepo,
		sectorRepo:     sectorRepo,
		dividendsRepo:  dividendsRepo,
		cbRateRepo:     cbRateRepo,
		newsRepo:       newsRepo,
		marketService:  moexDataProvider,
		eventPublisher: eventPublisher,
		kafkaProducer:  kafkaProducer,
		pool:           pool,
		redisClient:    redisClient,
	}, nil
}

func (f *FinData) Run() error {
	f.prepareRouter()
	return f.runRouter()
}

func (f *FinData) prepareRouter() error {
	m, err := myMiddleware.NewMiddlewareConfig()
	if err != nil {
		return fmt.Errorf("create middleware config %w", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(m.CORSMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	routers.RegisterRatiosRoutes(r, f.ratiosRepo, m)
	routers.RegisterRawDataRoutes(r, f.rawDataRepo, m)
	routers.RegisterCompanyRoutes(r, f.companyRepo, f.marketService, f.eventPublisher, m)
	routers.RegisterSectorRoutes(r, f.sectorRepo, m)
	routers.RegisterDividendsRoutes(r, f.dividendsRepo, m)
	routers.RegisterMacroRoutes(r, f.cbRateRepo, m)
	routers.RegisterNewsRoutes(r, f.newsRepo, m)
	routers.RegisterPriceRoutes(r, f.marketService, m)

	srv := &http.Server{
		Addr:         ":8082",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	f.srv = srv
	return nil
}

func (f *FinData) runRouter() error {
	if f.srv == nil {
		return errors.New("srv is not initialize")
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("starting server", "address", f.srv.Addr)
		serverErrors <- f.srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		f.kafkaProducer.Close()
		f.pool.Close()
		f.redisClient.Close()
		return err
	case sig := <-shutdown:
		slog.Info("Shutdown signal received", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := f.srv.Shutdown(ctx); err != nil {
			f.srv.Close()
			return err
		}

		f.kafkaProducer.Close()
		f.pool.Close()
		f.redisClient.Close()
		slog.Info("Server stopped gracefully")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
