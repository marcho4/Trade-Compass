package application

import (
	"context"
	"financial_data/internal/domain"
	"time"
)

type DividendsRepository interface {
	GetByTicker(ctx context.Context, ticker string) ([]domain.Dividends, error)
	GetByID(ctx context.Context, id int) (*domain.Dividends, error)
	Create(ctx context.Context, dividend *domain.Dividends) error
	Update(ctx context.Context, id int, dividend *domain.Dividends) error
	Delete(ctx context.Context, id int) error
}

type MacroDataRepository interface {
	GetCurrent(ctx context.Context) (*domain.CBRate, error)
	GetByDate(ctx context.Context, date time.Time) (*domain.CBRate, error)
	GetHistory(ctx context.Context, from, to time.Time) ([]domain.CBRate, error)
	Create(ctx context.Context, rate *domain.CBRate) error
	Update(ctx context.Context, date time.Time, rate float64) error
	Delete(ctx context.Context, date time.Time) error
}

type NewsRepository interface {
	GetByID(ctx context.Context, id int) (*domain.News, error)
	GetByTicker(ctx context.Context, ticker string) ([]domain.News, error)
	GetBySector(ctx context.Context, sectorID int) ([]domain.News, error)
	Create(ctx context.Context, news *domain.News) error
	Update(ctx context.Context, id int, news *domain.News) error
	Delete(ctx context.Context, id int) error
}

type RawDataRepository interface {
	GetByTickerAndPeriod(ctx context.Context, ticker string, year int, period domain.ReportPeriod) (*domain.RawData, error)
	GetLatestByTicker(ctx context.Context, ticker string) (*domain.RawData, error)
	GetHistoryByTicker(ctx context.Context, ticker string) ([]domain.RawData, error)
	GetDraftByTickerAndPeriod(ctx context.Context, ticker string, year int, period domain.ReportPeriod) (*domain.RawData, error)
	GetDraftsByTicker(ctx context.Context, ticker string) ([]domain.RawData, error)
	ConfirmDraft(ctx context.Context, ticker string, year int, period domain.ReportPeriod) error
	Create(ctx context.Context, rawData *domain.RawData) error
	Update(ctx context.Context, rawData *domain.RawData) error
	Delete(ctx context.Context, ticker string, year int, period domain.ReportPeriod) error
}

type CompanyRepository interface {
	GetByTicker(ctx context.Context, ticker string) (*domain.Company, error)
	GetAll(ctx context.Context) ([]domain.Company, error)
	GetBySector(ctx context.Context, sectorID int) ([]domain.Company, error)
	Create(ctx context.Context, company *domain.Company) error
	Update(ctx context.Context, ticker string, company *domain.Company) error
	Delete(ctx context.Context, ticker string) error
}

type SectorRepository interface {
	GetByID(ctx context.Context, id int) (*domain.SectorModel, error)
	GetAll(ctx context.Context) ([]domain.SectorModel, error)
	Create(ctx context.Context, sector *domain.SectorModel) error
	Update(ctx context.Context, id int, sector *domain.SectorModel) error
	Delete(ctx context.Context, id int) error
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
