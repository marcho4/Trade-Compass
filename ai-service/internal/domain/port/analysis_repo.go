package port

import (
	"ai-service/internal/domain/entity"
	"context"
)

type AnalysisRepository interface {
	SaveAnalysis(ctx context.Context, result, ticker string, year, period int) error
	GetAnalysis(ctx context.Context, ticker string, year, period int) (string, error)
	GetAvailablePeriods(ctx context.Context, ticker string) ([]entity.AvailablePeriod, error)
}
