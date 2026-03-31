package port

import (
	"ai-service/internal/domain/entity"
	"context"
	"time"
)

type NewsRepository interface {
	SaveNews(ctx context.Context, ticker string, news *entity.NewsResponse) error
	GetFreshNews(ctx context.Context, ticker string, ttl time.Duration) (*entity.NewsResponse, error)
}
