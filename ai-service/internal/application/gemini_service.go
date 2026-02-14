package application

import (
	"ai-service/internal/domain"
	"ai-service/internal/infrastructure/financialdata"
	"ai-service/internal/infrastructure/gemini"
	"ai-service/internal/infrastructure/s3"
	"context"
	"fmt"
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
	// rawData, err := g.finDataClient.GetDraft(ctx, ticker, year, period)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to get financial data: %w", err)
	// }

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

	prompt := BuildAnalysisPrompt(AnalysisContext{
		Ticker: ticker,
		Year:   year,
		Period: period,
		// RawData:   rawData,
		Candles:   candles,
		CBRate:    cbRate,
		MarketCap: marketCap,
	})

	pdfBytes, err := g.s3Client.DownloadPDF(ctx, reportUrl)
	if err != nil {
		return "", fmt.Errorf("failed to download PDF: %w", err)
	}

	response, err := g.geminiClient.AnalyzeWithPDF(ctx, pdfBytes, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate analysis: %w", err)
	}

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

func (g *GeminiService) ExtractDataFromReport(reportUrl string) (string, error) {
	return "", nil
}
