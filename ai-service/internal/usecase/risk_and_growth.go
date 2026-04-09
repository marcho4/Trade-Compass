package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	docs "ai-service/internal/docs"
	"ai-service/internal/domain/entity"
)

type RiskAndGrowthUsecase struct {
	ai        AIProvider
	rag       RiskAndGrowthRepository
	news      NewsRepository
	business  BusinessResearchRepository
	ttl       time.Duration
	publisher MessagePublisher
}

func NewRiskAndGrowthUsecase(
	ai AIProvider,
	rag RiskAndGrowthRepository,
	news NewsRepository,
	business BusinessResearchRepository,
	publisher MessagePublisher,
	ttl time.Duration,
) *RiskAndGrowthUsecase {
	return &RiskAndGrowthUsecase{
		ai:        ai,
		rag:       rag,
		news:      news,
		business:  business,
		ttl:       ttl,
		publisher: publisher,
	}
}

func (u *RiskAndGrowthUsecase) Execute(ctx context.Context, task entity.Task) error {
	logger := slog.With(
		slog.String("id", task.Id),
		slog.String("ticker", task.Ticker),
	)

	logger.Info("starting risk and growth task")

	existing, err := u.rag.GetFreshRiskAndGrowth(ctx, task.Ticker, u.ttl)
	if err != nil {
		return fmt.Errorf("check existing risk and growth: %w", err)
	}

	if existing != nil {
		logger.Info("fresh risk and growth already exists, skipping")
		if err := u.pusblishNextMessage(ctx, task); err != nil {
			return fmt.Errorf("send next msg: %w", err)
		}
		return nil
	}

	news, err := u.news.GetFreshNews(ctx, task.Ticker, u.ttl)
	if err != nil {
		return fmt.Errorf("get fresh news: %w", err)
	}

	business, err := u.business.GetBusinessResearch(ctx, task.Ticker)
	if err != nil {
		return fmt.Errorf("get business research: %w", err)
	}

	newsJSON, err := json.Marshal(news)
	if err != nil {
		return fmt.Errorf("marshal news: %w", err)
	}

	businessJSON, err := json.Marshal(business)
	if err != nil {
		return fmt.Errorf("marshal business: %w", err)
	}

	prompt := docs.RiskAndGrowthPrompt()
	prompt = strings.ReplaceAll(prompt, "{{TICKER}}", task.Ticker)
	prompt = strings.ReplaceAll(prompt, "{{BUSINESS_RESEARCH}}", string(businessJSON))
	prompt = strings.ReplaceAll(prompt, "{{NEWS}}", string(newsJSON))

	factorSchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"name":    {Type: TypeString},
			"type":    {Type: TypeString, Enum: []string{"growth", "risk"}},
			"horizon": {Type: TypeString, Enum: []string{"short_term", "medium_term"}},
			"impact":  {Type: TypeString, Enum: []string{"high", "medium", "low"}},
			"summary": {Type: TypeString},
			"source":  {Type: TypeString},
		},
		Required: []string{"name", "type", "horizon", "impact", "summary", "source"},
	}

	responseSchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"ticker":  {Type: TypeString},
			"factors": {Type: TypeArray, Items: factorSchema},
		},
		Required: []string{"ticker", "factors"},
	}

	logger.Info("calling AI for risk and growth analysis")

	text, err := u.ai.GenerateText(ctx, prompt, entity.Flash, GenerateParams{
		ResponseSchema: responseSchema,
	})
	if err != nil {
		return fmt.Errorf("extract risk and growth: %w", err)
	}

	var result entity.RiskAndGrowthResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		logger.Error("failed to parse risk and growth response", slog.String("ai_response", text))
		return fmt.Errorf("parse risk and growth response: %w", err)
	}

	result.Ticker = task.Ticker

	logger.Info("risk and growth analysis completed",
		slog.Int("factors_count", len(result.Factors)),
	)

	if err := u.rag.SaveRiskAndGrowth(ctx, &result); err != nil {
		return fmt.Errorf("save risk and growth: %w", err)
	}

	logger.Info("risk and growth completed and saved")
	if err := u.pusblishNextMessage(ctx, task); err != nil {
		return fmt.Errorf("send next msg: %w", err)
	}

	return nil
}

func (u *RiskAndGrowthUsecase) pusblishNextMessage(ctx context.Context, task entity.Task) error {
	nextTask := entity.Task{
		Id:     task.Id,
		Type:   entity.RiskAndGrowthSuccess,
		Ticker: task.Ticker,
	}

	marshalled, err := json.Marshal(nextTask)
	if err != nil {
		return fmt.Errorf("marshal risk-and-growth success task: %w", err)
	}

	if err := u.publisher.PublishMessage(ctx, marshalled); err != nil {
		return fmt.Errorf("publish risk-and-growth task: %w", err)
	}

	slog.Info("published risk-and-growth success task")

	return nil
}
