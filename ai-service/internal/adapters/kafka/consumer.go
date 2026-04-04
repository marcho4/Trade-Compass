package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"
	kafkalib "github.com/segmentio/kafka-go"
)

const maxRetries = 3

type MessageConsumer interface {
	StartConsuming(ctx context.Context, messages chan<- kafkalib.Message) error
	CommitMessage(ctx context.Context, msg kafkalib.Message) error
}

type Consumer struct {
	kafka      MessageConsumer
	dispatcher *TaskDispatcher
	numWorkers int
	taskChan   chan kafkalib.Message
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

func NewConsumer(kafka MessageConsumer, dispatcher *TaskDispatcher, numWorkers int) *Consumer {
	return &Consumer{
		kafka:      kafka,
		dispatcher: dispatcher,
		numWorkers: numWorkers,
		taskChan:   make(chan kafkalib.Message, numWorkers),
	}
}

func (c *Consumer) Start(ctx context.Context) {
	ctx, c.cancel = context.WithCancel(ctx)

	c.wg.Go(func() {
		defer close(c.taskChan)
		c.consumeWithRetry(ctx)
	})

	for range c.numWorkers {
		c.wg.Go(func() {
			c.worker(ctx)
		})
	}
}

func (c *Consumer) Stop(ctx context.Context) {
	c.cancel()

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("Kafka consumer stopped")
	case <-ctx.Done():
		slog.Warn("Kafka consumer stop timed out")
	}
}

func (c *Consumer) consumeWithRetry(ctx context.Context) {
	const maxBackoff = 30 * time.Second
	backoff := time.Second

	for {
		err := c.kafka.StartConsuming(ctx, c.taskChan)
		if err == nil || errors.Is(err, context.Canceled) {
			return
		}

		slog.Info("Kafka consumer error, retrying",
			slog.Any("error", err),
			slog.Duration("backoff", backoff),
		)

		select {
		case <-time.After(backoff):
			backoff = min(backoff*2, maxBackoff)
		case <-ctx.Done():
			return
		}
	}
}

func (c *Consumer) worker(ctx context.Context) {
	for {
		select {
		case msg, ok := <-c.taskChan:
			if !ok {
				return
			}

			var task entity.Task
			if err := json.Unmarshal(msg.Value, &task); err != nil {
				slog.Error("Failed to unmarshal task",
					slog.String("raw", string(msg.Value)),
					slog.Any("error", err),
				)
				c.commit(ctx, msg)
				continue
			}

			c.processWithRetry(ctx, task, msg)
		case <-ctx.Done():
			return
		}
	}
}

func (c *Consumer) processWithRetry(ctx context.Context, task entity.Task, msg kafkalib.Message) {
	var err error
	backoff := 5 * time.Second

	for attempt := range maxRetries {
		taskCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		err = c.dispatcher.Dispatch(taskCtx, task)
		cancel()

		if errors.Is(err, domain.ErrUnknownTaskType) {
			slog.Warn("Unknown task type", slog.String("type", string(task.Type)))
			c.commit(ctx, msg)
			return
		}

		if err == nil {
			break
		}

		slog.Error("Failed to process task",
			slog.String("type", string(task.Type)),
			slog.String("ticker", task.Ticker),
			slog.Int("attempt", attempt+1),
			slog.Any("error", err),
		)

		if attempt+1 < maxRetries {
			select {
			case <-time.After(backoff):
				backoff *= 2
			case <-ctx.Done():
				return
			}
		}
	}

	if err != nil {
		slog.Error("Task permanently failed, skipping",
			slog.String("type", string(task.Type)),
			slog.String("ticker", task.Ticker),
			slog.Any("error", err),
		)
	}

	c.commit(ctx, msg)
}

func (c *Consumer) commit(ctx context.Context, msg kafkalib.Message) {
	if err := c.kafka.CommitMessage(ctx, msg); err != nil {
		slog.Error("Failed to commit message", slog.Any("error", err))
	}
}
