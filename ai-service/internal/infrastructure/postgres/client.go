package postgres

import (
	"ai-service/internal/domain"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBRepo struct {
	conn *pgxpool.Pool
}

func NewDBRepo(ctx context.Context, postgresUrl string) (*DBRepo, error) {
	conn, err := pgxpool.New(ctx, postgresUrl)
	if err != nil {
		slog.Error("Error while creating DB Repo")
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		slog.Error("Unable to ping db", slog.Any("error", err))
		return nil, err
	}

	slog.Info("Successfully initialized DBRepo")

	return &DBRepo{conn: conn}, nil
}

func (d *DBRepo) Close() {
	if d.conn != nil {
		d.conn.Close()
	}
}

func (d *DBRepo) SaveAnalysis(ctx context.Context, result, ticker string, year, period int) error {
	slog.Info("Executing upsert",
		slog.String("ticker", ticker),
		slog.Int("year", year),
		slog.Int("period", period),
		slog.Int("result_length", len(result)),
	)

	query := `
		INSERT INTO analysis_reports (ticker, year, period, analysis)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (ticker, year, period)
		DO UPDATE SET analysis = EXCLUDED.analysis
	`
	_, err := d.conn.Exec(ctx, query, ticker, year, period, result)
	if err != nil {
		slog.Error("[SaveAnalysis] exec failed",
			slog.String("ticker", ticker),
			slog.Int("year", year),
			slog.Int("period", period),
			slog.Any("error", err),
		)
		return err
	}

	slog.Info("Successfully saved analysis")
	return nil
}

func (d *DBRepo) GetAnalysis(ctx context.Context, ticker string, year, period int) (string, error) {
	query := `SELECT analysis FROM analysis_reports WHERE ticker = $1 AND year = $2 AND period = $3`
	var analysis string
	err := d.conn.QueryRow(ctx, query, ticker, year, period).Scan(&analysis)
	if err != nil {
		return "", err
	}
	return analysis, nil
}

type AvailablePeriod struct {
	Year   int `json:"year"`
	Period int `json:"period"`
}

func (d *DBRepo) GetAvailablePeriods(ctx context.Context, ticker string) ([]AvailablePeriod, error) {
	query := `SELECT year, period
		FROM analysis_reports WHERE ticker = $1 ORDER BY year DESC, period DESC`
	rows, err := d.conn.Query(ctx, query, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periods []AvailablePeriod
	for rows.Next() {
		var p AvailablePeriod
		if err := rows.Scan(&p.Year, &p.Period); err != nil {
			return nil, err
		}
		periods = append(periods, p)
	}
	return periods, nil
}

func (d *DBRepo) SaveReportResults(ctx context.Context, result *domain.ReportResults, ticker string, year, period int) error {
	query := `
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
	_, err := d.conn.Exec(ctx, query,
		ticker, year, period,
		result.Health, result.Growth, result.Moat, result.Dividends, result.Value, result.Total,
	)
	if err != nil {
		slog.Error("SaveReportResults exec failed",
			slog.String("ticker", ticker),
			slog.Int("year", year),
			slog.Int("period", period),
			slog.Any("error", err),
		)
		return err
	}

	slog.Info("SaveReportResults success",
		slog.String("ticker", ticker),
		slog.Int("year", year),
		slog.Int("period", period),
	)
	return nil
}

func (d *DBRepo) GetReportResults(ctx context.Context, ticker string, year, period int) (*domain.ReportResults, error) {
	query := `
		SELECT health, growth, moat, dividends, value, total
		FROM report_results
		WHERE ticker = $1 AND year = $2 AND period = $3
	`
	var r domain.ReportResults
	err := d.conn.QueryRow(ctx, query, ticker, year, period).Scan(
		&r.Health, &r.Growth, &r.Moat, &r.Dividends, &r.Value, &r.Total,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (d *DBRepo) GetLatestReportResults(ctx context.Context, ticker string) (*domain.ReportResults, error) {
	query := `
		SELECT health, growth, moat, dividends, value, total
		FROM report_results
		WHERE ticker = $1
		ORDER BY year DESC, period DESC
		LIMIT 1
	`
	var r domain.ReportResults
	err := d.conn.QueryRow(ctx, query, ticker).Scan(
		&r.Health, &r.Growth, &r.Moat, &r.Dividends, &r.Value, &r.Total,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
