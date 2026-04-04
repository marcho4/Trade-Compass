package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"ai-service/internal/domain/entity"
)

type RiskAndGrowthUsecase struct {
	ai       AIService
	rag      RiskAndGrowthRepository
	news     NewsRepository
	business BusinessResearchRepository
	ttl      time.Duration
}

func NewRiskAndGrowthUsecase(
	ai AIService,
	rag RiskAndGrowthRepository,
	news NewsRepository,
	business BusinessResearchRepository,
	ttl time.Duration,
) *RiskAndGrowthUsecase {
	return &RiskAndGrowthUsecase{
		ai:       ai,
		rag:      rag,
		news:     news,
		business: business,
		ttl:      ttl,
	}
}

func (u *RiskAndGrowthUsecase) Execute(ctx context.Context, task entity.Task) error {
	logger := slog.With(slog.String("ticker", task.Ticker))

	logger.Info("starting risk and growth task")

	existing, err := u.rag.GetFreshRiskAndGrowth(ctx, task.Ticker, u.ttl)
	if err != nil {
		return fmt.Errorf("check existing risk and growth: %w", err)
	}

	if existing != nil {
		logger.Info("fresh risk and growth already exists, skipping")
		return nil
	}

	newsData, err := u.news.GetFreshNews(ctx, task.Ticker, u.ttl)
	if err != nil {
		return fmt.Errorf("get fresh news: %w", err)
	}

	businessData, err := u.business.GetBusinessResearch(ctx, task.Ticker)
	if err != nil {
		return fmt.Errorf("get business research: %w", err)
	}

	result, err := u.ai.ExtractRiskAndGrowth(ctx, task.Ticker, newsData, businessData)
	if err != nil {
		return fmt.Errorf("extract risk and growth: %w", err)
	}

	if err := u.rag.SaveRiskAndGrowth(ctx, result); err != nil {
		return fmt.Errorf("save risk and growth: %w", err)
	}

	logger.Info("risk and growth completed and saved")
	return nil
}
