package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"financial_data/internal/domain"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	companyCacheTTL = 24 * time.Hour
	companiesAllKey = "companies:all"
)

type CompanyRepository struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

func NewCompanyRepository(pool *pgxpool.Pool, redisClient *redis.Client) *CompanyRepository {
	return &CompanyRepository{pool: pool, redis: redisClient}
}

func (r *CompanyRepository) companyKey(ticker string) string {
	return fmt.Sprintf("company:%s", ticker)
}

func (r *CompanyRepository) sectorKey(sectorID int) string {
	return fmt.Sprintf("companies:sector:%d", sectorID)
}

func (r *CompanyRepository) invalidateListCaches(ctx context.Context, sectorID int) {
	keys := []string{companiesAllKey}
	if sectorID > 0 {
		keys = append(keys, r.sectorKey(sectorID))
	}
	if err := r.redis.Del(ctx, keys...).Err(); err != nil {
		slog.Warn("redis del list caches failed", "error", err)
	}
}

func (r *CompanyRepository) GetByTicker(ctx context.Context, ticker string) (*domain.Company, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	key := r.companyKey(ticker)
	cached, err := r.redis.Get(ctx, key).Result()
	if err == nil {
		var company domain.Company
		if jsonErr := json.Unmarshal([]byte(cached), &company); jsonErr == nil {
			return &company, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		slog.Warn("redis get failed", "key", key, "error", err)
	}

	query := `
		SELECT id, ticker, name, sector_id, lot_size, ceo
		FROM companies
		WHERE ticker = $1
	`

	var name *string
	company := &domain.Company{}
	err = r.pool.QueryRow(ctx, query, ticker).Scan(
		&company.ID, &company.Ticker, &name, &company.SectorID, &company.LotSize, &company.CEO,
	)
	if name != nil {
		company.Name = *name
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("company not found for ticker %s: %w", ticker, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if data, jsonErr := json.Marshal(company); jsonErr == nil {
		if setErr := r.redis.Set(ctx, key, data, companyCacheTTL).Err(); setErr != nil {
			slog.Warn("redis set failed", "key", key, "error", setErr)
		}
	}

	return company, nil
}

func (r *CompanyRepository) GetAll(ctx context.Context) ([]domain.Company, error) {
	query := `
		SELECT id, ticker, name, sector_id, lot_size, ceo
		FROM companies
		ORDER BY ticker
	`

	companies := make([]domain.Company, 0)
	cached, err := r.redis.Get(ctx, companiesAllKey).Result()
	if err == nil {
		if jsonErr := json.Unmarshal([]byte(cached), &companies); jsonErr == nil {
			return companies, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		slog.Warn("redis get failed", "key", companiesAllKey, "error", err)
	}

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query companies: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var company domain.Company
		var name *string
		err := rows.Scan(
			&company.ID, &company.Ticker, &name, &company.SectorID, &company.LotSize, &company.CEO,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}
		if name != nil {
			company.Name = *name
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating companies: %w", err)
	}

	if data, jsonErr := json.Marshal(companies); jsonErr == nil {
		if setErr := r.redis.Set(ctx, companiesAllKey, data, companyCacheTTL).Err(); setErr != nil {
			slog.Warn("redis set failed", "key", companiesAllKey, "error", setErr)
		}
	}

	return companies, nil
}

func (r *CompanyRepository) GetBySector(ctx context.Context, sectorID int) ([]domain.Company, error) {
	key := r.sectorKey(sectorID)
	cached, err := r.redis.Get(ctx, key).Result()
	if err == nil {
		var companies []domain.Company
		if jsonErr := json.Unmarshal([]byte(cached), &companies); jsonErr == nil {
			return companies, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		slog.Warn("redis get failed", "key", key, "error", err)
	}

	query := `
		SELECT id, ticker, name, sector_id, lot_size, ceo
		FROM companies
		WHERE sector_id = $1
		ORDER BY ticker
	`

	rows, err := r.pool.Query(ctx, query, sectorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query companies by sector: %w", err)
	}
	defer rows.Close()

	var companies []domain.Company
	for rows.Next() {
		var company domain.Company
		var name *string
		err := rows.Scan(
			&company.ID, &company.Ticker, &name, &company.SectorID, &company.LotSize, &company.CEO,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}
		if name != nil {
			company.Name = *name
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating companies: %w", err)
	}

	if data, jsonErr := json.Marshal(companies); jsonErr == nil {
		if setErr := r.redis.Set(ctx, key, data, companyCacheTTL).Err(); setErr != nil {
			slog.Warn("redis set failed", "key", key, "error", setErr)
		}
	}

	return companies, nil
}

func (r *CompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	if company == nil {
		return fmt.Errorf("company is nil: %w", domain.ErrInvalidInput)
	}
	if company.Ticker == "" {
		return fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := `
		INSERT INTO companies (ticker, name, sector_id, lot_size, ceo)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		company.Ticker, company.Name, company.SectorID, company.LotSize, company.CEO,
	).Scan(&company.ID)

	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}

	r.invalidateListCaches(ctx, company.SectorID)

	return nil
}

func (r *CompanyRepository) Update(ctx context.Context, ticker string, company *domain.Company) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}
	if company == nil {
		return fmt.Errorf("company is nil: %w", domain.ErrInvalidInput)
	}

	query := `
		UPDATE companies SET
			name = $2, sector_id = $3, lot_size = $4, ceo = $5
		WHERE ticker = $1
	`

	result, err := r.pool.Exec(ctx, query,
		ticker, company.Name, company.SectorID, company.LotSize, company.CEO,
	)

	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("company not found for ticker %s: %w", ticker, domain.ErrNotFound)
	}

	if err := r.redis.Del(ctx, r.companyKey(ticker)).Err(); err != nil {
		slog.Warn("redis del failed on update", "ticker", ticker, "error", err)
	}
	r.invalidateListCaches(ctx, company.SectorID)

	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, ticker string) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := `DELETE FROM companies WHERE ticker = $1 RETURNING sector_id`

	var sectorID int
	err := r.pool.QueryRow(ctx, query, ticker).Scan(&sectorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("company not found for ticker %s: %w", ticker, domain.ErrNotFound)
		}
		return fmt.Errorf("failed to delete company: %w", err)
	}

	if err := r.redis.Del(ctx, r.companyKey(ticker)).Err(); err != nil {
		slog.Warn("redis del failed on delete", "ticker", ticker, "error", err)
	}
	r.invalidateListCaches(ctx, sectorID)

	return nil
}
