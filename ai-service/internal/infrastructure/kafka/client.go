package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewKafkaClient(kafkaUrl, consumeTopic string) *KafkaClient {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaUrl},
		GroupID:  "ai-service-group",
		Topic:    consumeTopic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaUrl),
		Topic:    consumeTopic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaClient{reader: reader, writer: writer}
}

func (c *KafkaClient) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}
	if err := c.writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	return nil
}

func (c *KafkaClient) PublishMessage(ctx context.Context, value []byte) error {
	err := c.writer.WriteMessages(ctx, kafka.Message{Value: value})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
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
