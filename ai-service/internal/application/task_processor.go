package application

import (
	"ai-service/internal/domain"
	kafkaclient "ai-service/internal/infrastructure/kafka"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type TaskProcessor struct {
	taskChan      chan kafka.Message
	numWorkers    int
	kafkaClient   *kafkaclient.KafkaClient
	geminiService domain.GeminiService
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewTaskProcessor(numWorkers int, geminiService domain.GeminiService, kafkaClient *kafkaclient.KafkaClient) *TaskProcessor {
	taskChan := make(chan kafka.Message, numWorkers)

	return &TaskProcessor{
		taskChan:      taskChan,
		numWorkers:    numWorkers,
		kafkaClient:   kafkaClient,
		geminiService: geminiService,
	}
}

func (p *TaskProcessor) Start(ctx context.Context) {
	ctx, p.cancel = context.WithCancel(ctx)

	p.wg.Go(func() {
		defer close(p.taskChan)
		p.consumeWithRetry(ctx)
	})

	for range p.numWorkers {
		p.wg.Go(func() {
			p.worker(ctx)
		})
	}
}

func (p *TaskProcessor) Stop(ctx context.Context) {
	p.cancel()

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("Task processor stopped")
	case <-ctx.Done():
		slog.Info("Task processor stop timed out, forcing shutdown")
	}
}

func (p *TaskProcessor) consumeWithRetry(ctx context.Context) {
	const maxBackoff = 30 * time.Second
	backoff := time.Second

	for {
		err := p.kafkaClient.StartConsuming(ctx, p.taskChan)
		if err == nil || errors.Is(err, context.Canceled) {
			return
		}

		slog.Info("Kafka consumer error. Retrying...", slog.Any("error", err), slog.Any("backoff", backoff))

		select {
		case <-time.After(backoff):
			backoff = min(backoff*2, maxBackoff)
		case <-ctx.Done():
			return
		}
	}
}

func (p *TaskProcessor) worker(ctx context.Context) {
	for {
		select {
		case msg, ok := <-p.taskChan:
			if !ok {
				return
			}
			var task domain.Task

			if err := json.Unmarshal(msg.Value, &task); err != nil {
				slog.Error("Failed to unmarshal task", slog.Any("error", err))
				if err := p.kafkaClient.CommitMessage(ctx, msg); err != nil {
					slog.Error("Failed to commit bad message", slog.Any("error", err))
				}
				continue
			}
			p.processTask(ctx, task, msg)
		case <-ctx.Done():
			return
		}
	}
}

func (p *TaskProcessor) processTask(ctx context.Context, task domain.Task, msg kafka.Message) {
	const maxRetries = 3

	slog.Info("Processing task", slog.Any("task", task))

	var err error
	backoff := 5 * time.Second

	for attempt := range maxRetries {
		taskCtx, cancel := context.WithTimeout(ctx, 120*time.Second)

		switch task.Type {
		case domain.Analyze:
			err = p.processAnalyzeTask(taskCtx, task)
		case domain.Extract:
			err = p.processExtractTask(taskCtx, task)
		default:
			slog.Warn("Unknown task type", slog.String("type", string(task.Type)))
			cancel()
			if err := p.kafkaClient.CommitMessage(ctx, msg); err != nil {
				slog.Error("Failed to commit message", slog.Any("error", err))
			}
			return
		}
		cancel()

		if err == nil {
			break
		}

		slog.Error("Failed to process task",
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

	if err := p.kafkaClient.CommitMessage(ctx, msg); err != nil {
		slog.Error("Failed to commit message", slog.Any("error", err))
	}
}

func (p *TaskProcessor) processAnalyzeTask(ctx context.Context, task domain.Task) error {
	_, err := p.geminiService.AnalyzeReport(ctx, task.Ticker, task.ReportURL, task.Year, domain.ReportPeriod(task.Period))
	return err
}

func (p *TaskProcessor) processExtractTask(ctx context.Context, task domain.Task) error {
	_, err := p.geminiService.ExtractDataFromReport(ctx, task.Ticker, task.ReportURL, task.Year, domain.ReportPeriod(task.Period))
	return err
}
