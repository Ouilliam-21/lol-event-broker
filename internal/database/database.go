package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(user, password, host, port, database string) (*Database, error) {

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, database)

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Unable to ping database: %w", err)
	}

	return &Database{Pool: pool}, nil
}
