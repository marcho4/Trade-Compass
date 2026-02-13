package kafka

import (
	"context"
	"encoding/json"
	"fmt"
)

type CompanyCreatedEvent struct {
	Ticker string `json:"ticker"`
}

type KafkaEventPublisher struct {
	producer *Producer
}

func NewKafkaEventPublisher(producer *Producer) *KafkaEventPublisher {
	return &KafkaEventPublisher{producer: producer}
}

func (p *KafkaEventPublisher) PublishCompanyCreated(ctx context.Context, ticker string) error {
	event := CompanyCreatedEvent{Ticker: ticker}
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal company created event: %w", err)
	}
	return p.producer.Publish(ctx, []byte(ticker), value)
}
