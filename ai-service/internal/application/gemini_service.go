package application

import (
	"ai-service/internal/domain"
	"ai-service/internal/infrastructure/financialdata"
	"ai-service/internal/infrastructure/gemini"
	"ai-service/internal/infrastructure/postgres"
	"ai-service/internal/infrastructure/s3"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

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

	slog.Info("[AnalyzeReport] fetching daily prices", slog.String("ticker", ticker))
	candles, err := g.finDataClient.GetDailyPrices(ctx, ticker)
	if err != nil {
		return "", fmt.Errorf("failed to get price history: %w", err)
	}
	slog.Info("[AnalyzeReport] daily prices fetched", slog.String("ticker", ticker), slog.Int("candles_count", len(candles)))

	slog.Info("[AnalyzeReport] fetching CB rates")
	cbRate, err := g.finDataClient.GetCBRates(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get CB rate: %w", err)
	}
	slog.Info("[AnalyzeReport] CB rate fetched", slog.Float64("rate", cbRate.Rate))

	slog.Info("[AnalyzeReport] fetching market cap", slog.String("ticker", ticker))
	marketCap, err := g.finDataClient.GetMarketCap(ctx, ticker)
	if err != nil {
		return "", fmt.Errorf("failed to get market cap: %w", err)
	}
	slog.Info("[AnalyzeReport] market cap fetched", slog.Float64("market_cap", marketCap))

	prompt := BuildAnalysisPrompt(AnalysisContext{
		Ticker:    ticker,
		Year:      year,
		Period:    period,
		Candles:   candles,
		CBRate:    cbRate,
		MarketCap: marketCap,
	})
	slog.Info("[AnalyzeReport] prompt built", slog.Int("prompt_length", len(prompt)))

	slog.Info("[AnalyzeReport] downloading PDF", slog.String("report_url", reportUrl))
	pdfBytes, err := g.s3Client.DownloadPDF(ctx, reportUrl)
	if err != nil {
		return "", fmt.Errorf("failed to download PDF: %w", err)
	}
	slog.Info("[AnalyzeReport] PDF downloaded", slog.Int("pdf_size_bytes", len(pdfBytes)))

	slog.Info("[AnalyzeReport] calling Gemini API")
	response, err := g.geminiClient.AnalyzeWithPDF(ctx, pdfBytes, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate analysis: %w", err)
	}
	slog.Info("[AnalyzeReport] Gemini response received", slog.Int("response_length", len(response)))

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
