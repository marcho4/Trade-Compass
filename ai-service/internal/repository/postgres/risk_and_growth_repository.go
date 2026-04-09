package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ai-service/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RiskAndGrowthRepository struct {
	db *pgxpool.Pool
}

func NewRiskAndGrowthRepository(db *pgxpool.Pool) *RiskAndGrowthRepository {
	return &RiskAndGrowthRepository{db: db}
}

func (r *RiskAndGrowthRepository) SaveRiskAndGrowth(ctx context.Context, response *entity.RiskAndGrowthResponse) error {
	db := Executor(ctx, r.db)

	factorsJSON, err := json.Marshal(response.Factors)
	if err != nil {
		return fmt.Errorf("marshal factors: %w", err)
	}

	_, err = db.Exec(ctx, `
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

func (r *RiskAndGrowthRepository) GetFreshRiskAndGrowth(ctx context.Context, ticker string, ttl time.Duration) (*entity.RiskAndGrowthResponse, error) {
	db := Executor(ctx, r.db)

	var factorsJSON []byte
	err := db.QueryRow(ctx, `
		SELECT factors
		FROM risk_and_growth
		WHERE ticker = $1 AND created_at > NOW() - make_interval(secs => $2)
	`, ticker, ttl.Seconds()).Scan(&factorsJSON)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get fresh risk and growth: %w", err)
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
