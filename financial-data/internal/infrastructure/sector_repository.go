package infrastructure

import (
	"context"
	"errors"
	"financial_data/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SectorRepository struct {
	pool *pgxpool.Pool
}

func NewSectorRepository(pool *pgxpool.Pool) *SectorRepository {
	return &SectorRepository{pool: pool}
}

func (r *SectorRepository) GetByID(ctx context.Context, id int) (*domain.SectorModel, error) {
	if id < 1 || id > 19 {
		return nil, fmt.Errorf("invalid sector ID: %d (allowed 1-19)", id)
	}

	query := `SELECT id, name FROM sectors WHERE id = $1`

	sector := &domain.SectorModel{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&sector.ID, &sector.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("sector not found for ID %d", id)
		}
		return nil, fmt.Errorf("failed to get sector: %w", err)
	}

	return sector, nil
}

func (r *SectorRepository) GetAll(ctx context.Context) ([]domain.SectorModel, error) {
	query := `SELECT id, name FROM sectors ORDER BY id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sectors: %w", err)
	}
	defer rows.Close()

	var sectors []domain.SectorModel
	for rows.Next() {
		var sector domain.SectorModel
		err := rows.Scan(&sector.ID, &sector.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sector: %w", err)
		}
		sectors = append(sectors, sector)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sectors: %w", err)
	}

	return sectors, nil
}

func (r *SectorRepository) Create(ctx context.Context, sector *domain.SectorModel) error {
	if sector == nil {
		return fmt.Errorf("sector is nil")
	}
	if sector.Name == "" {
		return fmt.Errorf("sector name is empty")
	}

	query := `INSERT INTO sectors (name) VALUES ($1) RETURNING id`

	err := r.pool.QueryRow(ctx, query, sector.Name).Scan(&sector.ID)
	if err != nil {
		return fmt.Errorf("failed to create sector: %w", err)
	}

	return nil
}

func (r *SectorRepository) Update(ctx context.Context, id int, sector *domain.SectorModel) error {
	if id < 1 {
		return fmt.Errorf("invalid sector ID: %d", id)
	}
	if sector == nil {
		return fmt.Errorf("sector is nil")
	}
	if sector.Name == "" {
		return fmt.Errorf("sector name is empty")
	}

	query := `UPDATE sectors SET name = $2 WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id, sector.Name)
	if err != nil {
		return fmt.Errorf("failed to update sector: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("sector not found for ID %d", id)
	}

	return nil
}

func (r *SectorRepository) Delete(ctx context.Context, id int) error {
	if id < 1 {
		return fmt.Errorf("invalid sector ID: %d", id)
	}

	query := `DELETE FROM sectors WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sector: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("sector not found for ID %d", id)
	}

	return nil
}
