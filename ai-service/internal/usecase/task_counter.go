package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"ai-service/internal/domain/entity"
)

const (
	waitTaskCount = 2
)

type TaskCounter interface {
	Increment(ctx context.Context, task entity.Task) error
	Decrement(ctx context.Context, task entity.Task) error
}

type TaskCounterUsecase struct {
	tasks      TasksRepository
	transactor Transactor
	publisher  MessagePublisher
}

func NewTaskCounterUsecase(tasks TasksRepository, transactor Transactor, publisher MessagePublisher) *TaskCounterUsecase {
	return &TaskCounterUsecase{
		tasks:      tasks,
		transactor: transactor,
		publisher:  publisher,
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

func successToExpect(t entity.TaskType) string {
	switch t {
	case entity.RawDataSuccess:
		return string(entity.RawDataExpect)
	case entity.RiskAndGrowthSuccess:
		return string(entity.RiskAndGrowthExpect)
	default:
		return string(t)
	}
}

func (u *TaskCounterUsecase) decrement(ctx context.Context, task entity.Task) error {
	return u.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		_, err := u.tasks.DecrementPending(txCtx, task.Id, successToExpect(task.Type))
		if err != nil {
			return fmt.Errorf("decrement pending: %w", err)
		}

		ok, err := u.tasks.CheckIfTaskIsReady(ctx, task.Id, waitTaskCount)
		if err != nil {
			return fmt.Errorf("check task: %w", err)
		}

		if ok {
			u.sendGenerateScenariosMessage(ctx, task)
		}

		return nil
	})
}

func (u *TaskCounterUsecase) increment(ctx context.Context, task entity.Task) error {
	return u.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		if err := u.tasks.IncrementPending(txCtx, task.Id, string(task.Type), 1); err != nil {
			return fmt.Errorf("increment pending: %w", err)
		}

		ok, err := u.tasks.CheckIfTaskIsReady(ctx, task.Id, waitTaskCount)
		if err != nil {
			return fmt.Errorf("check task: %w", err)
		}

		if ok {
			u.sendGenerateScenariosMessage(ctx, task)
		}

		return nil
	})
}

func (u *TaskCounterUsecase) sendGenerateScenariosMessage(ctx context.Context, task entity.Task) error {
	msg := entity.Task{
		Id:     task.Id,
		Type:   entity.GenerateScenarios,
		Ticker: task.Ticker,
	}

	marshalled, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal risk-and-growth success task: %w", err)
	}

	if err := u.publisher.PublishMessage(ctx, marshalled); err != nil {
		return fmt.Errorf("publish risk-and-growth task: %w", err)
	}

	return nil
}
