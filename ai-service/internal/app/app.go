package app

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	httpserver "ai-service/internal/adapters/http"
	kafkaadapter "ai-service/internal/adapters/kafka"
	"ai-service/internal/config"
	financialdata "ai-service/internal/gateway/financial_data"
	geminigw "ai-service/internal/gateway/gemini"
	kafkagw "ai-service/internal/gateway/kafka"
	"ai-service/internal/gateway/parser"
	"ai-service/internal/gateway/redis"
	"ai-service/internal/gateway/s3"
	"ai-service/internal/repository/postgres"
	"ai-service/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg         *config.Config
	server      *httpserver.HttpServer
	consumer    *kafkaadapter.Consumer
	kafkaClient *kafkagw.KafkaClient
	redisClient *redis.Client
	pool        *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	if err := postgres.RunMigrations(); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	// repositories
	analysisRepo := postgres.NewAnalysisRepository(pool)
	businessResearchRepo := postgres.NewBusinessResearchRepository(pool)
	newsRepo := postgres.NewNewsRepository(pool)
	reportResultsRepo := postgres.NewReportResultsRepository(pool)
	riskAndGrowthRepo := postgres.NewRiskAndGrowthRepository(pool)
	taskRepo := postgres.NewTasksRepository(pool)
	scenarioRepo := postgres.NewScenarioRepository(pool)
	dcfRepo := postgres.NewDCFResultsRepository(pool)
	transactor := postgres.NewPgxTransactor(pool)

	// gateways
	geminiClient, err := geminigw.NewClient(cfg.GeminiAPIKey, cfg.GeminiProxyURL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create gemini client: %w", err)
	}

	s3Client, err := s3.NewClient(cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3BucketName, cfg.S3Endpoint)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create s3 client: %w", err)
	}

	redisClient, err := redis.NewClient(ctx, cfg.RedisURL, cfg.RedisPassword)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create redis client: %w", err)
	}

	fdClient := financialdata.NewClient(cfg.FinancialDataURL, cfg.FinancialDataAPIKey)
	parserClient := parser.NewClient(cfg.ParserURL)
	kafkaClient := kafkagw.NewKafkaClient(cfg.KafkaURL, cfg.KafkaTopic)

	// usecases
	analysisUC := usecase.NewGetAnalysisUsecase(analysisRepo)
	reportResultsUC := usecase.NewGetReportResultsUsecase(reportResultsRepo)
	businessResearchUC := usecase.NewBusinessResearchUsecase(geminiClient, businessResearchRepo, kafkaClient)

	analyzeReportUC := usecase.NewAnalyzeReportUsecase(geminiClient, analysisRepo, kafkaClient, fdClient, s3Client, newsRepo, businessResearchRepo, riskAndGrowthRepo, scenarioRepo, dcfRepo)
	extractRawDataUC := usecase.NewExtractRawDataUsecase(geminiClient, fdClient, parserClient, kafkaClient, s3Client)
	extractResultUC := usecase.NewExtractResultUsecase(geminiClient, reportResultsRepo, analysisRepo, taskRepo)
	newsResearchUC := usecase.NewNewsResearchUsecase(geminiClient, newsRepo, kafkaClient, cfg.NewsTTL)
	riskAndGrowthUC := usecase.NewRiskAndGrowthUsecase(geminiClient, riskAndGrowthRepo, newsRepo, businessResearchRepo, kafkaClient, cfg.NewsTTL)
	scenarioGeneratorUC := usecase.NewScenarioGenerator(geminiClient, fdClient, parserClient, riskAndGrowthRepo, scenarioRepo, dcfRepo, transactor, kafkaClient)
	taskCounterUC := usecase.NewTaskCounterUsecase(taskRepo, transactor, kafkaClient)

	// adapters
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("parse port: %w", err)
	}

	analysisHandler := httpserver.NewAnalysisHandler(analysisUC, reportResultsUC, businessResearchUC)
	server := httpserver.NewHttpServer(analysisHandler)
	server.RegisterRoutes(port, cfg.APIKey)

	dispatcher := kafkaadapter.NewTaskDispatcher(
		analyzeReportUC,
		extractRawDataUC,
		extractResultUC,
		businessResearchUC,
		newsResearchUC,
		riskAndGrowthUC,
		scenarioGeneratorUC,
		taskCounterUC,
	)
	consumer := kafkaadapter.NewConsumer(kafkaClient, dispatcher, 10)

	return &App{
		cfg:         cfg,
		server:      server,
		consumer:    consumer,
		kafkaClient: kafkaClient,
		redisClient: redisClient,
		pool:        pool,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.consumer.Start(ctx)

	return a.server.RunServer(ctx)
}

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	a.consumer.Stop(ctx)

	if err := a.kafkaClient.Close(); err != nil {
		slog.Error("failed to close kafka client", slog.Any("error", err))
	}

	if err := a.redisClient.Close(); err != nil {
		slog.Error("failed to close redis client", slog.Any("error", err))
	}

	a.pool.Close()

	slog.Info("app stopped gracefully")
}
