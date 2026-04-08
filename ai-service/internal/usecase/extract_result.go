package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"ai-service/internal/domain/entity"
)

type ExtractResultUsecase struct {
	ai           AIProvider
	results      ReportResultsSaver
	analysisRepo AnalysisRepository
	taskRepo     TasksRepository
}

func NewExtractResultUsecase(ai AIProvider, results ReportResultsSaver, analysisRepo AnalysisRepository, taskRepo TasksRepository) *ExtractResultUsecase {
	return &ExtractResultUsecase{
		ai:           ai,
		results:      results,
		analysisRepo: analysisRepo,
		taskRepo:     taskRepo,
	}
}

func (u *ExtractResultUsecase) Execute(ctx context.Context, task entity.Task) error {
	reportText, err := u.analysisRepo.GetAnalysis(ctx, task.Ticker, task.Year, entity.PeriodToMonths[string(task.Period)])
	if err != nil {
		return fmt.Errorf("get analysis: %w", err)
	}

	slog.Info("extracted report from database", slog.String("ticker", task.Ticker))

	text, err := u.ai.GenerateText(ctx, buildExtractPrompt(reportText), entity.Flash, GenerateParams{})
	if err != nil {
		return fmt.Errorf("generate extract result: %w", err)
	}

	var res entity.ReportResults
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		slog.Error("unable to parse extract result", slog.String("ai_response", text))
		return fmt.Errorf("unmarshal report results: %w", err)
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

	if err := u.results.SaveReportResults(ctx, &res, task.Ticker, task.Year, periodMonths); err != nil {
		return fmt.Errorf("save report results: %w", err)
	}

	if err := u.taskRepo.DeleteTask(ctx, task.Id); err != nil {
		slog.Warn("Error occured while deleting task", "err", err)
	}

	return nil
}
