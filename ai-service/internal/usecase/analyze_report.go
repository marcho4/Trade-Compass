package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"ai-service/internal/domain/entity"
)

type AnalyzeReportUsecase struct {
	ai        AIService
	analysis  AnalysisRepository
	publisher MessagePublisher
}

func NewAnalyzeReportUsecase(ai AIService, analysis AnalysisRepository, publisher MessagePublisher) *AnalyzeReportUsecase {
	return &AnalyzeReportUsecase{
		ai:        ai,
		analysis:  analysis,
		publisher: publisher,
	}
}

func (u *AnalyzeReportUsecase) Execute(ctx context.Context, task entity.Task) error {
	logger := slog.With(
		slog.String("ticker", task.Ticker),
		slog.Int("year", task.Year),
		slog.String("period", task.Period),
	)

	logger.Info("starting analyze task")

	result, err := u.ai.AnalyzeReport(ctx, task.Ticker, task.ReportURL, task.Year, entity.ReportPeriod(task.Period))
	if err != nil {
		return fmt.Errorf("analyze report: %w", err)
	}

	logger.Info("analyze report succeeded", slog.Int("result_length", len(result)))

	periodMonths, ok := entity.PeriodToMonths[task.Period]
	if !ok {
		return fmt.Errorf("unknown period: %s", task.Period)
	}

	if err := u.analysis.SaveAnalysis(ctx, result, task.Ticker, task.Year, periodMonths); err != nil {
		return fmt.Errorf("save analysis: %w", err)
	}

	nextTask := entity.Task{
		Ticker:    task.Ticker,
		Year:      task.Year,
		Period:    task.Period,
		ReportURL: task.ReportURL,
		Type:      entity.ExtractResult,
	}

	payload, err := json.Marshal(nextTask)
	if err != nil {
		return fmt.Errorf("marshal extract-result task: %w", err)
	}

	if err := u.publisher.PublishMessage(ctx, payload); err != nil {
		return fmt.Errorf("publish extract-result task: %w", err)
	}

	logger.Info("published extract-result task")
	return nil
}
