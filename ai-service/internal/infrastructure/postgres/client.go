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

type AnalysisReport struct {
	ID        int64  `json:"id"`
	Ticker    string `json:"ticker"`
	Year      int    `json:"year"`
	Period    int    `json:"period"`
	Analysis  string `json:"analysis"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (d *DBRepo) GetAnalysesByTicker(ctx context.Context, ticker string) ([]AnalysisReport, error) {
	query := `SELECT id, ticker, year, period, analysis, created_at, updated_at
		FROM analysis_reports WHERE ticker = $1 ORDER BY year DESC, period DESC`
	rows, err := d.conn.Query(ctx, query, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []AnalysisReport
	for rows.Next() {
		var r AnalysisReport
		if err := rows.Scan(&r.ID, &r.Ticker, &r.Year, &r.Period, &r.Analysis, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}
