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
		return nil, NewDbError(fmt.Sprintf("invalid news ID: %d", id), 0)
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
			return nil, NewDbError(fmt.Sprintf("news not found for ID %d", id), 0)
		}
		return nil, NewDbError(fmt.Sprintf("failed to get news: %v", err), 0)
	}

	return news, nil
}

func (r *NewsRepository) GetByTicker(ctx context.Context, ticker string) ([]domain.News, error) {
	if ticker == "" {
		return nil, NewDbError("ticker is empty", 0)
	}

	query := `
		SELECT id, ticker, sector_id, date, title, content, source, url
		FROM news
		WHERE ticker = $1
		ORDER BY date DESC
	`

	rows, err := r.pool.Query(ctx, query, ticker)
	if err != nil {
		return nil, NewDbError(fmt.Sprintf("failed to query news by ticker: %v", err), 0)
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
			return nil, NewDbError(fmt.Sprintf("failed to scan news: %v", err), 0)
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, NewDbError(fmt.Sprintf("error iterating news: %v", err), 0)
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
		return nil, NewDbError(fmt.Sprintf("failed to query news by sector: %v", err), 0)
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
			return nil, NewDbError(fmt.Sprintf("failed to scan news: %v", err), 0)
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, NewDbError(fmt.Sprintf("error iterating news: %v", err), 0)
	}

	return newsList, nil
}

func (r *NewsRepository) Create(ctx context.Context, news *domain.News) error {
	if news == nil {
		return NewDbError("news is nil", 0)
	}
	if news.Title == "" {
		return NewDbError("news title is empty", 0)
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
		return NewDbError(fmt.Sprintf("failed to create news: %v", err), 0)
	}

	return nil
}

func (r *NewsRepository) Update(ctx context.Context, id int, news *domain.News) error {
	if id < 1 {
		return NewDbError(fmt.Sprintf("invalid news ID: %d", id), 0)
	}
	if news == nil {
		return NewDbError("news is nil", 0)
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
		return NewDbError(fmt.Sprintf("failed to update news: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("news not found for ID %d", id), 0)
	}

	return nil
}

func (r *NewsRepository) Delete(ctx context.Context, id int) error {
	if id < 1 {
		return NewDbError(fmt.Sprintf("invalid news ID: %d", id), 0)
	}

	query := `DELETE FROM news WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to delete news: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("news not found for ID %d", id), 0)
	}

	return nil
}
