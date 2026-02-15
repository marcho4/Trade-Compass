package postgres

import (
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

func (d *DBRepo) SaveAnalysis(result, ticker string, year, period int) error {
	slog.Info("[SaveAnalysis] executing upsert",
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
	res, err := d.conn.Exec(context.Background(), query, ticker, year, period, result)
	if err != nil {
		slog.Error("[SaveAnalysis] exec failed",
			slog.String("ticker", ticker),
			slog.Int("year", year),
			slog.Int("period", period),
			slog.Any("error", err),
		)
		return err
	}

	slog.Info("[SaveAnalysis] success",
		slog.String("ticker", ticker),
		slog.Int("year", year),
		slog.Int("period", period),
		slog.Int64("rows_affected", res.RowsAffected()),
	)
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
