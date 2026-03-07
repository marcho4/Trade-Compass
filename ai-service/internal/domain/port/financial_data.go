package port

import (
	"ai-service/internal/domain/entity"
	"context"
)

type FinancialDataProvider interface {
	GetDailyPrices(ctx context.Context, ticker string) ([]entity.Candle, error)
	GetCBRates(ctx context.Context) (*entity.CBRate, error)
	GetMarketCap(ctx context.Context, ticker string) (float64, error)
	GetRawData(ctx context.Context, ticker string, year int, period entity.ReportPeriod) (*entity.RawData, error)
	GetRawDataHistory(ctx context.Context, ticker string) ([]entity.RawData, error)
	SaveDraft(ctx context.Context, rawData *entity.RawData) error
	GetDraft(ctx context.Context, ticker string, year int, period entity.ReportPeriod) (*entity.RawData, error)
}
