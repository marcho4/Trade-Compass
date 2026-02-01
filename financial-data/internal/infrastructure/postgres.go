package infrastructure

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	connString := os.Getenv("DB_URL")
	if connString == "" {
		return nil, fmt.Errorf("DB_URL environment variable is required")
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse postgres connection string: %w", err)
	}

	config.MaxConns = getEnvAsInt("DB_MAX_CONNS", 25)
	config.MinConns = getEnvAsInt("DB_MIN_CONNS", 5)
	config.MaxConnLifetime = time.Duration(getEnvAsInt("DB_MAX_CONN_LIFETIME_MIN", 30)) * time.Minute
	config.MaxConnIdleTime = time.Duration(getEnvAsInt("DB_MAX_CONN_IDLE_TIME_MIN", 5)) * time.Minute
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = 10 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create postgres connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}

func getEnvAsInt(key string, defaultValue int) int32 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return int32(defaultValue)
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return int32(defaultValue)
	}
	return int32(value)
}
