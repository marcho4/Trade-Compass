package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type contextTxKey struct{}

func getContextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, contextTxKey{}, tx)
}

func Executor(ctx context.Context, pool *pgxpool.Pool) Querier {
	if tx, ok := ctx.Value(contextTxKey{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}

func RunInTx[T any](ctx context.Context, db *pgxpool.Pool, fn func(context.Context) (T, error)) (T, error) {
	var zero T

	tx, err := db.Begin(ctx)
	if err != nil {
		return zero, fmt.Errorf("begin tx: %w", err)
	}

	ctx = getContextWithTx(ctx, tx)

	data, err := fn(ctx)
	if err != nil {
		return zero, errors.Join(err, tx.Rollback(ctx))
	}

	if err = tx.Commit(ctx); err != nil {
		return zero, fmt.Errorf("commit tx: %w", err)
	}

	return data, nil
}

type PgxTransactor struct {
	db *pgxpool.Pool
}

func NewPgxTransactor(db *pgxpool.Pool) *PgxTransactor {
	return &PgxTransactor{db: db}
}

func (t *PgxTransactor) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	return RunInTxVoid(ctx, t.db, fn)
}

func RunInTxVoid(ctx context.Context, db *pgxpool.Pool, fn func(context.Context) error) error {
	_, err := RunInTx(ctx, db, func(ctx context.Context) (struct{}, error) {
		return struct{}{}, fn(ctx)
	})
	return err
}
