package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase() (*Database, error) {

	user := os.Getenv("DATABASE_USER")
	if user == "" {
		return nil, fmt.Errorf("DATABASE_USER is not set")
	}

	password := os.Getenv("DATABASE_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("DATABASE_PASSWORD is not set")
	}

	host := os.Getenv("DATABASE_HOST")
	if host == "" {
		return nil, fmt.Errorf("DATABASE_HOST is not set")
	}

	port := os.Getenv("DATABASE_PORT")
	if port == "" {
		return nil, fmt.Errorf("DATABASE_PORT is not set")
	}

	database := os.Getenv("DATABASE_NAME")
	if database == "" {
		return nil, fmt.Errorf("DATABASE_NAME is not set")
	}

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
