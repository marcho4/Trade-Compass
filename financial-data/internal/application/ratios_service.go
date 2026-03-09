package application

import (
	"context"
	"errors"
	"financial_data/internal/application/routers"
	"financial_data/internal/domain"
	"fmt"
	"log/slog"
)

type RatiosService struct {
	rawDataRepo routers.RawDataRepository
	ratiosRepo  routers.RatiosRepository
	companyRepo routers.CompanyRepository
}

func NewRatiosService(rawDataRepo routers.RawDataRepository, ratiosRepo routers.RatiosRepository, companyRepo routers.CompanyRepository) *RatiosService {
	return &RatiosService{
		rawDataRepo: rawDataRepo,
		ratiosRepo:  ratiosRepo,
		companyRepo: companyRepo,
	}
}

func (s *RatiosService) CalculateAndSave(ctx context.Context, rawData *domain.RawData) error {
	company, err := s.companyRepo.GetByTicker(ctx, rawData.Ticker)
	if err != nil {
		slog.Error("ratios: failed to get company", "ticker", rawData.Ticker, "error", err)
		return err
	}

	previous, err := s.rawDataRepo.GetByTickerAndPeriod(ctx, rawData.Ticker, rawData.Year-1, rawData.Period)
	if err != nil {
		previous = nil
	}

	ratios := domain.CalculateRatios(rawData, previous)
	ratios.Ticker = rawData.Ticker
	ratios.Year = rawData.Year
	ratios.Period = rawData.Period

	err = s.ratiosRepo.Update(ctx, ratios)
	if errors.Is(err, domain.ErrNotFound) {
		err = s.ratiosRepo.Create(ctx, domain.Sector(company.SectorID), ratios)
	}

	if err != nil {
		slog.Error("ratios: failed to save ratios", "ticker", rawData.Ticker, "error", err)
		return err
	}

	slog.Info("ratios: calculated and saved", "ticker", rawData.Ticker, "year", rawData.Year, "period", rawData.Period)
	return nil
}

func (s *RatiosService) RecalculateAll(ctx context.Context, ticker string) error {
	history, err := s.rawDataRepo.GetHistoryByTicker(ctx, ticker)
	if err != nil {
		return fmt.Errorf("ratios: failed to get history for %s: %w", ticker, err)
	}

	var calculated int
	for i := range history {
		if err := s.CalculateAndSave(ctx, &history[i]); err != nil {
			slog.Error("ratios: failed to recalculate", "ticker", ticker, "year", history[i].Year, "period", history[i].Period, "error", err)
			continue
		}
		calculated++
	}

	slog.Info("ratios: recalculated all", "ticker", ticker, "total", len(history), "calculated", calculated)
	return nil
}
