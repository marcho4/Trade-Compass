package infrastructure

import (
	"context"
	"fmt"

	"financial_data/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RatiosRepository struct {
	pool *pgxpool.Pool
}

func NewRatiosRepository(pool *pgxpool.Pool) *RatiosRepository {
	return &RatiosRepository{pool: pool}
}

func (r *RatiosRepository) GetByTicker(ctx context.Context, ticker string) (*domain.Ratios, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty")
	}
	return nil, fmt.Errorf("not implemented")
}
