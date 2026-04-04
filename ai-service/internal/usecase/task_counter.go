package usecase

import (
	"context"
	"fmt"

	"ai-service/internal/domain/entity"
)

type TaskCounter interface {
	Increment(ctx context.Context, task entity.Task) error
	Decrement(ctx context.Context, task entity.Task) error
}

type TaskCounterUsecase struct {
	tasks      TasksRepository
	transactor Transactor
}

func NewTaskCounterUsecase(tasks TasksRepository, transactor Transactor) *TaskCounterUsecase {
	return &TaskCounterUsecase{
		tasks:      tasks,
		transactor: transactor,
	}
}

func (u *TaskCounterUsecase) Increment(ctx context.Context, task entity.Task) error {
	return u.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		if err := u.tasks.IncrementPending(txCtx, task.Id, string(task.Type), 1); err != nil {
			return fmt.Errorf("increment pending: %w", err)
		}
		return nil
	})
}

func (u *TaskCounterUsecase) Decrement(ctx context.Context, task entity.Task) error {
	return u.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		_, err := u.tasks.DecrementPending(txCtx, task.Id, string(task.Type))
		if err != nil {
			return fmt.Errorf("decrement pending: %w", err)
		}
		return nil
	})
}
