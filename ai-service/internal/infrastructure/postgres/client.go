package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBRepo struct {
	conn *pgxpool.Pool
}

func NewDBRepo(ctx context.Context, postgresUrl string) (*DBRepo, error) {
	conn, err := pgxpool.New(ctx, postgresUrl)
	if err != nil {
		slog.Error("Error while creating ")
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		slog.Error("Unable to ping db: %v", err)
		return nil, err
	}

	return &DBRepo{conn: conn}, nil
}

func (d *DBRepo) Close() {
	if d.conn != nil {
		d.conn.Close()
	}
}
