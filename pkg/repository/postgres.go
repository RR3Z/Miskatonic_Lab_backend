package repository

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, cfg config.PostgresDBConfig) (*pgxpool.Pool, error) {
	databaseURL := config.DatabaseURL(cfg)

	return NewPostgresDBFromURL(ctx, databaseURL)
}

func NewPostgresDBFromURL(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	dbConnection, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := dbConnection.Ping(ctx); err != nil {
		dbConnection.Close()
		return nil, err
	}

	return dbConnection, nil
}
