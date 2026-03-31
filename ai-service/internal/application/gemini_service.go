package application

import (
	"ai-service/internal/domain/entity"
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
	newsTTL       time.Duration
}

func NewGeminiService(
	client *gemini.Client,
	s3Client *s3.Client,
	finDataClient *financialdata.Client,
	db *postgres.DBRepo,
	newsTTL time.Duration,
) *GeminiService {
	return &GeminiService{
		geminiClient:  client,
		finDataClient: finDataClient,
		s3Client:      s3Client,
		db:            db,
		newsTTL:       newsTTL,
	}
}

func (g *GeminiService) AnalyzeReport(ctx context.Context, ticker, reportUrl string, year int, period entity.ReportPeriod) (string, error) {
	logger := slog.With(slog.String("ticker", ticker),
		slog.Int("year", year),
		slog.String("period", string(period)),
		slog.String("report_url", reportUrl))

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

	news, err := g.db.GetFreshNews(ctx, ticker, g.newsTTL)
	if err != nil {
		logger.Warn("failed to get news from DB, continuing without it",
			slog.Any("error", err),
		)
	}
	if news == nil {
		logger.Warn("no fresh news found in DB")
	}

	rawDataHistory, err := g.finDataClient.GetRawDataHistory(ctx, ticker)
	if err != nil {
		logger.Warn("[AnalyzeReport] failed to get raw data history, continuing without it",
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

	logger.Info("downloading PDF")
	pdfBytes, err := g.s3Client.DownloadPDF(ctx, reportUrl)
	if err != nil {
		return "", fmt.Errorf("failed to download PDF: %w", err)
	}
	logger.Info("PDF downloaded", slog.Int("pdf_size_bytes", len(pdfBytes)))

	logger.Info("calling Gemini API..")
	start := time.Now()
	response, err := g.geminiClient.AnalyzeWithPDF(ctx, pdfBytes, prompt, entity.Pro)
	if err != nil {
		return "", fmt.Errorf("failed to generate analysis: %w", err)
	}
	logger.Info("Gemini response received",
		slog.Int("response_length", len(response)),
		slog.String("duration", time.Since(start).String()),
	)

	return response, nil
}

func (g *GeminiService) GetChatResponse(ctx context.Context, prompt string, chatContext entity.ChatContext) (string, error) {
	return "", nil
}

func (g *GeminiService) AnalyzeSector(ctx context.Context, sectorId int) (string, error) {
	return "", nil
}

func (g *GeminiService) GetCompanyHistory(ctx context.Context, ticker string) (string, error) {
	return "", nil
}

func (g *GeminiService) ExtractResultFromReport(ctx context.Context, ticker string, year int, period entity.ReportPeriod) (*entity.ReportResults, error) {
	reportText, err := g.db.GetAnalysis(ctx, ticker, year, entity.PeriodToMonths[string(period)])
	if err != nil {
		return nil, err
	}
	slog.Info("Extracted report from database")

	prompt := BuildExtractPrompt(reportText)
	slog.Info("Generating Report results...")

	text, err := g.geminiClient.GenerateText(ctx, prompt, entity.Flash)
	if err != nil {
		return nil, err
	}

	var res entity.ReportResults
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		slog.Error("Unable to parse text", slog.String("ai response", text))
		return nil, err
	}
	return &res, nil
}

func (g *GeminiService) CollectNews(ctx context.Context, ticker string) (*entity.NewsResponse, error) {
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
	text, err := g.geminiClient.GenerateText(ctx, prompt, entity.Flash,
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

	var res entity.NewsResponse
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		slog.Error("failed to parse news response", slog.String("ai_response", text))
		return nil, fmt.Errorf("parse news response: %w", err)
	}

	return &res, nil
}

func (g *GeminiService) ResearchBusiness(ctx context.Context, ticker, companyName string) (*entity.BusinessResearchResponse, error) {
	prompt := docs.BusinessResearcherPrompt() + "\n\n## Компания для анализа\nТикер: " + ticker + "\nНазвание: " + companyName + "\n\nВАЖНО: В поле ticker ответа используй СТРОГО \"" + ticker + "\". Не заменяй тикер на альтернативный."

	marketSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"market": {Type: genai.TypeString},
			"role":   {Type: genai.TypeString},
		},
		Required: []string{"market", "role"},
	}

	revenueSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"segment":     {Type: genai.TypeString},
			"share_pct":   {Type: genai.TypeNumber},
			"approximate": {Type: genai.TypeBoolean},
			"description": {Type: genai.TypeString},
			"trend":       {Type: genai.TypeString, Enum: []string{"growing", "stable", "declining"}},
		},
		Required: []string{"segment", "share_pct", "approximate", "description", "trend"},
	}

	dependencySchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"factor":      {Type: genai.TypeString},
			"type":        {Type: genai.TypeString, Enum: []string{"commodity", "currency", "regulation", "macro", "technology", "geopolitics", "infrastructure", "demand"}},
			"severity":    {Type: genai.TypeString, Enum: []string{"critical", "high", "moderate"}},
			"description": {Type: genai.TypeString},
		},
		Required: []string{"factor", "type", "severity", "description"},
	}

	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"ticker":       {Type: genai.TypeString},
			"company_name": {Type: genai.TypeString},
			"profile": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"description":           {Type: genai.TypeString},
					"products_and_services": {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}},
					"markets":               {Type: genai.TypeArray, Items: marketSchema},
					"key_clients":           {Type: genai.TypeString},
					"business_model":        {Type: genai.TypeString},
				},
				Required: []string{"description", "products_and_services", "markets", "key_clients", "business_model"},
			},
			"revenue_sources": {Type: genai.TypeArray, Items: revenueSchema},
			"dependencies":    {Type: genai.TypeArray, Items: dependencySchema},
		},
		Required: []string{"ticker", "company_name", "profile", "revenue_sources", "dependencies"},
	}

	slog.Info("[ResearchBusiness] calling Gemini Flash", slog.String("ticker", ticker))

	text, err := g.geminiClient.GenerateText(ctx, prompt, entity.Flash,
		gemini.WithGoogleSearch(),
		gemini.WithResponseSchema(responseSchema),
	)
	if err != nil {
		return nil, fmt.Errorf("research business: %w", err)
	}

	var res entity.BusinessResearchResponse
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		slog.Error("[ResearchBusiness] failed to parse response", slog.String("ai_response", text))
		return nil, fmt.Errorf("parse business research response: %w", err)
	}

	res.Ticker = ticker

	slog.Info("[ResearchBusiness] completed", slog.String("ticker", ticker), slog.String("company_name", res.CompanyName))
	return &res, nil
}

func (g *GeminiService) ExtractRawData(ctx context.Context, ticker, reportUrl string, year int, period entity.ReportPeriod) (*entity.RawData, error) {
	prompt := docs.RawDataAgentPrompt() + "\n<ticker>" + ticker + "</ticker>"

	slog.Info("[Extract Raw Data] downloading PDF", slog.String("report_url", reportUrl))
	pdfBytes, err := g.s3Client.DownloadPDF(ctx, reportUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}
	slog.Info("[Extract Raw Data] PDF downloaded", slog.Int("pdf_size_bytes", len(pdfBytes)))

	text, err := g.geminiClient.AnalyzeWithPDF(ctx, pdfBytes, prompt, entity.Flash)
	if err != nil {
		return nil, fmt.Errorf("failed to extract from PDF: %w", err)
	}

	var rawData entity.RawData
	err = json.Unmarshal([]byte(text), &rawData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal raw data: %w", err)
	}

	return &rawData, nil
}
