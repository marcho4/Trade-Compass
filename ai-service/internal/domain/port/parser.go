package port

import (
	"ai-service/internal/domain/entity"
	"context"
)

type ParserProvider interface {
	GetReportS3Path(ctx context.Context, ticker, period string, year int) (string, error)
	GetLatestReportYear(ctx context.Context, ticker, period string) (int, error)
	GetReports(ctx context.Context, ticker string) ([]entity.Report, error)
	IsLatestReport(ctx context.Context, ticker string, year int, periodMonths int) (bool, error)
}
