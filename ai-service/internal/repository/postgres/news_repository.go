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

	latestJSON, err := json.Marshal(news.LatestNews)
	if err != nil {
		return fmt.Errorf("marshal latest_news: %w", err)
	}

	importantJSON, err := json.Marshal(news.ImportantNews)
	if err != nil {
		return fmt.Errorf("marshal important_news: %w", err)
	}

	_, err = db.Exec(ctx, `
		INSERT INTO company_news (ticker, latest_news, important_news)
		VALUES ($1, $2, $3)
		ON CONFLICT (ticker) DO UPDATE SET
			latest_news = EXCLUDED.latest_news,
			important_news = EXCLUDED.important_news,
			created_at = NOW()
	`, ticker, latestJSON, importantJSON)
	if err != nil {
		return fmt.Errorf("upsert company_news: %w", err)
	}

	return nil
}

func (r *NewsRepository) GetFreshNews(ctx context.Context, ticker string, ttl time.Duration) (*entity.NewsResponse, error) {
	db := Executor(ctx, r.db)

	var latestJSON, importantJSON []byte
	err := db.QueryRow(ctx, `
		SELECT latest_news, important_news
		FROM company_news
		WHERE ticker = $1 AND created_at > NOW() - $2::interval
	`, ticker, ttl.String()).Scan(&latestJSON, &importantJSON)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get fresh news: %w", err)
	}

	var news entity.NewsResponse
	if err := json.Unmarshal(latestJSON, &news.LatestNews); err != nil {
		return nil, fmt.Errorf("unmarshal latest_news: %w", err)
	}
	if err := json.Unmarshal(importantJSON, &news.ImportantNews); err != nil {
		return nil, fmt.Errorf("unmarshal important_news: %w", err)
	}

	return &news, nil
}
