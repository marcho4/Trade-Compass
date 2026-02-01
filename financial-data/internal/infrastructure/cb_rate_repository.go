package infrastructure

import (
	"context"
	"errors"
	"financial_data/internal/domain"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CBRateRepository struct {
	pool *pgxpool.Pool
}

func NewCBRateRepository(pool *pgxpool.Pool) *CBRateRepository {
	return &CBRateRepository{pool: pool}
}

func (r *CBRateRepository) GetCurrent(ctx context.Context) (*domain.CBRate, error) {
	query := `SELECT date, rate FROM cb_rates ORDER BY date DESC LIMIT 1`

	rate := &domain.CBRate{}
	err := r.pool.QueryRow(ctx, query).Scan(&rate.Date, &rate.Rate)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDbError("no CB rates found", 0)
		}
		return nil, NewDbError(fmt.Sprintf("failed to get current CB rate: %v", err), 0)
	}

	return rate, nil
}

func (r *CBRateRepository) GetByDate(ctx context.Context, date time.Time) (*domain.CBRate, error) {
	query := `SELECT date, rate FROM cb_rates WHERE date = $1`

	rate := &domain.CBRate{}
	err := r.pool.QueryRow(ctx, query, date).Scan(&rate.Date, &rate.Rate)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDbError(fmt.Sprintf("CB rate not found for date %s", date.Format("2006-01-02")), 0)
		}
		return nil, NewDbError(fmt.Sprintf("failed to get CB rate: %v", err), 0)
	}

	return rate, nil
}

func (r *CBRateRepository) GetHistory(ctx context.Context, from, to time.Time) ([]domain.CBRate, error) {
	query := `
		SELECT date, rate 
		FROM cb_rates 
		WHERE date >= $1 AND date <= $2
		ORDER BY date DESC
	`

	rows, err := r.pool.Query(ctx, query, from, to)
	if err != nil {
		return nil, NewDbError(fmt.Sprintf("failed to query CB rates history: %v", err), 0)
	}
	defer rows.Close()

	var rates []domain.CBRate
	for rows.Next() {
		var rate domain.CBRate
		err := rows.Scan(&rate.Date, &rate.Rate)
		if err != nil {
			return nil, NewDbError(fmt.Sprintf("failed to scan CB rate: %v", err), 0)
		}
		rates = append(rates, rate)
	}

	if err := rows.Err(); err != nil {
		return nil, NewDbError(fmt.Sprintf("error iterating CB rates: %v", err), 0)
	}

	if len(rates) == 0 {
		return nil, NewDbError("no CB rates found for the specified period", 0)
	}

	return rates, nil
}

func (r *CBRateRepository) Create(ctx context.Context, rate *domain.CBRate) error {
	if rate == nil {
		return NewDbError("rate is nil", 0)
	}

	query := `INSERT INTO cb_rates (date, rate) VALUES ($1, $2)`

	_, err := r.pool.Exec(ctx, query, rate.Date, rate.Rate)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to create CB rate: %v", err), 0)
	}

	return nil
}

func (r *CBRateRepository) Update(ctx context.Context, date time.Time, rate float64) error {
	query := `UPDATE cb_rates SET rate = $2 WHERE date = $1`

	result, err := r.pool.Exec(ctx, query, date, rate)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to update CB rate: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("CB rate not found for date %s", date.Format("2006-01-02")), 0)
	}

	return nil
}

func (r *CBRateRepository) Delete(ctx context.Context, date time.Time) error {
	query := `DELETE FROM cb_rates WHERE date = $1`

	result, err := r.pool.Exec(ctx, query, date)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to delete CB rate: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("CB rate not found for date %s", date.Format("2006-01-02")), 0)
	}

	return nil
}
