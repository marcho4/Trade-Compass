package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"ai-service/internal/domain/entity"
)

type ExtractResultUsecase struct {
	ai      AIService
	results ReportResultsSaver
}

func NewExtractResultUsecase(ai AIService, results ReportResultsSaver) *ExtractResultUsecase {
	return &ExtractResultUsecase{
		ai:      ai,
		results: results,
	}
}

func (u *ExtractResultUsecase) Execute(ctx context.Context, task entity.Task) error {
	result, err := u.ai.ExtractResultFromReport(ctx, task.Ticker, task.Year, entity.ReportPeriod(task.Period))
	if err != nil {
		return fmt.Errorf("extract results from report: %w", err)
	}

	slog.Info("extracted results from report analysis",
		slog.String("ticker", task.Ticker),
		slog.Int("year", task.Year),
		slog.String("period", task.Period),
	)

	periodMonths, ok := entity.PeriodToMonths[task.Period]
	if !ok {
		return fmt.Errorf("unknown period: %s", task.Period)
	}

	if err := u.results.SaveReportResults(ctx, result, task.Ticker, task.Year, periodMonths); err != nil {
		return fmt.Errorf("save report results: %w", err)
	}

	return nil
}
