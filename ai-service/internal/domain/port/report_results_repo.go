package port

import (
	"ai-service/internal/domain/entity"
	"context"
)

type ReportResultsRepository interface {
	SaveReportResults(ctx context.Context, result *entity.ReportResults, ticker string, year, period int) error
	GetReportResults(ctx context.Context, ticker string, year, period int) (*entity.ReportResults, error)
	GetLatestReportResults(ctx context.Context, ticker string) (*entity.ReportResults, error)
}
