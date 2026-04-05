package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DCFResultsRepository struct {
	db *pgxpool.Pool
}

func NewDCFResultsRepository(db *pgxpool.Pool) *DCFResultsRepository {
	return &DCFResultsRepository{db: db}
}

func (r *DCFResultsRepository) SaveDCFResults(ctx context.Context, ticker string, result entity.DCFResult) error {
	db := Executor(ctx, r.db)

	for _, s := range result.Scenarios {
		fcfsJSON, err := json.Marshal(s.YearlyFCFs)
		if err != nil {
			return fmt.Errorf("marshal yearly fcfs: %w", err)
		}

		_, err = db.Exec(ctx, `
			INSERT INTO dcf_results (id, ticker, scenario_type, probability, enterprise_value, equity_value, price_per_share, terminal_value, yearly_fcfs)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id, ticker, scenario_type) DO UPDATE SET
				probability = EXCLUDED.probability,
				enterprise_value = EXCLUDED.enterprise_value,
				equity_value = EXCLUDED.equity_value,
				price_per_share = EXCLUDED.price_per_share,
				terminal_value = EXCLUDED.terminal_value,
				yearly_fcfs = EXCLUDED.yearly_fcfs,
				created_at = NOW()
		`, result.ID, ticker, s.ScenarioID, s.Probability, s.EnterpriseValue, s.EquityValue, s.PricePerShare, s.TerminalValue, fcfsJSON)
		if err != nil {
			return fmt.Errorf("upsert dcf result for scenario %s: %w", s.ScenarioID, err)
		}
	}

	return nil
}

func (r *DCFResultsRepository) GetDCFResults(ctx context.Context, ticker string) (*entity.DCFResult, error) {
	db := Executor(ctx, r.db)

	rows, err := db.Query(ctx, `
		SELECT id, scenario_type, probability, enterprise_value, equity_value, price_per_share, terminal_value, yearly_fcfs
		FROM dcf_results
		WHERE ticker = $1
	`, ticker)
	if err != nil {
		return nil, fmt.Errorf("query dcf results: %w", err)
	}
	defer rows.Close()

	result := &entity.DCFResult{}
	for rows.Next() {
		var s entity.ScenarioDCFResult
		var fcfsJSON []byte

		if err := rows.Scan(&result.ID, &s.ScenarioID, &s.Probability, &s.EnterpriseValue, &s.EquityValue, &s.PricePerShare, &s.TerminalValue, &fcfsJSON); err != nil {
			return nil, fmt.Errorf("scan dcf result: %w", err)
		}

		if err := json.Unmarshal(fcfsJSON, &s.YearlyFCFs); err != nil {
			return nil, fmt.Errorf("unmarshal yearly fcfs: %w", err)
		}

		result.Scenarios = append(result.Scenarios, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if len(result.Scenarios) == 0 {
		return nil, domain.ErrDCFResultsNotFound
	}

	result.ComputeWeighted()

	return result, nil
}
