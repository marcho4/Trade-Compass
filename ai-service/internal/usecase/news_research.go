package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"ai-service/internal/domain/entity"
)

type NewsResearchUsecase struct {
	ai        AIService
	news      NewsRepository
	publisher MessagePublisher
	newsTTL   time.Duration
}

func NewNewsResearchUsecase(ai AIService, news NewsRepository, publisher MessagePublisher, newsTTL time.Duration) *NewsResearchUsecase {
	return &NewsResearchUsecase{
		ai:        ai,
		news:      news,
		publisher: publisher,
		newsTTL:   newsTTL,
	}
}

func (u *NewsResearchUsecase) Execute(ctx context.Context, task entity.Task) error {
	logger := slog.With(slog.String("ticker", task.Ticker))

	logger.Info("starting news research task")

	existing, err := u.news.GetFreshNews(ctx, task.Ticker, u.newsTTL)
	if err != nil {
		return fmt.Errorf("check existing news: %w", err)
	}

	if existing != nil {
		logger.Info("fresh news already exist, skipping")
		return nil
	}

	news, err := u.ai.CollectNews(ctx, task.Ticker)
	if err != nil {
		return fmt.Errorf("collect news: %w", err)
	}

	if err := u.news.SaveNews(ctx, task.Ticker, news); err != nil {
		return fmt.Errorf("save news: %w", err)
	}

	logger.Info("news research completed and saved")

	nextTask := entity.Task{
		Id:     task.Id,
		Ticker: task.Ticker,
		Type:   entity.RiskAndGrowth,
	}

	payload, err := json.Marshal(nextTask)
	if err != nil {
		return fmt.Errorf("marshal risk-and-growth task: %w", err)
	}

	if err := u.publisher.PublishMessage(ctx, payload); err != nil {
		return fmt.Errorf("publish risk-and-growth task: %w", err)
	}

	logger.Info("published risk-and-growth task")
	return nil
}
