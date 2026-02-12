package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	reader *kafka.Reader
}

func NewKafkaClient(kafkaUrl, topic string) *KafkaClient {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaUrl},
		GroupID:  "ai-service-group",
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &KafkaClient{reader: reader}
}

func (c *KafkaClient) Close() error {
	err := c.reader.Close()
	if err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}
	return nil
}

func (c *KafkaClient) StartConsuming(ctx context.Context, messages chan<- kafka.Message) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}
		select {
		case messages <- msg:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *KafkaClient) CommitMessage(ctx context.Context, msg kafka.Message) error {
	return c.reader.CommitMessages(ctx, msg)
}
