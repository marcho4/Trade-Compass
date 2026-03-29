package kafka

import (
	"context"
	"encoding/json"
	"fmt"
)

type CompanyCreatedEvent struct {
	Id     string `json:"id"`
	Ticker string `json:"ticker"`
	Name   string `json:"name"`
}

type AITask struct {
	Id     string `json:"id"`
	Ticker string `json:"ticker"`
	Type   string `json:"type"`
}

type KafkaEventPublisher struct {
	producer   *Producer
	aiProducer *Producer
}

func NewKafkaEventPublisher(producer *Producer, aiProducer *Producer) *KafkaEventPublisher {
	return &KafkaEventPublisher{producer: producer, aiProducer: aiProducer}
}

func (p *KafkaEventPublisher) PublishCompanyCreated(ctx context.Context, ticker, name, id string) error {
	event := CompanyCreatedEvent{
		Ticker: ticker,
		Name:   name,
		Id:     id,
	}

	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal company created event: %w", err)
	}

	return p.producer.Publish(ctx, []byte(ticker), value)
}

func (p *KafkaEventPublisher) PublishBusinessResearchTask(ctx context.Context, ticker, id string) error {
	task := AITask{
		Id:     id,
		Ticker: ticker,
		Type:   "business-research",
	}

	value, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal business research task: %w", err)
	}

	return p.aiProducer.Publish(ctx, []byte(ticker), value)
}

func (p *KafkaEventPublisher) PublishExpectRiskAndGrowthAnalysis(ctx context.Context, ticker, id string) error {
	task := AITask{
		Id:     id,
		Ticker: ticker,
		Type:   "expect-risk-and-growth",
	}

	value, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal business research task: %w", err)
	}

	return p.aiProducer.Publish(ctx, []byte(ticker), value)
}
