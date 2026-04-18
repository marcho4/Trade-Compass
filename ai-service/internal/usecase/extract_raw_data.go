package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"time"

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
		slog.String("id", task.Id),
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
		prompt := docs.RawDataAgentPrompt() + "\n<ticker>" + task.Ticker + "</ticker>"

		pdfBytes, err := u.storage.DownloadPDF(ctx, task.ReportURL)
		if err != nil {
			return fmt.Errorf("download PDF: %w", err)
		}

		text, err := u.ai.AnalyzeWithPDF(ctx, pdfBytes, prompt, entity.Pro)
		if err != nil {
			return fmt.Errorf("ai call with pdf: %w", err)
		}

		var rawData entity.RawData
		if err := json.Unmarshal([]byte(text), &rawData); err != nil {
			return fmt.Errorf("unmarshal raw data: %w", err)
		}

		stockInfo, err := u.fd.GetStockInfo(ctx, task.Ticker)
		if err != nil {
			logger.Warn("failed to get stock info, skipping market cap calculation", slog.Any("error", err))
		} else {
			periodEnd := periodEndDate(task.Year, period)
			price, err := u.fd.GetPriceAt(ctx, task.Ticker, periodEnd)
			if err != nil {
				logger.Warn("failed to get price at period end, skipping market cap calculation",
					slog.String("date", periodEnd.Format("2006-01-02")),
					slog.Any("error", err),
				)
			} else {
				unitDivisor := unitDivisorForReportUnits(rawData.ReportUnits)
				shares := int64(stockInfo.NumberOfShares)
				rawData.SharesOutstanding = &shares
				marketCap := int64(math.Round(price*float64(stockInfo.NumberOfShares))) / unitDivisor
				rawData.MarketCap = &marketCap
			}
		}

		rawData.ComputeDerivedFields()
		rawData.Status = entity.RawDataStatusConfirmed

		if err := u.fd.SaveDraft(ctx, &rawData); err != nil {
			return fmt.Errorf("save raw data: %w", err)
		}
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

	logger.Info("raw data extraction succeed")

	return nil
}

func unitDivisorForReportUnits(units string) int64 {
	switch units {
	case "thousands":
		return 1_000
	case "millions":
		return 1_000_000
	case "billions":
		return 1_000_000_000
	default:
		return 1
	}
}

func periodEndDate(year int, period entity.ReportPeriod) time.Time {
	months, ok := entity.PeriodToMonths[string(period)]
	if !ok {
		months = 12
	}
	lastDay := time.Date(year, time.Month(months)+1, 0, 0, 0, 0, 0, time.UTC)
	return lastDay
}
