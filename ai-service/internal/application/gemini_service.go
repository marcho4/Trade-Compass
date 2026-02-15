package application

import (
	"ai-service/internal/domain"
	"ai-service/internal/infrastructure/financialdata"
	"ai-service/internal/infrastructure/gemini"
	"ai-service/internal/infrastructure/s3"
	"context"
	"fmt"
	"log/slog"
)

type GeminiService struct {
	geminiClient  *gemini.Client
	s3Client      *s3.Client
	finDataClient *financialdata.Client
}

func NewGeminiService(
	client *gemini.Client,
	s3Client *s3.Client,
	finDataClient *financialdata.Client,
) *GeminiService {
	return &GeminiService{
		geminiClient:  client,
		finDataClient: finDataClient,
		s3Client:      s3Client,
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

func (g *GeminiService) GetChatResponse(prompt string, chatContext domain.ChatContext) (string, error) {
	return "", nil
}

func (g *GeminiService) AnalyzeSector(sectorId int) (string, error) {
	return "", nil
}

func (g *GeminiService) GetCompanyHistory(ticker string) (string, error) {
	return "", nil
}

func (g *GeminiService) ExtractDataFromReport(ctx context.Context, ticker, reportUrl string, year int, period domain.ReportPeriod) (string, error) {
	return "", nil
}
