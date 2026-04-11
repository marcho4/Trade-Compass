package postgres

import (
	"context"
	"errors"
	"fmt"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportResultsRepository struct {
	db *pgxpool.Pool
}

func NewReportResultsRepository(db *pgxpool.Pool) *ReportResultsRepository {
	return &ReportResultsRepository{db: db}
}

func (r *ReportResultsRepository) SaveReportResults(ctx context.Context, result *entity.ReportResults, ticker string, year, period int) error {
	db := Executor(ctx, r.db)

	sql := `
		INSERT INTO report_results (ticker, year, period, health, growth, moat, dividends, value, total)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (ticker, year, period)
		DO UPDATE SET
			health = EXCLUDED.health,
			growth = EXCLUDED.growth,
			moat = EXCLUDED.moat,
			dividends = EXCLUDED.dividends,
			value = EXCLUDED.value,
			total = EXCLUDED.total
	`

	_, err := db.Exec(ctx, sql,
		ticker, year, period,
		result.Health, result.Growth, result.Moat, result.Dividends, result.Value, result.Total,
	)
	if err != nil {
		return fmt.Errorf("save report results: %w", err)
	}

	return nil
}

func (r *ReportResultsRepository) GetReportResults(ctx context.Context, ticker string, year, period int) (*entity.ReportResults, error) {
	db := Executor(ctx, r.db)

	sql := `
		SELECT health, growth, moat, dividends, value, total
		FROM report_results
		WHERE ticker = $1 AND year = $2 AND period = $3
	`

	var res entity.ReportResults
	err := db.QueryRow(ctx, sql, ticker, year, period).Scan(
		&res.Health, &res.Growth, &res.Moat, &res.Dividends, &res.Value, &res.Total,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: report results for %s year=%d period=%d", domain.ErrNotFound, ticker, year, period)
		}
		return nil, fmt.Errorf("get report results: %w", err)
	}

	return &res, nil
}

func (r *ReportResultsRepository) GetLatestReportResults(ctx context.Context, ticker string) (*entity.ReportResults, error) {
	db := Executor(ctx, r.db)

	sql := `
		SELECT health, growth, moat, dividends, value, total
		FROM report_results
		WHERE ticker = $1
		ORDER BY year DESC, period DESC
		LIMIT 1
	`

	var res entity.ReportResults
	err := db.QueryRow(ctx, sql, ticker).Scan(
		&res.Health, &res.Growth, &res.Moat, &res.Dividends, &res.Value, &res.Total,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: latest report results for %s", domain.ErrNotFound, ticker)
		}
		return nil, fmt.Errorf("get latest report results: %w", err)
	}

	return &res, nil
}
