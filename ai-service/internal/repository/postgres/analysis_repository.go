package postgres

import (
	"context"
	"encoding/json"
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

func (r *AnalysisRepository) GetReportResults(ctx context.Context, ticker string, year, period int) (*entity.ReportResults, error) {
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

func (r *AnalysisRepository) GetLatestReportResults(ctx context.Context, ticker string) (*entity.ReportResults, error) {
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

func (r *AnalysisRepository) GetBusinessResearch(ctx context.Context, ticker string) (*entity.BusinessResearchResult, error) {
	db := Executor(ctx, r.db)

	var profile entity.CompanyProfile
	var productsJSON, marketsJSON []byte

	err := db.QueryRow(ctx, `
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
		return nil, fmt.Errorf("get company profile: %w", err)
	}

	if err := json.Unmarshal(productsJSON, &profile.ProductsAndServices); err != nil {
		return nil, fmt.Errorf("unmarshal products: %w", err)
	}
	if err := json.Unmarshal(marketsJSON, &profile.Markets); err != nil {
		return nil, fmt.Errorf("unmarshal markets: %w", err)
	}

	rows, err := db.Query(ctx, `
		SELECT ticker, segment, share_pct, approximate, description, trend
		FROM revenue_sources WHERE ticker = $1
	`, ticker)
	if err != nil {
		return nil, fmt.Errorf("get revenue sources: %w", err)
	}
	defer rows.Close()

	var revenues []entity.RevenueSource
	for rows.Next() {
		var rs entity.RevenueSource
		if err := rows.Scan(&rs.Ticker, &rs.Segment, &rs.SharePct, &rs.Approximate, &rs.Description, &rs.Trend); err != nil {
			return nil, fmt.Errorf("scan revenue source: %w", err)
		}
		revenues = append(revenues, rs)
	}

	depRows, err := db.Query(ctx, `
		SELECT ticker, factor, type, severity, description
		FROM company_dependencies WHERE ticker = $1
	`, ticker)
	if err != nil {
		return nil, fmt.Errorf("get company dependencies: %w", err)
	}
	defer depRows.Close()

	var deps []entity.CompanyDependency
	for depRows.Next() {
		var dep entity.CompanyDependency
		if err := depRows.Scan(&dep.Ticker, &dep.Factor, &dep.Type, &dep.Severity, &dep.Description); err != nil {
			return nil, fmt.Errorf("scan company dependency: %w", err)
		}
		deps = append(deps, dep)
	}

	return &entity.BusinessResearchResult{
		Profile:      profile,
		Revenue:      revenues,
		Dependencies: deps,
	}, nil
}
