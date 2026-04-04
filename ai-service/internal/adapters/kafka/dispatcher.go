package kafka

import (
	"context"
	"fmt"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"
)

type TaskExecutor interface {
	Execute(ctx context.Context, task entity.Task) error
}

type TaskDispatcher struct {
	handlers map[entity.TaskType]TaskExecutor
}

func NewTaskDispatcher(
	analyzeReport TaskExecutor,
	extractRawData TaskExecutor,
	extractResult TaskExecutor,
	businessResearch TaskExecutor,
	newsResearch TaskExecutor,
	riskAndGrowth TaskExecutor,
	counter TaskExecutor,
) *TaskDispatcher {
	handlers := map[entity.TaskType]TaskExecutor{
		entity.Analyze:              analyzeReport,
		entity.Extract:              extractRawData,
		entity.ExtractResult:        extractResult,
		entity.BusinessResearch:     businessResearch,
		entity.NewsResearch:         newsResearch,
		entity.RiskAndGrowth:        riskAndGrowth,
		entity.RiskAndGrowthSuccess: counter,
		entity.RiskAndGrowthExpect:  counter,
		entity.RawDataExpect:        counter,
		entity.RawDataSuccess:       counter,
	}

	return &TaskDispatcher{
		handlers: handlers,
	}
}

func (d *TaskDispatcher) Dispatch(ctx context.Context, task entity.Task) error {
	handler, ok := d.handlers[task.Type]
	if !ok {
		return fmt.Errorf("%w: %s", domain.ErrUnknownTaskType, task.Type)
	}

	return handler.Execute(ctx, task)
}
