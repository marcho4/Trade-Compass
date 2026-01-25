package infrastructure

import (
	"context"
	"errors"
	"financial_data/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompanyRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{pool: pool}
}

func (r *CompanyRepository) GetByTicker(ctx context.Context, ticker string) (*domain.Company, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty")
	}

	query := `
		SELECT id, inn, ticker, owner, sector_id, lot_size, ceo, employees
		FROM companies
		WHERE ticker = $1
	`

	company := &domain.Company{}
	err := r.pool.QueryRow(ctx, query, ticker).Scan(
		&company.ID, &company.INN, &company.Ticker, &company.Owner,
		&company.SectorID, &company.LotSize, &company.CEO, &company.Employees,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("company not found for ticker %s", ticker)
		}
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	return company, nil
}

func (r *CompanyRepository) GetAll(ctx context.Context) ([]domain.Company, error) {
	query := `
		SELECT id, inn, ticker, owner, sector_id, lot_size, ceo, employees
		FROM companies
		ORDER BY ticker
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query companies: %w", err)
	}
	defer rows.Close()

	var companies []domain.Company
	for rows.Next() {
		var company domain.Company
		err := rows.Scan(
			&company.ID, &company.INN, &company.Ticker, &company.Owner,
			&company.SectorID, &company.LotSize, &company.CEO, &company.Employees,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating companies: %w", err)
	}

	return companies, nil
}

func (r *CompanyRepository) GetBySector(ctx context.Context, sectorID int) ([]domain.Company, error) {
	query := `
		SELECT id, inn, ticker, owner, sector_id, lot_size, ceo, employees
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
		err := rows.Scan(
			&company.ID, &company.INN, &company.Ticker, &company.Owner,
			&company.SectorID, &company.LotSize, &company.CEO, &company.Employees,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating companies: %w", err)
	}

	return companies, nil
}

func (r *CompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	if company == nil {
		return fmt.Errorf("company is nil")
	}
	if company.INN == "" {
		return fmt.Errorf("INN is empty")
	}
	if company.Ticker == "" {
		return fmt.Errorf("ticker is empty")
	}

	query := `
		INSERT INTO companies (inn, ticker, owner, sector_id, lot_size, ceo, employees)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		company.INN, company.Ticker, company.Owner, company.SectorID,
		company.LotSize, company.CEO, company.Employees,
	).Scan(&company.ID)

	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}

	return nil
}

func (r *CompanyRepository) Update(ctx context.Context, ticker string, company *domain.Company) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty")
	}
	if company == nil {
		return fmt.Errorf("company is nil")
	}

	query := `
		UPDATE companies SET
			inn = $2, owner = $3, sector_id = $4, lot_size = $5, ceo = $6, employees = $7
		WHERE ticker = $1
	`

	result, err := r.pool.Exec(ctx, query,
		ticker, company.INN, company.Owner, company.SectorID,
		company.LotSize, company.CEO, company.Employees,
	)

	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("company not found for ticker %s", ticker)
	}

	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, ticker string) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty")
	}

	query := `DELETE FROM companies WHERE ticker = $1`

	result, err := r.pool.Exec(ctx, query, ticker)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("company not found for ticker %s", ticker)
	}

	return nil
}
