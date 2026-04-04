package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TasksRepository struct {
	db *pgxpool.Pool
}

func NewTasksRepository(db *pgxpool.Pool) *TasksRepository {
	return &TasksRepository{
		db: db,
	}
}

func (r *TasksRepository) IncrementPending(ctx context.Context, taskID string, taskType string, count int) error {
	db := Executor(ctx, r.db)

	sql := `
        INSERT INTO tasks (task_id, task_type, pending_count)
        VALUES ($1, $2, $3)
        ON CONFLICT (task_id, task_type)
        DO UPDATE SET pending_count = tasks.pending_count + $3
    `
	_, err := db.Exec(ctx, sql, taskID, taskType, count)
	if err != nil {
		return fmt.Errorf("increment: %w", err)
	}

	return nil
}

func (r *TasksRepository) DecrementPending(ctx context.Context, taskID string, taskType string) (int, error) {
	db := Executor(ctx, r.db)

	ok, err := r.ensureTaskExists(ctx, taskID, taskType)
	if err != nil {
		return 0, fmt.Errorf("check task exists: %w", err)
	}
	if !ok {
		return 0, errors.New("task not found")
	}

	sql := `update tasks set pending_count = pending_count - 1 where task_id = $1 and task_type = $2 returning pending_count`

	row := db.QueryRow(ctx, sql, taskID, taskType)

	var pendingCount int
	err = row.Scan(&pendingCount)
	if err != nil {
		return 0, fmt.Errorf("update tasks: %w", err)
	}

	return pendingCount, nil
}

func (r *TasksRepository) CheckIfTaskIsReady(ctx context.Context, taskID string, expectedTasks int) (bool, error) {
	db := Executor(ctx, r.db)

	sql := `select count(*) from tasks where task_id = $1 and pending_count = 0`

	rows := db.QueryRow(ctx, sql, taskID)

	var count int
	err := rows.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("select pending count: %w", err)
	}

	return count == expectedTasks, nil
}

func (r *TasksRepository) ensureTaskExists(ctx context.Context, taskID string, taskType string) (bool, error) {
	db := Executor(ctx, r.db)

	sql := `select exists(select 1 from tasks where task_id = $1 and task_type = $2)`

	var exists bool
	err := db.QueryRow(ctx, sql, taskID, taskType).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("exist check: %w", err)
	}

	return exists, nil
}

func (r *TasksRepository) DeleteTask(ctx context.Context, taskID string) error {
	db := Executor(ctx, r.db)

	sql := `delete from tasks where task_id = $1`
	_, err := db.Exec(ctx, sql, taskID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}
