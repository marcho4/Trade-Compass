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
		return nil, NewDbError(fmt.Sprintf("invalid sector ID: %d (allowed 1-19)", id), 0)
	}

	query := `SELECT id, name FROM sectors WHERE id = $1`

	sector := &domain.SectorModel{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&sector.ID, &sector.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDbError(fmt.Sprintf("sector not found for ID %d", id), 0)
		}
		return nil, NewDbError(fmt.Sprintf("failed to get sector: %v", err), 0)
	}

	return sector, nil
}

func (r *SectorRepository) GetAll(ctx context.Context) ([]domain.SectorModel, error) {
	query := `SELECT id, name FROM sectors ORDER BY id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, NewDbError(fmt.Sprintf("failed to query sectors: %v", err), 0)
	}
	defer rows.Close()

	var sectors []domain.SectorModel
	for rows.Next() {
		var sector domain.SectorModel
		err := rows.Scan(&sector.ID, &sector.Name)
		if err != nil {
			return nil, NewDbError(fmt.Sprintf("failed to scan sector: %v", err), 0)
		}
		sectors = append(sectors, sector)
	}

	if err := rows.Err(); err != nil {
		return nil, NewDbError(fmt.Sprintf("error iterating sectors: %v", err), 0)
	}

	return sectors, nil
}

func (r *SectorRepository) Create(ctx context.Context, sector *domain.SectorModel) error {
	if sector == nil {
		return NewDbError("sector is nil", 0)
	}
	if sector.Name == "" {
		return NewDbError("sector name is empty", 0)
	}

	query := `INSERT INTO sectors (name) VALUES ($1) RETURNING id`

	err := r.pool.QueryRow(ctx, query, sector.Name).Scan(&sector.ID)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to create sector: %v", err), 0)
	}

	return nil
}

func (r *SectorRepository) Update(ctx context.Context, id int, sector *domain.SectorModel) error {
	if id < 1 {
		return NewDbError(fmt.Sprintf("invalid sector ID: %d", id), 0)
	}
	if sector == nil {
		return NewDbError("sector is nil", 0)
	}
	if sector.Name == "" {
		return NewDbError("sector name is empty", 0)
	}

	query := `UPDATE sectors SET name = $2 WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id, sector.Name)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to update sector: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("sector not found for ID %d", id), 0)
	}

	return nil
}

func (r *SectorRepository) Delete(ctx context.Context, id int) error {
	if id < 1 {
		return NewDbError(fmt.Sprintf("invalid sector ID: %d", id), 0)
	}

	query := `DELETE FROM sectors WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to delete sector: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("sector not found for ID %d", id), 0)
	}

	return nil
}
