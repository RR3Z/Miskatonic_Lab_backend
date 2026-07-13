package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/publishers"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func connectPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	databaseURL := config.DatabaseURLWithFallback(os.Getenv("DATABASE_URL"), config.PostgresDBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	})

	dbConnection, err := repository.NewPostgresDBFromURL(ctx, databaseURL)
	if err != nil {
		slog.Error(
			"database connection failed",
			"component", "main",
			"database", "postgres",
			"error", err,
		)
		return nil, err
	}

	return dbConnection, nil
}

func configureClerk(ctx context.Context) (func(http.Handler) http.Handler, bool) {
	clerkSecretKey := strings.TrimSpace(os.Getenv("CLERK_SECRET_KEY"))
	if clerkSecretKey == "" {
		slog.Error(
			"clerk secret key is not set",
			"component", "main",
			"env", "CLERK_SECRET_KEY",
		)
		return nil, false
	}

	clerk.SetKey(clerkSecretKey)
	jwksClient, err := middleware.NewClerkJWKSClient(clerkSecretKey)
	if err != nil {
		slog.Error("clerk configuration is invalid", "component", "main", "error", err)
		return nil, false
	}

	preflightCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	keys, err := middleware.PreflightClerkJWKS(preflightCtx, jwksClient)
	if err != nil {
		slog.Error("clerk JWKS preflight failed", "component", "main", "error", err)
		return nil, false
	}
	keyIDs := make([]string, 0, len(keys))
	for _, key := range keys {
		keyIDs = append(keyIDs, key.KeyID)
	}
	slog.Info("clerk JWKS preflight succeeded",
		"component", "main",
		"signing_keys", len(keys),
		"kids", keyIDs,
	)

	authorizedParties := config.ParseAllowedOrigins(os.Getenv("CLERK_AUTHORIZED_PARTIES"))
	if len(authorizedParties) == 0 {
		slog.Error("clerk authorized parties are not set",
			"component", "main",
			"env", "CLERK_AUTHORIZED_PARTIES",
		)
		return nil, false
	}

	return middleware.NewClerkAuthMiddleware(middleware.ClerkAuthConfig{
		JWKSClient:        jwksClient,
		AuthorizedParties: authorizedParties,
		Logger:            slog.Default(),
		Leeway:            middleware.DefaultClerkAuthLeeway,
	}), true
}

func newEventBus(ctx context.Context) *events.EventBus {
	syncPublisher := publishers.NewSyncPublisher()
	asyncPublisher := publishers.NewAsyncPublisher(100, slog.Default())
	asyncPublisher.Start(ctx, 4)

	return events.NewEventBus(syncPublisher, asyncPublisher)
}
