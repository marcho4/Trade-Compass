package infrastructure

import (
	"context"
	"errors"
	"financial_data/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DividendsRepository struct {
	pool *pgxpool.Pool
}

func NewDividendsRepository(pool *pgxpool.Pool) *DividendsRepository {
	return &DividendsRepository{pool: pool}
}

func (r *DividendsRepository) GetByTicker(ctx context.Context, ticker string) ([]domain.Dividends, error) {
	if ticker == "" {
		return nil, NewDbError("ticker is empty", 0)
	}

	query := `
		SELECT id, ticker, ex_dividend_date, payment_date, amount_per_share, 
		       dividend_yield, payout_ratio, currency
		FROM dividends
		WHERE ticker = $1
		ORDER BY ex_dividend_date DESC
	`

	rows, err := r.pool.Query(ctx, query, ticker)
	if err != nil {
		return nil, NewDbError(fmt.Sprintf("failed to query dividends: %v", err), 0)
	}
	defer rows.Close()

	var dividends []domain.Dividends
	for rows.Next() {
		var div domain.Dividends
		err := rows.Scan(
			&div.ID, &div.Ticker, &div.ExDividendDate, &div.PaymentDate,
			&div.AmountPerShare, &div.DividendYield, &div.PayoutRatio, &div.Currency,
		)
		if err != nil {
			return nil, NewDbError(fmt.Sprintf("failed to scan dividend: %v", err), 0)
		}
		dividends = append(dividends, div)
	}

	if err := rows.Err(); err != nil {
		return nil, NewDbError(fmt.Sprintf("error iterating dividends: %v", err), 0)
	}

	return dividends, nil
}

func (r *DividendsRepository) GetByID(ctx context.Context, id int) (*domain.Dividends, error) {
	if id < 1 {
		return nil, NewDbError(fmt.Sprintf("invalid dividend ID: %d", id), 0)
	}

	query := `
		SELECT id, ticker, ex_dividend_date, payment_date, amount_per_share, 
		       dividend_yield, payout_ratio, currency
		FROM dividends
		WHERE id = $1
	`

	div := &domain.Dividends{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&div.ID, &div.Ticker, &div.ExDividendDate, &div.PaymentDate,
		&div.AmountPerShare, &div.DividendYield, &div.PayoutRatio, &div.Currency,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDbError(fmt.Sprintf("dividend not found for ID %d", id), 0)
		}
		return nil, NewDbError(fmt.Sprintf("failed to get dividend: %v", err), 0)
	}

	return div, nil
}

func (r *DividendsRepository) Create(ctx context.Context, dividend *domain.Dividends) error {
	if dividend == nil {
		return NewDbError("dividend is nil", 0)
	}
	if dividend.Ticker == "" {
		return NewDbError("ticker is empty", 0)
	}

	query := `
		INSERT INTO dividends (ticker, ex_dividend_date, payment_date, amount_per_share, 
		                       dividend_yield, payout_ratio, currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		dividend.Ticker, dividend.ExDividendDate, dividend.PaymentDate,
		dividend.AmountPerShare, dividend.DividendYield, dividend.PayoutRatio, dividend.Currency,
	).Scan(&dividend.ID)

	if err != nil {
		return NewDbError(fmt.Sprintf("failed to create dividend: %v", err), 0)
	}

	return nil
}

func (r *DividendsRepository) Update(ctx context.Context, id int, dividend *domain.Dividends) error {
	if id < 1 {
		return NewDbError(fmt.Sprintf("invalid dividend ID: %d", id), 0)
	}
	if dividend == nil {
		return NewDbError("dividend is nil", 0)
	}

	query := `
		UPDATE dividends SET
			ex_dividend_date = $2, payment_date = $3, amount_per_share = $4,
			dividend_yield = $5, payout_ratio = $6, currency = $7
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		id, dividend.ExDividendDate, dividend.PaymentDate,
		dividend.AmountPerShare, dividend.DividendYield, dividend.PayoutRatio, dividend.Currency,
	)

	if err != nil {
		return NewDbError(fmt.Sprintf("failed to update dividend: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("dividend not found for ID %d", id), 0)
	}

	return nil
}

func (r *DividendsRepository) Delete(ctx context.Context, id int) error {
	if id < 1 {
		return NewDbError(fmt.Sprintf("invalid dividend ID: %d", id), 0)
	}

	query := `DELETE FROM dividends WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to delete dividend: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("dividend not found for ID %d", id), 0)
	}

	return nil
}
