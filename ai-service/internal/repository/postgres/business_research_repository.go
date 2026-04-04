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

type BusinessResearchRepository struct {
	db *pgxpool.Pool
}

func NewBusinessResearchRepository(db *pgxpool.Pool) *BusinessResearchRepository {
	return &BusinessResearchRepository{db: db}
}

func (r *BusinessResearchRepository) SaveBusinessResearch(ctx context.Context, research *entity.BusinessResearchResponse) error {
	db := Executor(ctx, r.db)

	productsJSON, err := json.Marshal(research.Profile.ProductsAndServices)
	if err != nil {
		return fmt.Errorf("marshal products: %w", err)
	}

	marketsJSON, err := json.Marshal(research.Profile.Markets)
	if err != nil {
		return fmt.Errorf("marshal markets: %w", err)
	}

	_, err = db.Exec(ctx, `
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

	_, err = db.Exec(ctx, `DELETE FROM revenue_sources WHERE ticker = $1`, research.Ticker)
	if err != nil {
		return fmt.Errorf("delete old revenue_sources: %w", err)
	}

	for _, rs := range research.RevenueSources {
		_, err = db.Exec(ctx, `
			INSERT INTO revenue_sources (ticker, segment, share_pct, approximate, description, trend)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, research.Ticker, rs.Segment, rs.SharePct, rs.Approximate, rs.Description, rs.Trend)
		if err != nil {
			return fmt.Errorf("insert revenue_source: %w", err)
		}
	}

	_, err = db.Exec(ctx, `DELETE FROM company_dependencies WHERE ticker = $1`, research.Ticker)
	if err != nil {
		return fmt.Errorf("delete old company_dependencies: %w", err)
	}

	for _, dep := range research.Dependencies {
		_, err = db.Exec(ctx, `
			INSERT INTO company_dependencies (ticker, factor, type, severity, description)
			VALUES ($1, $2, $3, $4, $5)
		`, research.Ticker, dep.Factor, dep.Type, dep.Severity, dep.Description)
		if err != nil {
			return fmt.Errorf("insert company_dependency: %w", err)
		}
	}

	return nil
}

func (r *BusinessResearchRepository) GetBusinessResearch(ctx context.Context, ticker string) (*entity.BusinessResearchResult, error) {
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
