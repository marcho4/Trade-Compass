package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	docs "ai-service/internal/docs"
	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"
)

type BusinessResearchUsecase struct {
	ai        AIProvider
	repo      BusinessResearchRepository
	publisher MessagePublisher
}

func NewBusinessResearchUsecase(ai AIProvider, repo BusinessResearchRepository, publisher MessagePublisher) *BusinessResearchUsecase {
	return &BusinessResearchUsecase{
		ai:        ai,
		repo:      repo,
		publisher: publisher,
	}
}

func (u *BusinessResearchUsecase) Execute(ctx context.Context, task entity.Task) error {
	logger := slog.With(
		slog.String("id", task.Id),
		slog.String("ticker", task.Ticker),
	)

	logger.Info("starting business research task")

	existing, err := u.repo.GetBusinessResearch(ctx, task.Ticker)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("check existing business research: %w", err)
	}

	if existing != nil {
		logger.Info("business research already exists, skipping")

		if err := u.publishNextTask(ctx, task); err != nil {
			return fmt.Errorf("publish next task: %w", err)
		}

		return nil
	}

	prompt := docs.BusinessResearcherPrompt() +
		"\n\n## Компания для анализа\nТикер: " + task.Ticker +
		"\n\nВАЖНО: В поле ticker ответа используй СТРОГО \"" + task.Ticker + "\". Не заменяй тикер на альтернативный."

	marketSchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"market": {Type: TypeString},
			"role":   {Type: TypeString},
		},
		Required: []string{"market", "role"},
	}

	revenueSchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"segment":     {Type: TypeString},
			"share_pct":   {Type: TypeNumber},
			"approximate": {Type: TypeBoolean},
			"description": {Type: TypeString},
			"trend":       {Type: TypeString, Enum: []string{"growing", "stable", "declining"}},
		},
		Required: []string{"segment", "share_pct", "approximate", "description", "trend"},
	}

	dependencySchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"factor":      {Type: TypeString},
			"type":        {Type: TypeString, Enum: []string{"commodity", "currency", "regulation", "macro", "technology", "geopolitics", "infrastructure", "demand"}},
			"severity":    {Type: TypeString, Enum: []string{"critical", "high", "moderate"}},
			"description": {Type: TypeString},
		},
		Required: []string{"factor", "type", "severity", "description"},
	}

	responseSchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"ticker":       {Type: TypeString},
			"company_name": {Type: TypeString},
			"profile": {
				Type: TypeObject,
				Properties: map[string]*Schema{
					"description":           {Type: TypeString},
					"products_and_services": {Type: TypeArray, Items: &Schema{Type: TypeString}},
					"markets":               {Type: TypeArray, Items: marketSchema},
					"key_clients":           {Type: TypeString},
					"business_model":        {Type: TypeString},
				},
				Required: []string{"description", "products_and_services", "markets", "key_clients", "business_model"},
			},
			"revenue_sources": {Type: TypeArray, Items: revenueSchema},
			"dependencies":    {Type: TypeArray, Items: dependencySchema},
		},
		Required: []string{"ticker", "company_name", "profile", "revenue_sources", "dependencies"},
	}

	logger.Info("calling AI for business research")
	text, err := u.ai.GenerateText(ctx, prompt, entity.Flash, GenerateParams{
		GoogleSearch:   true,
		ResponseSchema: responseSchema,
	})
	if err != nil {
		return fmt.Errorf("research business: %w", err)
	}

	var res entity.BusinessResearchResponse
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		logger.Error("failed to parse business research response", slog.String("ai_response", text))
		return fmt.Errorf("parse business research response: %w", err)
	}

	res.Ticker = task.Ticker

	logger.Info("business research completed")

	if err := u.repo.SaveBusinessResearch(ctx, &res); err != nil {
		return fmt.Errorf("save business research: %w", err)
	}

	logger.Info("business research completed and saved")

	if err := u.publishNextTask(ctx, task); err != nil {
		return fmt.Errorf("publish next task: %w", err)
	}

	return nil
}

func (u *BusinessResearchUsecase) GetBusinessResearch(ctx context.Context, ticker string) (*entity.BusinessResearchResult, error) {
	return u.repo.GetBusinessResearch(ctx, ticker)
}

func (u *BusinessResearchUsecase) publishNextTask(ctx context.Context, task entity.Task) error {
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

	slog.Info("published news-research task")

	return nil
}
