package kafka

import (
	"context"
	"fmt"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"
	"ai-service/internal/usecase"
)

type TaskDispatcher struct {
	handlers map[entity.TaskType]usecase.TaskExecutor
	counter  usecase.TaskCounter
}

func NewTaskDispatcher(
	analyzeReport usecase.TaskExecutor,
	extractRawData usecase.TaskExecutor,
	extractResult usecase.TaskExecutor,
	businessResearch usecase.TaskExecutor,
	newsResearch usecase.TaskExecutor,
	riskAndGrowth usecase.TaskExecutor,
	counter usecase.TaskCounter,
) *TaskDispatcher {
	handlers := map[entity.TaskType]usecase.TaskExecutor{
		entity.Analyze:          analyzeReport,
		entity.Extract:          extractRawData,
		entity.ExtractResult:    extractResult,
		entity.BusinessResearch: businessResearch,
		entity.NewsResearch:     newsResearch,
		entity.RiskAndGrowth:    riskAndGrowth,
	}

	return &TaskDispatcher{
		handlers: handlers,
		counter:  counter,
	}
}

func (d *TaskDispatcher) Dispatch(ctx context.Context, task entity.Task) error {
	switch task.Type {
	case entity.RawDataExpect, entity.RiskAndGrowthExpect:
		return d.counter.Increment(ctx, task)
	case entity.RawDataSuccess, entity.RiskAndGrowthSuccess:
		return d.counter.Decrement(ctx, task)
	}

	handler, ok := d.handlers[task.Type]
	if !ok {
		return fmt.Errorf("%w: %s", domain.ErrUnknownTaskType, task.Type)
	}

	return handler.Execute(ctx, task)
}
