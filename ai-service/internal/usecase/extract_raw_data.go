package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"ai-service/internal/domain/entity"
)

type ExtractRawDataUsecase struct {
	ai        AIService
	fd        FinancialDataGateway
	parser    ParserGateway
	publisher MessagePublisher
}

func NewExtractRawDataUsecase(
	ai AIService,
	fd FinancialDataGateway,
	parser ParserGateway,
	publisher MessagePublisher,
) *ExtractRawDataUsecase {
	return &ExtractRawDataUsecase{
		ai:        ai,
		fd:        fd,
		parser:    parser,
		publisher: publisher,
	}
}

func (u *ExtractRawDataUsecase) Execute(ctx context.Context, task entity.Task) error {
	logger := slog.With(
		slog.String("ticker", task.Ticker),
		slog.Int("year", task.Year),
		slog.String("period", task.Period),
	)

	period := entity.ReportPeriod(task.Period)

	logger.Info("extracting raw data")

	existing, err := u.fd.GetRawData(ctx, task.Ticker, task.Year, period)
	if err != nil {
		return fmt.Errorf("check existing raw data: %w", err)
	}

	if existing == nil {
		logger.Info("raw data not found, extracting")

		result, err := u.ai.ExtractRawData(ctx, task.Ticker, task.ReportURL, task.Year, period)
		if err != nil {
			return fmt.Errorf("extract raw data: %w", err)
		}
		result.Status = entity.RawDataStatusConfirmed

		if err := u.fd.SaveDraft(ctx, result); err != nil {
			return fmt.Errorf("save raw data: %w", err)
		}

		logger.Info("raw data saved")
	}

	nextTask := entity.Task{
		Id:        task.Id,
		Ticker:    task.Ticker,
		Year:      task.Year,
		Period:    task.Period,
		ReportURL: task.ReportURL,
		Type:      entity.RawDataSuccess,
	}

	payload, err := json.Marshal(nextTask)
	if err != nil {
		return fmt.Errorf("marshal analyze task: %w", err)
	}

	if err := u.publisher.PublishMessage(ctx, payload); err != nil {
		return fmt.Errorf("publish analyze task: %w", err)
	}

	logger.Info("report is the latest, published analyze task")

	return nil
}
