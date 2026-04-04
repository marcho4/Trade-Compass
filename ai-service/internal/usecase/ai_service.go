package usecase

import (
	"context"

	"ai-service/internal/domain/entity"
)

type AIService interface {
	AnalyzeReport(ctx context.Context, ticker, reportURL string, year int, period entity.ReportPeriod) (string, error)
	ExtractRawData(ctx context.Context, ticker, reportURL string, year int, period entity.ReportPeriod) (*entity.RawData, error)
	ExtractResultFromReport(ctx context.Context, ticker string, year int, period entity.ReportPeriod) (*entity.ReportResults, error)
	CollectNews(ctx context.Context, ticker string) (*entity.NewsResponse, error)
	ResearchBusiness(ctx context.Context, ticker string) (*entity.BusinessResearchResponse, error)
	ExtractRiskAndGrowth(ctx context.Context, ticker string, news *entity.NewsResponse, business *entity.BusinessResearchResult) (*entity.RiskAndGrowthResponse, error)
	GenerateScenarios(ctx context.Context, ticker string, years int, history []entity.RawData, cbRate *entity.CBRate, wacc float64, riskAndGrowth *entity.RiskAndGrowthResponse) ([]entity.Scenario, error)
}
