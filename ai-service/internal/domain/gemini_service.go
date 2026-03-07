package domain

import (
	"ai-service/internal/domain/entity"
	"context"
)

type GeminiService interface {
	AnalyzeReport(ctx context.Context, ticker, reportUrl string, year int, period entity.ReportPeriod) (string, error)
	GetChatResponse(ctx context.Context, prompt string, chatContext entity.ChatContext) (string, error)
	AnalyzeSector(ctx context.Context, sectorId int) (string, error)
	GetCompanyHistory(ctx context.Context, ticker string) (string, error)
	ExtractResultFromReport(ctx context.Context, ticker string, year int, period entity.ReportPeriod) (*entity.ReportResults, error)
	CollectNews(ctx context.Context, ticker string) (*entity.NewsResponse, error)
	ExtractRawData(ctx context.Context, ticker, reportUrl string, year int, period entity.ReportPeriod) (*entity.RawData, error)
}
