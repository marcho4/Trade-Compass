package usecase

import (
	"context"

	"ai-service/internal/domain/entity"
)

type TaskExecutor interface {
	Execute(ctx context.Context, task entity.Task) error
}
