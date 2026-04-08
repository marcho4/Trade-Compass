package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	docs "ai-service/internal/docs"
	"ai-service/internal/domain/entity"
)

type ExtractRawDataUsecase struct {
	ai        AIProvider
	fd        FinancialDataGateway
	parser    ParserGateway
	publisher MessagePublisher
	storage   StorageClient
}

func NewExtractRawDataUsecase(
	ai AIProvider,
	fd FinancialDataGateway,
	parser ParserGateway,
	publisher MessagePublisher,
	storage StorageClient,
) *ExtractRawDataUsecase {
	return &ExtractRawDataUsecase{
		ai:        ai,
		fd:        fd,
		parser:    parser,
		publisher: publisher,
		storage:   storage,
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

		prompt := docs.RawDataAgentPrompt() + "\n<ticker>" + task.Ticker + "</ticker>"

		logger.Info("downloading PDF for raw data extraction", slog.String("report_url", task.ReportURL))

		pdfBytes, err := u.storage.DownloadPDF(ctx, task.ReportURL)
		if err != nil {
			return fmt.Errorf("download PDF: %w", err)
		}
		logger.Info("PDF downloaded", slog.Int("pdf_size_bytes", len(pdfBytes)))

		text, err := u.ai.AnalyzeWithPDF(ctx, pdfBytes, prompt, entity.Flash)
		if err != nil {
			return fmt.Errorf("ai call with pdf: %w", err)
		}

		var rawData entity.RawData
		if err := json.Unmarshal([]byte(text), &rawData); err != nil {
			return fmt.Errorf("unmarshal raw data: %w", err)
		}

		rawData.ComputeDerivedFields()
		rawData.Status = entity.RawDataStatusConfirmed

		if err := u.fd.SaveDraft(ctx, &rawData); err != nil {
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
