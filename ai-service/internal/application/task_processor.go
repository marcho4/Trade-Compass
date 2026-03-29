package application

import (
	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"
	"ai-service/internal/infrastructure/financialdata"
	kafkaclient "ai-service/internal/infrastructure/kafka"
	"ai-service/internal/infrastructure/parser"
	"ai-service/internal/infrastructure/postgres"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type TaskProcessor struct {
	taskChan      chan kafka.Message
	numWorkers    int
	kafkaClient   *kafkaclient.KafkaClient
	fdClient      *financialdata.Client
	parserClient  *parser.Client
	geminiService domain.GeminiService
	dbRepo        *postgres.DBRepo
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewTaskProcessor(
	numWorkers int,
	geminiService domain.GeminiService,
	kafkaClient *kafkaclient.KafkaClient,
	fdClient *financialdata.Client,
	parserClient *parser.Client,
	dbRepo *postgres.DBRepo,
) *TaskProcessor {
	taskChan := make(chan kafka.Message, numWorkers)

	return &TaskProcessor{
		taskChan:      taskChan,
		numWorkers:    numWorkers,
		kafkaClient:   kafkaClient,
		geminiService: geminiService,
		parserClient:  parserClient,
		dbRepo:        dbRepo,
		fdClient:      fdClient,
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
			var task entity.Task
			if err := json.Unmarshal(msg.Value, &task); err != nil {
				slog.Error("Failed to unmarshal task", slog.String("raw", string(msg.Value)), slog.Any("error", err))
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

func (p *TaskProcessor) processTask(ctx context.Context, task entity.Task, msg kafka.Message) {
	const maxRetries = 3

	slog.Info("Processing task", slog.Any("task", task))

	var err error
	backoff := 5 * time.Second

	for attempt := range maxRetries {
		err = p.dispatchTask(ctx, task)
		if errors.Is(err, domain.ErrUnknownTaskType) {
			slog.Warn("Unknown task type", slog.String("type", string(task.Type)))
			if err := p.kafkaClient.CommitMessage(ctx, msg); err != nil {
				slog.Error("Failed to commit message", slog.Any("error", err))
			}
			return
		}

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
		slog.Error("Failed to commit message", slog.String("ticker", task.Ticker), slog.Any("error", err))
	} else {
		slog.Info("Kafka message committed", slog.String("ticker", task.Ticker))
	}
}

func (p *TaskProcessor) processAnalyzeTask(ctx context.Context, task entity.Task) error {
	slog.Info("Starting analyze task",
		slog.String("ticker", task.Ticker),
		slog.Int("year", task.Year),
		slog.String("period", task.Period),
	)

	result, err := p.geminiService.AnalyzeReport(ctx, task.Ticker, task.ReportURL, task.Year, entity.ReportPeriod(task.Period))
	if err != nil {
		slog.Error("AnalyzeReport failed",
			slog.String("ticker", task.Ticker),
			slog.Any("error", err),
		)
		return err
	}

	slog.Info("AnalyzeReport succeeded",
		slog.String("ticker", task.Ticker),
		slog.Int("result_length", len(result)),
	)

	periodMonths, ok := entity.PeriodToMonths[task.Period]
	if !ok {
		return fmt.Errorf("unknown period: %s", task.Period)
	}

	if err := p.dbRepo.SaveAnalysis(ctx, result, task.Ticker, task.Year, periodMonths); err != nil {
		slog.Error("SaveAnalysis failed",
			slog.String("ticker", task.Ticker),
			slog.Int("year", task.Year),
			slog.Int("period_months", periodMonths),
			slog.Any("error", err),
		)
		return err
	}

	slog.Info("Analysis saved successfully",
		slog.String("ticker", task.Ticker),
		slog.Int("year", task.Year),
		slog.Int("period_months", periodMonths),
	)

	extractResultTask := entity.Task{
		Ticker:    task.Ticker,
		Year:      task.Year,
		Period:    task.Period,
		ReportURL: task.ReportURL,
		Type:      entity.ExtractResult,
	}

	payload, err := json.Marshal(extractResultTask)
	if err != nil {
		return fmt.Errorf("failed to marshal extract-result task: %w", err)
	}

	if err := p.kafkaClient.PublishMessage(ctx, payload); err != nil {
		return fmt.Errorf("failed to publish extract-result task: %w", err)
	}

	slog.Info("Published extract-result task",
		slog.String("ticker", task.Ticker),
		slog.Int("year", task.Year),
		slog.String("period", task.Period),
	)
	return nil
}

var errUnknownTaskType = domain.ErrUnknownTaskType

func (p *TaskProcessor) dispatchTask(ctx context.Context, task entity.Task) error {
	taskCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	switch task.Type {
	case entity.Analyze:
		return p.processAnalyzeTask(taskCtx, task)
	case entity.Extract:
		return p.processExtractTask(taskCtx, task)
	case entity.ExtractResult:
		return p.processExtractResultTask(taskCtx, task)
	case entity.BusinessResearch:
		return p.processBusinessResearchTask(taskCtx, task)
	default:
		return errUnknownTaskType
	}
}

func (p *TaskProcessor) processExtractTask(ctx context.Context, task entity.Task) error {
	period := entity.ReportPeriod(task.Period)

	slog.Info("extracting raw data...",
		slog.String("ticker", task.Ticker),
		slog.Int("year", task.Year),
		slog.String("period", task.Period),
	)

	existing, err := p.fdClient.GetRawData(ctx, task.Ticker, task.Year, period)
	if err != nil {
		return fmt.Errorf("check existing raw data: %w", err)
	}

	if existing == nil {
		slog.Info("raw data not found. Start extracting...",
			slog.String("ticker", task.Ticker),
			slog.Int("year", task.Year),
			slog.String("period", task.Period),
		)

		result, err := p.geminiService.ExtractRawData(ctx, task.Ticker, task.ReportURL, task.Year, period)
		if err != nil {
			return fmt.Errorf("extract raw data: %w", err)
		}
		result.Status = "confirmed"

		slog.Info("saving new raw data...", slog.String("ticker", task.Ticker))
		if err := p.fdClient.SaveDraft(ctx, result); err != nil {
			return fmt.Errorf("save raw data: %w", err)
		}

		slog.Info("raw data saved successfully", slog.String("ticker", task.Ticker))
	}

	periodMonths, ok := entity.PeriodToMonths[task.Period]
	if !ok {
		return fmt.Errorf("unknown period: %s", task.Period)
	}

	latest, err := p.parserClient.IsLatestReport(ctx, task.Ticker, task.Year, periodMonths)
	if err != nil {
		slog.Warn("failed to check if report is latest, skipping analyze",
			slog.String("ticker", task.Ticker),
			slog.Any("error", err),
		)
		return nil
	}

	if !latest {
		slog.Info("report is not the latest, skipping analyze",
			slog.String("ticker", task.Ticker),
			slog.Int("year", task.Year),
			slog.String("period", task.Period),
		)
		return nil
	}

	analyzeTask := entity.Task{
		Ticker:    task.Ticker,
		Year:      task.Year,
		Period:    task.Period,
		ReportURL: task.ReportURL,
		Type:      entity.Analyze,
	}

	payload, err := json.Marshal(analyzeTask)
	if err != nil {
		return fmt.Errorf("marshal analyze task: %w", err)
	}

	if err := p.kafkaClient.PublishMessage(ctx, payload); err != nil {
		return fmt.Errorf("publish analyze task: %w", err)
	}

	slog.Info("report is the latest, published analyze task",
		slog.String("ticker", task.Ticker),
		slog.Int("year", task.Year),
		slog.String("period", task.Period),
	)

	return nil
}

func (p *TaskProcessor) processBusinessResearchTask(ctx context.Context, task entity.Task) error {
	slog.Info("Starting business research task", slog.String("ticker", task.Ticker))

	existing, err := p.dbRepo.GetBusinessResearch(ctx, task.Ticker)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("check existing business research: %w", err)
	}

	if existing != nil {
		slog.Info("Business research already exists, skipping", slog.String("ticker", task.Ticker))
		return nil
	}

	result, err := p.geminiService.ResearchBusiness(ctx, task.Ticker, task.Ticker)
	if err != nil {
		return fmt.Errorf("research business: %w", err)
	}

	if err := p.dbRepo.SaveBusinessResearch(ctx, result); err != nil {
		return fmt.Errorf("save business research: %w", err)
	}

	slog.Info("Business research completed and saved", slog.String("ticker", task.Ticker))
	return nil
}

func (p *TaskProcessor) processExtractResultTask(ctx context.Context, task entity.Task) error {
	result, err := p.geminiService.ExtractResultFromReport(ctx, task.Ticker, task.Year, entity.ReportPeriod(task.Period))
	if err != nil {
		return fmt.Errorf("extract results from report: %w", err)
	}

	slog.Info("Successfully extracted results from report analysis")

	periodMonths, ok := entity.PeriodToMonths[task.Period]
	if !ok {
		return fmt.Errorf("unknown period: %s", task.Period)
	}

	return p.dbRepo.SaveReportResults(ctx, result, task.Ticker, task.Year, periodMonths)
}

func (p *TaskProcessor) processNewsResearchTask(ctx context.Context, task entity.Task) error {
	return nil
}

func (p *TaskProcessor) processRiskAndGrowthTask(ctx context.Context, task entity.Task) error {
	return nil
}
