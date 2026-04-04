package postgres

import (
	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBRepo struct {
	conn *pgxpool.Pool
}

func NewDBRepo(ctx context.Context, db *pgxpool.Pool) (*DBRepo, error) {
	if err := db.Ping(ctx); err != nil {
		slog.Error("Unable to ping db", slog.Any("error", err))
		return nil, err
	}
	slog.Info("Successfully initialized DBRepo")

	return &DBRepo{conn: db}, nil
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
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%w: analysis for %s year=%d period=%d", domain.ErrNotFound, ticker, year, period)
		}
		return "", err
	}
	return analysis, nil
}

func (d *DBRepo) GetAvailablePeriods(ctx context.Context, ticker string) ([]entity.AvailablePeriod, error) {
	query := `SELECT year, period
		FROM analysis_reports WHERE ticker = $1 ORDER BY year DESC, period DESC`
	rows, err := d.conn.Query(ctx, query, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periods []entity.AvailablePeriod
	for rows.Next() {
		var p entity.AvailablePeriod
		if err := rows.Scan(&p.Year, &p.Period); err != nil {
			return nil, err
		}
		periods = append(periods, p)
	}
	return periods, nil
}

func (d *DBRepo) SaveReportResults(ctx context.Context, result *entity.ReportResults, ticker string, year, period int) error {
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

func (d *DBRepo) GetReportResults(ctx context.Context, ticker string, year, period int) (*entity.ReportResults, error) {
	query := `
		SELECT health, growth, moat, dividends, value, total
		FROM report_results
		WHERE ticker = $1 AND year = $2 AND period = $3
	`
	var r entity.ReportResults
	err := d.conn.QueryRow(ctx, query, ticker, year, period).Scan(
		&r.Health, &r.Growth, &r.Moat, &r.Dividends, &r.Value, &r.Total,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: report results for %s year=%d period=%d", domain.ErrNotFound, ticker, year, period)
		}
		return nil, err
	}
	return &r, nil
}

func (d *DBRepo) SaveBusinessResearch(ctx context.Context, research *entity.BusinessResearchResponse) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	productsJSON, _ := json.Marshal(research.Profile.ProductsAndServices)
	marketsJSON, _ := json.Marshal(research.Profile.Markets)

	_, err = tx.Exec(ctx, `
		INSERT INTO company_profiles (ticker, company_name, description, products_and_services, markets, key_clients, business_model)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (ticker) DO UPDATE SET
			company_name = EXCLUDED.company_name,
			description = EXCLUDED.description,
			products_and_services = EXCLUDED.products_and_services,
			markets = EXCLUDED.markets,
			key_clients = EXCLUDED.key_clients,
			business_model = EXCLUDED.business_model,
			updated_at = NOW()
	`, research.Ticker, research.CompanyName, research.Profile.Description,
		productsJSON, marketsJSON, research.Profile.KeyClients, research.Profile.BusinessModel)
	if err != nil {
		return fmt.Errorf("upsert company_profiles: %w", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM revenue_sources WHERE ticker = $1`, research.Ticker)
	if err != nil {
		return fmt.Errorf("delete old revenue_sources: %w", err)
	}

	for _, rs := range research.RevenueSources {
		_, err = tx.Exec(ctx, `
			INSERT INTO revenue_sources (ticker, segment, share_pct, approximate, description, trend)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, research.Ticker, rs.Segment, rs.SharePct, rs.Approximate, rs.Description, rs.Trend)
		if err != nil {
			return fmt.Errorf("insert revenue_source: %w", err)
		}
	}

	_, err = tx.Exec(ctx, `DELETE FROM company_dependencies WHERE ticker = $1`, research.Ticker)
	if err != nil {
		return fmt.Errorf("delete old company_dependencies: %w", err)
	}

	for _, dep := range research.Dependencies {
		_, err = tx.Exec(ctx, `
			INSERT INTO company_dependencies (ticker, factor, type, severity, description)
			VALUES ($1, $2, $3, $4, $5)
		`, research.Ticker, dep.Factor, dep.Type, dep.Severity, dep.Description)
		if err != nil {
			return fmt.Errorf("insert company_dependency: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (d *DBRepo) GetBusinessResearch(ctx context.Context, ticker string) (*entity.BusinessResearchResult, error) {
	var profile entity.CompanyProfile
	var productsJSON, marketsJSON []byte

	err := d.conn.QueryRow(ctx, `
		SELECT ticker, company_name, description, products_and_services, markets, key_clients, business_model
		FROM company_profiles WHERE ticker = $1
	`, ticker).Scan(
		&profile.Ticker, &profile.CompanyName, &profile.Description,
		&productsJSON, &marketsJSON, &profile.KeyClients, &profile.BusinessModel,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: business research for %s", domain.ErrNotFound, ticker)
		}
		return nil, err
	}

	json.Unmarshal(productsJSON, &profile.ProductsAndServices)
	json.Unmarshal(marketsJSON, &profile.Markets)

	rows, err := d.conn.Query(ctx, `
		SELECT ticker, segment, share_pct, approximate, description, trend
		FROM revenue_sources WHERE ticker = $1
	`, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var revenues []entity.RevenueSource
	for rows.Next() {
		var rs entity.RevenueSource
		if err := rows.Scan(&rs.Ticker, &rs.Segment, &rs.SharePct, &rs.Approximate, &rs.Description, &rs.Trend); err != nil {
			return nil, err
		}
		revenues = append(revenues, rs)
	}

	depRows, err := d.conn.Query(ctx, `
		SELECT ticker, factor, type, severity, description
		FROM company_dependencies WHERE ticker = $1
	`, ticker)
	if err != nil {
		return nil, err
	}
	defer depRows.Close()

	var deps []entity.CompanyDependency
	for depRows.Next() {
		var dep entity.CompanyDependency
		if err := depRows.Scan(&dep.Ticker, &dep.Factor, &dep.Type, &dep.Severity, &dep.Description); err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}

	return &entity.BusinessResearchResult{
		Profile:      profile,
		Revenue:      revenues,
		Dependencies: deps,
	}, nil
}

func (d *DBRepo) SaveNews(ctx context.Context, ticker string, news *entity.NewsResponse) error {
	latestJSON, err := json.Marshal(news.LatestNews)
	if err != nil {
		return fmt.Errorf("marshal latest_news: %w", err)
	}

	importantJSON, err := json.Marshal(news.ImportantNews)
	if err != nil {
		return fmt.Errorf("marshal important_news: %w", err)
	}

	_, err = d.conn.Exec(ctx, `
		INSERT INTO company_news (ticker, latest_news, important_news)
		VALUES ($1, $2, $3)
		ON CONFLICT (ticker) DO UPDATE SET
			latest_news = EXCLUDED.latest_news,
			important_news = EXCLUDED.important_news,
			created_at = NOW()
	`, ticker, latestJSON, importantJSON)
	if err != nil {
		return fmt.Errorf("upsert company_news: %w", err)
	}

	return nil
}

func (d *DBRepo) GetFreshNews(ctx context.Context, ticker string, ttl time.Duration) (*entity.NewsResponse, error) {
	var latestJSON, importantJSON []byte
	err := d.conn.QueryRow(ctx, `
		SELECT latest_news, important_news
		FROM company_news
		WHERE ticker = $1 AND created_at > NOW() - $2::interval
	`, ticker, ttl.String()).Scan(&latestJSON, &importantJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get fresh news: %w", err)
	}

	var news entity.NewsResponse
	if err := json.Unmarshal(latestJSON, &news.LatestNews); err != nil {
		return nil, fmt.Errorf("unmarshal latest_news: %w", err)
	}
	if err := json.Unmarshal(importantJSON, &news.ImportantNews); err != nil {
		return nil, fmt.Errorf("unmarshal important_news: %w", err)
	}

	return &news, nil
}

func (d *DBRepo) SaveRiskAndGrowth(ctx context.Context, response *entity.RiskAndGrowthResponse) error {
	factorsJSON, err := json.Marshal(response.Factors)
	if err != nil {
		return fmt.Errorf("marshal factors: %w", err)
	}

	_, err = d.conn.Exec(ctx, `
		INSERT INTO risk_and_growth (ticker, factors)
		VALUES ($1, $2)
		ON CONFLICT (ticker) DO UPDATE SET
			factors = EXCLUDED.factors,
			created_at = NOW()
	`, response.Ticker, factorsJSON)
	if err != nil {
		return fmt.Errorf("upsert risk_and_growth: %w", err)
	}

	return nil
}

func (d *DBRepo) GetRiskAndGrowth(ctx context.Context, ticker string, ttl time.Duration) (*entity.RiskAndGrowthResponse, error) {
	var factorsJSON []byte
	err := d.conn.QueryRow(ctx, `
		SELECT factors
		FROM risk_and_growth
		WHERE ticker = $1 AND created_at > NOW() - $2::interval
	`, ticker, ttl.String()).Scan(&factorsJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get risk and growth: %w", err)
	}

	var factors []entity.RiskAndGrowthFactor
	if err := json.Unmarshal(factorsJSON, &factors); err != nil {
		return nil, fmt.Errorf("unmarshal factors: %w", err)
	}

	return &entity.RiskAndGrowthResponse{
		Ticker:  ticker,
		Factors: factors,
	}, nil
}

func (d *DBRepo) GetLatestReportResults(ctx context.Context, ticker string) (*entity.ReportResults, error) {
	query := `
		SELECT health, growth, moat, dividends, value, total
		FROM report_results
		WHERE ticker = $1
		ORDER BY year DESC, period DESC
		LIMIT 1
	`
	var r entity.ReportResults
	err := d.conn.QueryRow(ctx, query, ticker).Scan(
		&r.Health, &r.Growth, &r.Moat, &r.Dividends, &r.Value, &r.Total,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: latest report results for %s", domain.ErrNotFound, ticker)
		}
		return nil, err
	}
	return &r, nil
}
