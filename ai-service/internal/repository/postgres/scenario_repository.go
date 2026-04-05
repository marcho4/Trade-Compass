package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ScenarioRepository struct {
	db *pgxpool.Pool
}

func NewScenarioRepository(db *pgxpool.Pool) *ScenarioRepository {
	return &ScenarioRepository{db: db}
}

func (r *ScenarioRepository) SaveScenarios(ctx context.Context, ticker string, scenarios []entity.Scenario) error {
	db := Executor(ctx, r.db)

	for _, s := range scenarios {
		growthJSON, err := json.Marshal(s.GrowthFactorsApplied)
		if err != nil {
			return fmt.Errorf("marshal growth factors: %w", err)
		}

		risksJSON, err := json.Marshal(s.RisksApplied)
		if err != nil {
			return fmt.Errorf("marshal risks: %w", err)
		}

		assumptionsJSON, err := json.Marshal(s.Assumptions)
		if err != nil {
			return fmt.Errorf("marshal assumptions: %w", err)
		}

		_, err = db.Exec(ctx, `
			INSERT INTO scenarios (id, ticker, name, description, probability, terminal_growth_rate, growth_factors_applied, risks_applied, assumptions)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id, ticker) DO UPDATE SET
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				probability = EXCLUDED.probability,
				terminal_growth_rate = EXCLUDED.terminal_growth_rate,
				growth_factors_applied = EXCLUDED.growth_factors_applied,
				risks_applied = EXCLUDED.risks_applied,
				assumptions = EXCLUDED.assumptions,
				created_at = NOW()
		`, s.ID, ticker, s.Name, s.Description, s.Probability, s.TerminalGrowthRate, growthJSON, risksJSON, assumptionsJSON)
		if err != nil {
			return fmt.Errorf("upsert scenario %s: %w", s.ID, err)
		}
	}

	return nil
}

func (r *ScenarioRepository) GetScenarios(ctx context.Context, ticker string) ([]entity.Scenario, error) {
	db := Executor(ctx, r.db)

	rows, err := db.Query(ctx, `
		SELECT id, name, description, probability, terminal_growth_rate, growth_factors_applied, risks_applied, assumptions
		FROM scenarios
		WHERE ticker = $1
	`, ticker)
	if err != nil {
		return nil, fmt.Errorf("query scenarios: %w", err)
	}
	defer rows.Close()

	var scenarios []entity.Scenario
	for rows.Next() {
		var s entity.Scenario
		var growthJSON, risksJSON, assumptionsJSON []byte

		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Probability, &s.TerminalGrowthRate, &growthJSON, &risksJSON, &assumptionsJSON); err != nil {
			return nil, fmt.Errorf("scan scenario: %w", err)
		}

		if growthJSON != nil {
			if err := json.Unmarshal(growthJSON, &s.GrowthFactorsApplied); err != nil {
				return nil, fmt.Errorf("unmarshal growth factors: %w", err)
			}
		}

		if risksJSON != nil {
			if err := json.Unmarshal(risksJSON, &s.RisksApplied); err != nil {
				return nil, fmt.Errorf("unmarshal risks: %w", err)
			}
		}

		if err := json.Unmarshal(assumptionsJSON, &s.Assumptions); err != nil {
			return nil, fmt.Errorf("unmarshal assumptions: %w", err)
		}

		scenarios = append(scenarios, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if len(scenarios) == 0 {
		return nil, domain.ErrScenariosNotFound
	}

	return scenarios, nil
}
