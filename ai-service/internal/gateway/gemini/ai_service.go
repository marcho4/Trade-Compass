package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"ai-service/internal/domain/entity"
	docs "ai-service/internal/docs"

	"google.golang.org/genai"
)

type StorageClient interface {
	DownloadPDF(ctx context.Context, url string) ([]byte, error)
}

type FinancialDataClient interface {
	GetDailyPrices(ctx context.Context, ticker string) ([]entity.Candle, error)
	GetCBRates(ctx context.Context) (*entity.CBRate, error)
	GetMarketCap(ctx context.Context, ticker string) (float64, error)
	GetRawDataHistory(ctx context.Context, ticker string) ([]entity.RawData, error)
}

type NewsReader interface {
	GetFreshNews(ctx context.Context, ticker string, ttl time.Duration) (*entity.NewsResponse, error)
}

type AnalysisReader interface {
	GetAnalysis(ctx context.Context, ticker string, year, period int) (string, error)
}

type AIService struct {
	client   *Client
	storage  StorageClient
	fd       FinancialDataClient
	news     NewsReader
	analysis AnalysisReader
	newsTTL  time.Duration
}

func NewAIService(
	client *Client,
	storage StorageClient,
	fd FinancialDataClient,
	news NewsReader,
	analysis AnalysisReader,
	newsTTL time.Duration,
) *AIService {
	return &AIService{
		client:   client,
		storage:  storage,
		fd:       fd,
		news:     news,
		analysis: analysis,
		newsTTL:  newsTTL,
	}
}

func (s *AIService) AnalyzeReport(ctx context.Context, ticker, reportURL string, year int, period entity.ReportPeriod) (string, error) {
	logger := slog.With(
		slog.String("ticker", ticker),
		slog.Int("year", year),
		slog.String("period", string(period)),
	)

	candles, err := s.fd.GetDailyPrices(ctx, ticker)
	if err != nil {
		return "", fmt.Errorf("get price history: %w", err)
	}

	cbRate, err := s.fd.GetCBRates(ctx)
	if err != nil {
		return "", fmt.Errorf("get CB rate: %w", err)
	}

	marketCap, err := s.fd.GetMarketCap(ctx, ticker)
	if err != nil {
		return "", fmt.Errorf("get market cap: %w", err)
	}

	news, err := s.news.GetFreshNews(ctx, ticker, s.newsTTL)
	if err != nil {
		logger.Warn("failed to get news from DB, continuing without it", slog.Any("error", err))
	}
	if news == nil {
		logger.Warn("no fresh news found in DB")
	}

	rawDataHistory, err := s.fd.GetRawDataHistory(ctx, ticker)
	if err != nil {
		logger.Warn("failed to get raw data history, continuing without it", slog.Any("error", err))
	}

	prompt := buildAnalysisPrompt(analysisContext{
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
	pdfBytes, err := s.storage.DownloadPDF(ctx, reportURL)
	if err != nil {
		return "", fmt.Errorf("download PDF: %w", err)
	}
	logger.Info("PDF downloaded", slog.Int("pdf_size_bytes", len(pdfBytes)))

	logger.Info("calling Gemini API")
	start := time.Now()
	response, err := s.client.AnalyzeWithPDF(ctx, pdfBytes, prompt, entity.Pro)
	if err != nil {
		return "", fmt.Errorf("generate analysis: %w", err)
	}
	logger.Info("Gemini response received",
		slog.Int("response_length", len(response)),
		slog.Duration("duration", time.Since(start)),
	)

	return response, nil
}

func (s *AIService) ExtractRawData(ctx context.Context, ticker, reportURL string, year int, period entity.ReportPeriod) (*entity.RawData, error) {
	prompt := docs.RawDataAgentPrompt() + "\n<ticker>" + ticker + "</ticker>"

	slog.Info("downloading PDF for raw data extraction", slog.String("report_url", reportURL))
	pdfBytes, err := s.storage.DownloadPDF(ctx, reportURL)
	if err != nil {
		return nil, fmt.Errorf("download PDF: %w", err)
	}
	slog.Info("PDF downloaded", slog.Int("pdf_size_bytes", len(pdfBytes)))

	text, err := s.client.AnalyzeWithPDF(ctx, pdfBytes, prompt, entity.Flash)
	if err != nil {
		return nil, fmt.Errorf("extract from PDF: %w", err)
	}

	var rawData entity.RawData
	if err := json.Unmarshal([]byte(text), &rawData); err != nil {
		return nil, fmt.Errorf("unmarshal raw data: %w", err)
	}

	return &rawData, nil
}

func (s *AIService) ExtractResultFromReport(ctx context.Context, ticker string, year int, period entity.ReportPeriod) (*entity.ReportResults, error) {
	reportText, err := s.analysis.GetAnalysis(ctx, ticker, year, entity.PeriodToMonths[string(period)])
	if err != nil {
		return nil, fmt.Errorf("get analysis: %w", err)
	}

	slog.Info("extracted report from database", slog.String("ticker", ticker))

	text, err := s.client.GenerateText(ctx, buildExtractPrompt(reportText), entity.Flash)
	if err != nil {
		return nil, fmt.Errorf("generate extract result: %w", err)
	}

	var res entity.ReportResults
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		slog.Error("unable to parse extract result", slog.String("ai_response", text))
		return nil, fmt.Errorf("unmarshal report results: %w", err)
	}

	return &res, nil
}

func (s *AIService) CollectNews(ctx context.Context, ticker string) (*entity.NewsResponse, error) {
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

	slog.Info("calling Gemini to collect news", slog.String("ticker", ticker))
	text, err := s.client.GenerateText(ctx, buildNewsAgentPrompt(ticker), entity.Flash,
		WithGoogleSearch(),
		WithResponseSchema(&genai.Schema{
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

func (s *AIService) ResearchBusiness(ctx context.Context, ticker, companyName string) (*entity.BusinessResearchResponse, error) {
	prompt := docs.BusinessResearcherPrompt() +
		"\n\n## Компания для анализа\nТикер: " + ticker +
		"\nНазвание: " + companyName +
		"\n\nВАЖНО: В поле ticker ответа используй СТРОГО \"" + ticker + "\". Не заменяй тикер на альтернативный."

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

	slog.Info("calling Gemini for business research", slog.String("ticker", ticker))
	text, err := s.client.GenerateText(ctx, prompt, entity.Flash,
		WithGoogleSearch(),
		WithResponseSchema(responseSchema),
	)
	if err != nil {
		return nil, fmt.Errorf("research business: %w", err)
	}

	var res entity.BusinessResearchResponse
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		slog.Error("failed to parse business research response", slog.String("ai_response", text))
		return nil, fmt.Errorf("parse business research response: %w", err)
	}

	res.Ticker = ticker

	slog.Info("business research completed", slog.String("ticker", ticker), slog.String("company_name", res.CompanyName))
	return &res, nil
}

func (s *AIService) ExtractRiskAndGrowth(ctx context.Context, ticker string, news *entity.NewsResponse, business *entity.BusinessResearchResult) (*entity.RiskAndGrowthResponse, error) {
	newsJSON, err := json.Marshal(news)
	if err != nil {
		return nil, fmt.Errorf("marshal news: %w", err)
	}

	businessJSON, err := json.Marshal(business)
	if err != nil {
		return nil, fmt.Errorf("marshal business: %w", err)
	}

	prompt := docs.RiskAndGrowthPrompt()
	prompt = strings.ReplaceAll(prompt, "{{BUSINESS_RESEARCH}}", string(businessJSON))
	prompt = strings.ReplaceAll(prompt, "{{NEWS}}", string(newsJSON))

	factorSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"name":    {Type: genai.TypeString},
			"type":    {Type: genai.TypeString, Enum: []string{"growth", "risk"}},
			"horizon": {Type: genai.TypeString, Enum: []string{"short_term", "medium_term"}},
			"impact":  {Type: genai.TypeString, Enum: []string{"high", "medium", "low"}},
			"summary": {Type: genai.TypeString},
			"source":  {Type: genai.TypeString},
		},
		Required: []string{"name", "type", "horizon", "impact", "summary", "source"},
	}

	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"ticker":  {Type: genai.TypeString},
			"factors": {Type: genai.TypeArray, Items: factorSchema},
		},
		Required: []string{"ticker", "factors"},
	}

	slog.Info("calling Gemini for risk and growth analysis", slog.String("ticker", ticker))
	text, err := s.client.GenerateText(ctx, prompt, entity.Flash,
		WithResponseSchema(responseSchema),
	)
	if err != nil {
		return nil, fmt.Errorf("extract risk and growth: %w", err)
	}

	var res entity.RiskAndGrowthResponse
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		slog.Error("failed to parse risk and growth response", slog.String("ai_response", text))
		return nil, fmt.Errorf("parse risk and growth response: %w", err)
	}

	res.Ticker = ticker

	slog.Info("risk and growth analysis completed",
		slog.String("ticker", ticker),
		slog.Int("factors_count", len(res.Factors)),
	)
	return &res, nil
}
