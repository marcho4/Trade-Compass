package subscriptions

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	pg *pgxpool.Pool
}

func NewSubscriptionRepo(pool *pgxpool.Pool) *Repository {
	return &Repository{pg: pool}
}
