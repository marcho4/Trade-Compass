package port

import (
	"ai-service/internal/domain/entity"
	"context"
	"time"
)

type TasksRepository interface {
	IncrementPending(ctx context.Context, taskID string, taskType string, count int) error
	DecrementPending(ctx context.Context, taskID string, taskType string) error
	GetReadyTasks(ctx context.Context) ([]entity.Task, error)
	DeleteExpired(ctx context.Context, ttl time.Duration) (int64, error)
	DeleteTask(ctx context.Context, taskID string, taskType string) error
}
