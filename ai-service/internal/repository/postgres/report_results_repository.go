package postgres

import (
	"context"
	"fmt"

	"ai-service/internal/domain/entity"

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
