package usecase

import (
	"context"
	"time"

	"ai-service/internal/domain/entity"
)

type AnalysisRepository interface {
	SaveAnalysis(ctx context.Context, result, ticker string, year, period int) error
	GetAnalysis(ctx context.Context, ticker string, year, period int) (string, error)
	GetAvailablePeriods(ctx context.Context, ticker string) ([]entity.AvailablePeriod, error)
}

type ReportResultsSaver interface {
	SaveReportResults(ctx context.Context, result *entity.ReportResults, ticker string, year, period int) error
}

type NewsRepository interface {
	SaveNews(ctx context.Context, ticker string, news *entity.NewsResponse) error
	GetFreshNews(ctx context.Context, ticker string, ttl time.Duration) (*entity.NewsResponse, error)
}

type BusinessResearchRepository interface {
	SaveBusinessResearch(ctx context.Context, research *entity.BusinessResearchResponse) error
	GetBusinessResearch(ctx context.Context, ticker string) (*entity.BusinessResearchResult, error)
}

type RiskAndGrowthRepository interface {
	SaveRiskAndGrowth(ctx context.Context, response *entity.RiskAndGrowthResponse) error
	GetFreshRiskAndGrowth(ctx context.Context, ticker string, ttl time.Duration) (*entity.RiskAndGrowthResponse, error)
}

type TasksRepository interface {
	IncrementPending(ctx context.Context, taskID, taskType string, count int) error
	DecrementPending(ctx context.Context, taskID, taskType string) (int, error)
	DeleteTask(ctx context.Context, taskID string) error
	CheckIfTaskIsReady(ctx context.Context, taskID string, expectedTasks int) (bool, error)
}

type ScenarioRepository interface {
	SaveScenarios(ctx context.Context, ticker string, scenarios []entity.Scenario) error
	GetScenarios(ctx context.Context, ticker string) ([]entity.Scenario, error)
}

type DCFResultsRepository interface {
	SaveDCFResults(ctx context.Context, ticker string, result entity.DCFResult) error
	GetDCFResults(ctx context.Context, ticker string) (*entity.DCFResult, error)
}

type Transactor interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
