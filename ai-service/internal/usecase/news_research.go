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
	ai        AIProvider
	news      NewsRepository
	publisher MessagePublisher
	newsTTL   time.Duration
}

func NewNewsResearchUsecase(ai AIProvider, news NewsRepository, publisher MessagePublisher, newsTTL time.Duration) *NewsResearchUsecase {
	return &NewsResearchUsecase{
		ai:        ai,
		news:      news,
		publisher: publisher,
		newsTTL:   newsTTL,
	}
}

func (u *NewsResearchUsecase) Execute(ctx context.Context, task entity.Task) error {
	logger := slog.With(
		slog.String("id", task.Id),
		slog.String("ticker", task.Ticker),
	)

	logger.Info("starting news research task")

	existing, err := u.news.GetFreshNews(ctx, task.Ticker, u.newsTTL)
	if err != nil {
		return fmt.Errorf("check existing news: %w", err)
	}

	if existing != nil {
		logger.Info("fresh news already exist, skipping")

		if err := u.sendNextTask(ctx, task); err != nil {
			return fmt.Errorf("send task: %w", err)
		}

		return nil
	}

	newsItemSchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"news":        {Type: TypeString},
			"date":        {Type: TypeString},
			"source":      {Type: TypeString},
			"severity":    {Type: TypeString, Enum: []string{"high", "medium", "low"}},
			"impact_type": {Type: TypeString, Enum: []string{"positive", "negative", "neutral"}},
		},
		Required: []string{"news", "date", "source", "severity", "impact_type"},
	}

	responseSchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"latest_news":    {Type: TypeArray, Items: newsItemSchema},
			"important_news": {Type: TypeArray, Items: newsItemSchema},
		},
		Required: []string{"latest_news", "important_news"},
	}

	logger.Info("calling AI to collect news")
	start := time.Now()

	text, err := u.ai.GenerateText(ctx, buildNewsAgentPrompt(task.Ticker), entity.Flash, GenerateParams{
		GoogleSearch:   true,
		ResponseSchema: responseSchema,
	})
	if err != nil {
		return fmt.Errorf("call ai provider for news: %w", err)
	}

	requestTime := time.Since(start)
	logger.Info("Gemini call lasted", "time", requestTime.String())

	var res entity.NewsResponse
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		slog.Error("failed to parse news response", slog.String("ai_response", text))
		return fmt.Errorf("parse news response: %w", err)
	}

	if err := u.news.SaveNews(ctx, task.Ticker, &res); err != nil {
		return fmt.Errorf("save news: %w", err)
	}

	logger.Info("news research completed and saved")

	if err := u.sendNextTask(ctx, task); err != nil {
		return fmt.Errorf("send task: %w", err)
	}

	logger.Info("published risk-and-growth task")

	return nil
}

func (u *NewsResearchUsecase) sendNextTask(ctx context.Context, task entity.Task) error {
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

	return nil
}
