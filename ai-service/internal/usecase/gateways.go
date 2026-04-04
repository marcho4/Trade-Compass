package usecase

import (
	"context"

	"ai-service/internal/domain/entity"
)

type FinancialDataGateway interface {
	GetDailyPrices(ctx context.Context, ticker string) ([]entity.Candle, error)
	GetCBRates(ctx context.Context) (*entity.CBRate, error)
	GetMarketCap(ctx context.Context, ticker string) (float64, error)
	GetRawData(ctx context.Context, ticker string, year int, period entity.ReportPeriod) (*entity.RawData, error)
	GetRawDataHistory(ctx context.Context, ticker string) ([]entity.RawData, error)
	SaveDraft(ctx context.Context, rawData *entity.RawData) error
}

type ParserGateway interface {
	IsLatestReport(ctx context.Context, ticker string, year, periodMonths int) (bool, error)
}

type MessagePublisher interface {
	PublishMessage(ctx context.Context, value []byte) error
}
