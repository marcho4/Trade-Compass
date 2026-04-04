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

type AnalysisRepository struct {
	db *pgxpool.Pool
}

func NewAnalysisRepository(db *pgxpool.Pool) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) GetAnalysis(ctx context.Context, ticker string, year, period int) (string, error) {
	db := Executor(ctx, r.db)

	sql := `SELECT analysis FROM analysis_reports WHERE ticker = $1 AND year = $2 AND period = $3`

	var analysis string
	err := db.QueryRow(ctx, sql, ticker, year, period).Scan(&analysis)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%w: analysis for %s year=%d period=%d", domain.ErrNotFound, ticker, year, period)
		}
		return "", fmt.Errorf("get analysis: %w", err)
	}

	return analysis, nil
}

func (r *AnalysisRepository) GetAvailablePeriods(ctx context.Context, ticker string) ([]entity.AvailablePeriod, error) {
	db := Executor(ctx, r.db)

	sql := `SELECT year, period FROM analysis_reports WHERE ticker = $1 ORDER BY year DESC, period DESC`

	rows, err := db.Query(ctx, sql, ticker)
	if err != nil {
		return nil, fmt.Errorf("get available periods: %w", err)
	}
	defer rows.Close()

	var periods []entity.AvailablePeriod
	for rows.Next() {
		var p entity.AvailablePeriod
		if err := rows.Scan(&p.Year, &p.Period); err != nil {
			return nil, fmt.Errorf("scan period: %w", err)
		}
		periods = append(periods, p)
	}

	return periods, nil
}

func (r *AnalysisRepository) SaveAnalysis(ctx context.Context, result, ticker string, year, period int) error {
	db := Executor(ctx, r.db)

	sql := `
		INSERT INTO analysis_reports (ticker, year, period, analysis)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (ticker, year, period)
		DO UPDATE SET analysis = EXCLUDED.analysis
	`

	_, err := db.Exec(ctx, sql, ticker, year, period, result)
	if err != nil {
		return fmt.Errorf("save analysis: %w", err)
	}

	return nil
}
