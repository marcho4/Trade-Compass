package application

import (
	"ai-service/internal/domain"
	docs "ai-service/internal/infrastructure/docs"
	"ai-service/internal/infrastructure/financialdata"
	"ai-service/internal/infrastructure/gemini"
	"ai-service/internal/infrastructure/postgres"
	"ai-service/internal/infrastructure/s3"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/genai"
)

type GeminiService struct {
	geminiClient  *gemini.Client
	s3Client      *s3.Client
	finDataClient *financialdata.Client
	db            *postgres.DBRepo
}

func NewGeminiService(
	client *gemini.Client,
	s3Client *s3.Client,
	finDataClient *financialdata.Client,
	db *postgres.DBRepo,
) *GeminiService {
	return &GeminiService{
		geminiClient:  client,
		finDataClient: finDataClient,
		s3Client:      s3Client,
		db:            db,
	}
}

func (g *GeminiService) AnalyzeReport(ctx context.Context, ticker, reportUrl string, year int, period domain.ReportPeriod) (string, error) {
	slog.Info("[AnalyzeReport] started",
		slog.String("ticker", ticker),
		slog.Int("year", year),
		slog.String("period", string(period)),
		slog.String("report_url", reportUrl),
	)

	candles, err := g.finDataClient.GetDailyPrices(ctx, ticker)
	if err != nil {
		return "", fmt.Errorf("failed to get price history: %w", err)
	}

	cbRate, err := g.finDataClient.GetCBRates(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get CB rate: %w", err)
	}

	marketCap, err := g.finDataClient.GetMarketCap(ctx, ticker)
	if err != nil {
		return "", fmt.Errorf("failed to get market cap: %w", err)
	}

	news, err := g.CollectNews(ctx, ticker)
	if err != nil {
		return "", fmt.Errorf("news: %w", err)
	}

	rawDataHistory, err := g.finDataClient.GetRawDataHistory(ctx, ticker)
	if err != nil {
		slog.Warn("[AnalyzeReport] failed to get raw data history, continuing without it",
			slog.String("ticker", ticker),
			slog.Any("error", err),
		)
	}

	prompt := BuildAnalysisPrompt(AnalysisContext{
		Ticker:         ticker,
		Year:           year,
		Period:         period,
		RawDataHistory: rawDataHistory,
		Candles:        candles,
		CBRate:         cbRate,
		MarketCap:      marketCap,
		News:           news,
	})

	slog.Info("[AnalyzeReport] downloading PDF", slog.String("report_url", reportUrl))
	pdfBytes, err := g.s3Client.DownloadPDF(ctx, reportUrl)
	if err != nil {
		return "", fmt.Errorf("failed to download PDF: %w", err)
	}
	slog.Info("[AnalyzeReport] PDF downloaded", slog.Int("pdf_size_bytes", len(pdfBytes)))

	slog.Info("[AnalyzeReport] calling Gemini API..", slog.String("ticker", ticker))
	start := time.Now()
	response, err := g.geminiClient.AnalyzeWithPDF(ctx, pdfBytes, prompt, domain.Pro)
	if err != nil {
		return "", fmt.Errorf("failed to generate analysis: %w", err)
	}
	slog.Info("[AnalyzeReport] Gemini response received",
		slog.Int("response_length", len(response)),
		slog.String("duration", time.Since(start).String()),
		slog.String("ticker", ticker),
	)

	return response, nil
}

func (g *GeminiService) GetChatResponse(ctx context.Context, prompt string, chatContext domain.ChatContext) (string, error) {
	return "", nil
}

func (g *GeminiService) AnalyzeSector(ctx context.Context, sectorId int) (string, error) {
	return "", nil
}

func (g *GeminiService) GetCompanyHistory(ctx context.Context, ticker string) (string, error) {
	return "", nil
}

func (g *GeminiService) ExtractResultFromReport(ctx context.Context, ticker string, year int, period domain.ReportPeriod) (*domain.ReportResults, error) {
	reportText, err := g.db.GetAnalysis(ctx, ticker, year, domain.PeriodToMonths[string(period)])
	if err != nil {
		return nil, err
	}
	slog.Info("Extracted report from database")

	prompt := BuildExtractPrompt(reportText)
	slog.Info("Generating Report results...")

	text, err := g.geminiClient.GenerateText(ctx, prompt, domain.Flash)
	if err != nil {
		return nil, err
	}

	var res domain.ReportResults
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		slog.Error("Unable to parse text", slog.String("ai response", text))
		return nil, err
	}
	return &res, nil
}

func (g *GeminiService) CollectNews(ctx context.Context, ticker string) (*domain.NewsResponse, error) {
	prompt := BuildNewsAgentPrompt(ticker)

	newsItemSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"news":        {Type: genai.TypeString},
			"date":        {Type: genai.TypeString},
			"source":      {Type: genai.TypeString},
			"severity":    {Type: genai.TypeString, Enum: []string{"high", "medium", "low"}},
			"impact_type": {Type: genai.TypeString, Enum: []string{"positive", "negative", "neutral"}},
		},
		Required: []string{"news", "date", "source", "severity", "impact_type"},
	}

	slog.Info("Calling Gemini to get news", slog.Any("ticker", ticker))
	text, err := g.geminiClient.GenerateText(ctx, prompt, domain.Flash,
		gemini.WithGoogleSearch(),
		gemini.WithResponseSchema(&genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"latest_news":    {Type: genai.TypeArray, Items: newsItemSchema},
				"important_news": {Type: genai.TypeArray, Items: newsItemSchema},
			},
			Required: []string{"latest_news", "important_news"},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("collect news: %w", err)
	}

	var res domain.NewsResponse
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		slog.Error("failed to parse news response", slog.String("ai_response", text))
		return nil, fmt.Errorf("parse news response: %w", err)
	}

	return &res, nil
}

func (g *GeminiService) ExtractRawData(ctx context.Context, ticker, reportUrl string, year int, period domain.ReportPeriod) (*domain.RawData, error) {
	prompt := docs.RawDataAgentPrompt() + "\n<ticker>" + ticker + "</ticker>"

	slog.Info("[Extract Raw Data] downloading PDF", slog.String("report_url", reportUrl))
	pdfBytes, err := g.s3Client.DownloadPDF(ctx, reportUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}
	slog.Info("[Extract Raw Data] PDF downloaded", slog.Int("pdf_size_bytes", len(pdfBytes)))

	text, err := g.geminiClient.AnalyzeWithPDF(ctx, pdfBytes, prompt, domain.Flash)
	if err != nil {
		return nil, fmt.Errorf("failed to extract from PDF: %w", err)
	}

	var rawData domain.RawData
	err = json.Unmarshal([]byte(text), &rawData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal raw data: %w", err)
	}

	return &rawData, nil
}
