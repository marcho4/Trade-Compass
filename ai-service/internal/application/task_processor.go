package application

import (
	"ai-service/internal/domain"
	kafkaclient "ai-service/internal/infrastructure/kafka"
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type TaskProcessor struct {
	taskChan    chan kafka.Message
	numWorkers  int
	kafkaClient *kafkaclient.KafkaClient
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

func NewTaskProcessor(numWorkers int, kafkaClient *kafkaclient.KafkaClient) *TaskProcessor {
	taskChan := make(chan kafka.Message, numWorkers)
	return &TaskProcessor{taskChan: taskChan, numWorkers: numWorkers, kafkaClient: kafkaClient}
}

func (p *TaskProcessor) Start(ctx context.Context) {
	ctx, p.cancel = context.WithCancel(ctx)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(p.taskChan)
		p.consumeWithRetry(ctx)
	}()

	for range p.numWorkers {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.worker(ctx)
		}()
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
		log.Println("Task processor stopped")
	case <-ctx.Done():
		log.Println("Task processor stop timed out, forcing shutdown")
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

		log.Printf("Kafka consumer error: %v, retrying in %v", err, backoff)

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
			var task domain.AnalyzeTask
			if err := json.Unmarshal(msg.Value, &task); err != nil {
				log.Printf("Failed to unmarshal task: %v", err)
				if err := p.kafkaClient.CommitMessage(ctx, msg); err != nil {
					log.Printf("Failed to commit bad message: %v", err)
				}
				continue
			}
			p.processTask(ctx, task, msg)
		case <-ctx.Done():
			return
		}
	}
}

func (p *TaskProcessor) processTask(ctx context.Context, task domain.AnalyzeTask, msg kafka.Message) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	log.Printf("Processing task: %+v\n", task)
	if err := p.kafkaClient.CommitMessage(ctx, msg); err != nil {
		log.Printf("Failed to commit message: %v", err)
	}
}
