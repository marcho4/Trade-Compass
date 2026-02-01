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
		return nil, NewDbError("ticker is empty", 0)
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
			return nil, NewDbError(fmt.Sprintf("company not found for ticker %s", ticker), 0)
		}
		return nil, NewDbError(fmt.Sprintf("failed to get company: %v", err), 0)
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
		return nil, NewDbError(fmt.Sprintf("failed to query companies: %v", err), 0)
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
			return nil, NewDbError(fmt.Sprintf("failed to scan company: %v", err), 0)
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, NewDbError(fmt.Sprintf("error iterating companies: %v", err), 0)
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
		return nil, NewDbError(fmt.Sprintf("failed to query companies by sector: %v", err), 0)
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
			return nil, NewDbError(fmt.Sprintf("failed to scan company: %v", err), 0)
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, NewDbError(fmt.Sprintf("error iterating companies: %v", err), 0)
	}

	return companies, nil
}

func (r *CompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	if company == nil {
		return NewDbError("company is nil", 0)
	}
	if company.INN == "" {
		return NewDbError("INN is empty", 0)
	}
	if company.Ticker == "" {
		return NewDbError("ticker is empty", 0)
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
		return NewDbError(fmt.Sprintf("failed to create company: %v", err), 0)
	}

	return nil
}

func (r *CompanyRepository) Update(ctx context.Context, ticker string, company *domain.Company) error {
	if ticker == "" {
		return NewDbError("ticker is empty", 0)
	}
	if company == nil {
		return NewDbError("company is nil", 0)
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
		return NewDbError(fmt.Sprintf("failed to update company: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("company not found for ticker %s", ticker), 0)
	}

	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, ticker string) error {
	if ticker == "" {
		return NewDbError("ticker is empty", 0)
	}

	query := `DELETE FROM companies WHERE ticker = $1`

	result, err := r.pool.Exec(ctx, query, ticker)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to delete company: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("company not found for ticker %s", ticker), 0)
	}

	return nil
}
