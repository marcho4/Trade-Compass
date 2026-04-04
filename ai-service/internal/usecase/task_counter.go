package usecase

import (
	"context"
	"errors"
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

func (u *TaskCounterUsecase) Execute(ctx context.Context, task entity.Task) error {
	switch task.Type {
	case entity.RawDataExpect, entity.RiskAndGrowthExpect:
		return u.increment(ctx, task)
	case entity.RawDataSuccess, entity.RiskAndGrowthSuccess:
		return u.decrement(ctx, task)
	default:
		return errors.New("wrong type of task was processed in TaskCounter usecase")
	}
}

func (u *TaskCounterUsecase) decrement(ctx context.Context, task entity.Task) error {
	return u.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		_, err := u.tasks.DecrementPending(txCtx, task.Id, string(task.Type))
		if err != nil {
			return fmt.Errorf("decrement pending: %w", err)
		}
		return nil
	})
}

func (u *TaskCounterUsecase) increment(ctx context.Context, task entity.Task) error {
	return u.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		if err := u.tasks.IncrementPending(txCtx, task.Id, string(task.Type), 1); err != nil {
			return fmt.Errorf("increment pending: %w", err)
		}
		return nil
	})
}
