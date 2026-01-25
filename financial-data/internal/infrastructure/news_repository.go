package infrastructure

import (
	"context"
	"errors"
	"financial_data/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepository struct {
	pool *pgxpool.Pool
}

func NewNewsRepository(pool *pgxpool.Pool) *NewsRepository {
	return &NewsRepository{pool: pool}
}

func (r *NewsRepository) GetByID(ctx context.Context, id int) (*domain.News, error) {
	if id < 1 {
		return nil, fmt.Errorf("invalid news ID: %d", id)
	}

	query := `
		SELECT id, ticker, sector_id, date, title, content, source, url
		FROM news
		WHERE id = $1
	`

	news := &domain.News{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&news.ID, &news.Ticker, &news.SectorID, &news.Date,
		&news.Title, &news.Content, &news.Source, &news.URL,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("news not found for ID %d", id)
		}
		return nil, fmt.Errorf("failed to get news: %w", err)
	}

	return news, nil
}

func (r *NewsRepository) GetByTicker(ctx context.Context, ticker string) ([]domain.News, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty")
	}

	query := `
		SELECT id, ticker, sector_id, date, title, content, source, url
		FROM news
		WHERE ticker = $1
		ORDER BY date DESC
	`

	rows, err := r.pool.Query(ctx, query, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to query news by ticker: %w", err)
	}
	defer rows.Close()

	var newsList []domain.News
	for rows.Next() {
		var news domain.News
		err := rows.Scan(
			&news.ID, &news.Ticker, &news.SectorID, &news.Date,
			&news.Title, &news.Content, &news.Source, &news.URL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan news: %w", err)
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating news: %w", err)
	}

	return newsList, nil
}

func (r *NewsRepository) GetBySector(ctx context.Context, sectorID int) ([]domain.News, error) {
	query := `
		SELECT id, ticker, sector_id, date, title, content, source, url
		FROM news
		WHERE sector_id = $1
		ORDER BY date DESC
	`

	rows, err := r.pool.Query(ctx, query, sectorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query news by sector: %w", err)
	}
	defer rows.Close()

	var newsList []domain.News
	for rows.Next() {
		var news domain.News
		err := rows.Scan(
			&news.ID, &news.Ticker, &news.SectorID, &news.Date,
			&news.Title, &news.Content, &news.Source, &news.URL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan news: %w", err)
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating news: %w", err)
	}

	return newsList, nil
}

func (r *NewsRepository) Create(ctx context.Context, news *domain.News) error {
	if news == nil {
		return fmt.Errorf("news is nil")
	}
	if news.Title == "" {
		return fmt.Errorf("news title is empty")
	}

	query := `
		INSERT INTO news (ticker, sector_id, date, title, content, source, url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		news.Ticker, news.SectorID, news.Date, news.Title,
		news.Content, news.Source, news.URL,
	).Scan(&news.ID)

	if err != nil {
		return fmt.Errorf("failed to create news: %w", err)
	}

	return nil
}

func (r *NewsRepository) Update(ctx context.Context, id int, news *domain.News) error {
	if id < 1 {
		return fmt.Errorf("invalid news ID: %d", id)
	}
	if news == nil {
		return fmt.Errorf("news is nil")
	}

	query := `
		UPDATE news SET
			ticker = $2, sector_id = $3, date = $4, title = $5,
			content = $6, source = $7, url = $8
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		id, news.Ticker, news.SectorID, news.Date,
		news.Title, news.Content, news.Source, news.URL,
	)

	if err != nil {
		return fmt.Errorf("failed to update news: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("news not found for ID %d", id)
	}

	return nil
}

func (r *NewsRepository) Delete(ctx context.Context, id int) error {
	if id < 1 {
		return fmt.Errorf("invalid news ID: %d", id)
	}

	query := `DELETE FROM news WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete news: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("news not found for ID %d", id)
	}

	return nil
}
