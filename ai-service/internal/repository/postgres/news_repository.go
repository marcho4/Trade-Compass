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

type NewsRepository struct {
	db *pgxpool.Pool
}

func NewNewsRepository(db *pgxpool.Pool) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) SaveNews(ctx context.Context, ticker string, news *entity.NewsResponse) error {
	db := Executor(ctx, r.db)

	data, err := json.Marshal(news)
	if err != nil {
		return fmt.Errorf("marshal news: %w", err)
	}

	_, err = db.Exec(ctx, `
		INSERT INTO company_news (ticker, data)
		VALUES ($1, $2)
		ON CONFLICT (ticker) DO UPDATE SET
			data = EXCLUDED.data,
			created_at = NOW()
	`, ticker, data)
	if err != nil {
		return fmt.Errorf("upsert company_news: %w", err)
	}

	return nil
}

func (r *NewsRepository) GetFreshNews(ctx context.Context, ticker string, ttl time.Duration) (*entity.NewsResponse, error) {
	db := Executor(ctx, r.db)

	var data []byte
	err := db.QueryRow(ctx, `
		SELECT data
		FROM company_news
		WHERE ticker = $1 AND created_at > NOW() - $2::interval
	`, ticker, ttl.String()).Scan(&data)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get fresh news: %w", err)
	}

	var news entity.NewsResponse
	if err := json.Unmarshal(data, &news); err != nil {
		return nil, fmt.Errorf("unmarshal news: %w", err)
	}

	return &news, nil
}
