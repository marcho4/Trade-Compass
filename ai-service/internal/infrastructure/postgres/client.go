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
	query := `
		INSERT INTO analysis_reports (ticker, year, period, analysis)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (ticker, year, period)
		DO UPDATE SET analysis = EXCLUDED.analysis
	`
	_, err := d.conn.Exec(context.Background(), query, ticker, year, period, result)
	return err
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
