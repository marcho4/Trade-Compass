package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"ai-service/internal/domain/entity"
)

type AnalyzeReportUsecase struct {
	analysis         AnalysisRepository
	publisher        MessagePublisher
	finData          FinancialDataGateway
	aiClient         AIProvider
	storage          StorageClient
	news             NewsRepository
	businessResearch BusinessResearchRepository
	riskAndGrowth    RiskAndGrowthRepository
	scenarios        ScenarioRepository
	dcf              DCFResultsRepository
}

func NewAnalyzeReportUsecase(
	ai AIProvider,
	analysis AnalysisRepository,
	publisher MessagePublisher,
	finData FinancialDataGateway,
	storage StorageClient,
	news NewsRepository,
	businessResearch BusinessResearchRepository,
	riskAndGrowth RiskAndGrowthRepository,
	scenarios ScenarioRepository,
	dcf DCFResultsRepository,
) *AnalyzeReportUsecase {
	return &AnalyzeReportUsecase{
		aiClient:         ai,
		analysis:         analysis,
		publisher:        publisher,
		finData:          finData,
		storage:          storage,
		news:             news,
		businessResearch: businessResearch,
		riskAndGrowth:    riskAndGrowth,
		scenarios:        scenarios,
		dcf:              dcf,
	}
}

func (u *AnalyzeReportUsecase) Execute(ctx context.Context, task entity.Task) error {
	logger := slog.With(
		slog.String("ticker", task.Ticker),
		slog.Int("year", task.Year),
		slog.String("period", task.Period),
	)

	logger.Info("starting analyze task")

	candles, err := u.finData.GetDailyPrices(ctx, task.Ticker)
	if err != nil {
		return fmt.Errorf("get price history: %w", err)
	}

	cbRate, err := u.finData.GetCBRates(ctx)
	if err != nil {
		return fmt.Errorf("get CB rate: %w", err)
	}

	marketCap, err := u.finData.GetMarketCap(ctx, task.Ticker)
	if err != nil {
		return fmt.Errorf("get market cap: %w", err)
	}

	news, err := u.news.GetFreshNews(ctx, task.Ticker, 72*time.Hour)
	if err != nil {
		logger.Warn("failed to get news from DB, continuing without it", slog.Any("error", err))
	}
	if news == nil {
		logger.Warn("no fresh news found in DB")
	}

	rawDataHistory, err := u.finData.GetRawDataHistory(ctx, task.Ticker)
	if err != nil {
		logger.Warn("failed to get raw data history, continuing without it", slog.Any("error", err))
	}

	businessResearch, err := u.businessResearch.GetBusinessResearch(ctx, task.Ticker)
	if err != nil {
		logger.Warn("failed to get business research, continuing without it", slog.Any("error", err))
	}

	risksAndGrowth, err := u.riskAndGrowth.GetFreshRiskAndGrowth(ctx, task.Ticker, 72*time.Hour)
	if err != nil {
		logger.Warn("failed to get risk and growth, continuing without it", slog.Any("error", err))
	}

	scenarios, err := u.scenarios.GetScenarios(ctx, task.Ticker)
	if err != nil {
		logger.Warn("failed to get scenarios, continuing without them", slog.Any("error", err))
	}

	dcfResult, err := u.dcf.GetDCFResults(ctx, task.Ticker)
	if err != nil {
		logger.Warn("failed to get dcf results, continuing without them", slog.Any("error", err))
	}

	prompt := buildAnalysisPrompt(analysisContext{
		Ticker:           task.Ticker,
		Year:             task.Year,
		Period:           entity.ReportPeriod(task.Period),
		RawDataHistory:   rawDataHistory,
		Candles:          candles,
		CBRate:           cbRate,
		MarketCap:        marketCap,
		News:             news,
		BusinessResearch: businessResearch,
		RisksAndGrowth:   risksAndGrowth,
		Scenarios:        scenarios,
		DCFResult:        dcfResult,
	})

	logger.Info("downloading PDF")

	pdfBytes, err := u.storage.DownloadPDF(ctx, task.ReportURL)
	if err != nil {
		return fmt.Errorf("download PDF: %w", err)
	}

	logger.Info("PDF downloaded", slog.Int("pdf_size_bytes", len(pdfBytes)))

	logger.Info("calling Gemini API")

	start := time.Now()
	result, err := u.aiClient.AnalyzeWithPDF(ctx, pdfBytes, prompt, entity.Pro)
	if err != nil {
		return fmt.Errorf("generate analysis: %w", err)
	}

	logger.Info("Gemini response received",
		slog.Int("response_length", len(result)),
		slog.Duration("duration", time.Since(start)),
	)

	logger.Info("analyze report succeeded", slog.Int("result_length", len(result)))

	periodMonths, ok := entity.PeriodToMonths[task.Period]
	if !ok {
		return fmt.Errorf("unknown period: %s", task.Period)
	}

	if err := u.analysis.SaveAnalysis(ctx, result, task.Ticker, task.Year, periodMonths); err != nil {
		return fmt.Errorf("save analysis: %w", err)
	}

	nextTask := entity.Task{
		Id:        task.Id,
		Ticker:    task.Ticker,
		Year:      task.Year,
		Period:    task.Period,
		ReportURL: task.ReportURL,
		Type:      entity.ExtractResult,
	}

	payload, err := json.Marshal(nextTask)
	if err != nil {
		return fmt.Errorf("marshal extract-result task: %w", err)
	}

	if err := u.publisher.PublishMessage(ctx, payload); err != nil {
		return fmt.Errorf("publish extract-result task: %w", err)
	}

	logger.Info("published extract-result task")

	return nil
}
