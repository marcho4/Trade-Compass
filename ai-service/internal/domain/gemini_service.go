package domain

import "context"

type GeminiService interface {
	AnalyzeReport(ctx context.Context, ticker, reportUrl string, year int, period ReportPeriod) (string, error)
	GetChatResponse(ctx context.Context, prompt string, chatContext ChatContext) (string, error)
	AnalyzeSector(ctx context.Context, sectorId int) (string, error)
	GetCompanyHistory(ctx context.Context, ticker string) (string, error)
	ExtractResultFromReport(ctx context.Context, ticker string, year int, period ReportPeriod) (*ReportResults, error)
	CollectNews(ctx context.Context, ticker string) (*NewsResponse, error)
}
