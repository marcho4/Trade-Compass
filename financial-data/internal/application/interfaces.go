package application

import (
	"context"
	"financial_data/internal/domain"
)

type DividendsRepository interface {
	GetDividends(ticker string) ([]domain.Dividends, error)
	GetDividendPolicy(ticker string) (string, error)
}

type MacroDataRepository interface {
	GetCurrentCentralBankRate() (float64, error)
	GetCentralBankRateHistory() (float64, error)
	GetCentralBankRateForecast() (float64, error)
	GetCurrentRubleDollarRate() (float64, error)
	GetRubleDollarForecast() (float64, error)
	GetRubleDollarHistory() (float64, error)
}

type NewsRepository interface {
	GetTickerNews(ticker string) ([]domain.News, error)
	GetSectorNews(sector domain.Sector) ([]domain.News, error)
}

type RawDataRepository interface {
	GetRawData(ticker string) (*domain.RawData, error)
}

type PriceRepository interface {
	GetStockPrice(ticker string, daysBackwards int, interval domain.Period) ([]domain.Candle, error)
}

type RatiosRepository interface {
	GetByTicker(ctx context.Context, ticker string) (*domain.Ratios, error)
	GetBySector(ctx context.Context, sector domain.Sector) (*domain.Ratios, error)
	Create(ctx context.Context, ticker string, sector domain.Sector, ratios *domain.Ratios) error
	Update(ctx context.Context, ticker string, ratios *domain.Ratios) error
	Delete(ctx context.Context, ticker string) error
}
