package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"
)

type BusinessResearchUsecase struct {
	ai        AIService
	repo      BusinessResearchRepository
	publisher MessagePublisher
}

func NewBusinessResearchUsecase(ai AIService, repo BusinessResearchRepository, publisher MessagePublisher) *BusinessResearchUsecase {
	return &BusinessResearchUsecase{
		ai:        ai,
		repo:      repo,
		publisher: publisher,
	}
}

func (u *BusinessResearchUsecase) Execute(ctx context.Context, task entity.Task) error {
	logger := slog.With(slog.String("ticker", task.Ticker))

	logger.Info("starting business research task")

	existing, err := u.repo.GetBusinessResearch(ctx, task.Ticker)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("check existing business research: %w", err)
	}

	if existing != nil {
		logger.Info("business research already exists, skipping")
		return nil
	}

	result, err := u.ai.ResearchBusiness(ctx, task.Ticker, task.Ticker)
	if err != nil {
		return fmt.Errorf("research business: %w", err)
	}

	if err := u.repo.SaveBusinessResearch(ctx, result); err != nil {
		return fmt.Errorf("save business research: %w", err)
	}

	logger.Info("business research completed and saved")

	nextTask := entity.Task{
		Id:     task.Id,
		Ticker: task.Ticker,
		Type:   entity.NewsResearch,
	}

	payload, err := json.Marshal(nextTask)
	if err != nil {
		return fmt.Errorf("marshal news-research task: %w", err)
	}

	if err := u.publisher.PublishMessage(ctx, payload); err != nil {
		return fmt.Errorf("publish news-research task: %w", err)
	}

	logger.Info("published news-research task")
	return nil
}

func (u *BusinessResearchUsecase) GetBusinessResearch(ctx context.Context, ticker string) (*entity.BusinessResearchResult, error) {
	return u.repo.GetBusinessResearch(ctx, ticker)
}
