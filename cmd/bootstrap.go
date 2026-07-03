package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/publishers"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	appService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func connectPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	dbConnection, err := repository.NewPostgresDB(ctx, config.PostgresDBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	})
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

func configureClerk() bool {
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		slog.Error(
			"clerk secret key is not set",
			"component", "main",
			"env", "CLERK_SECRET_KEY",
		)
		return false
	}

	clerk.SetKey(clerkSecretKey)
	return true
}

func newEventBus(ctx context.Context) *events.EventBus {
	syncPublisher := publishers.NewSyncPublisher()
	asyncPublisher := publishers.NewAsyncPublisher(100, slog.Default())
	asyncPublisher.Start(ctx, 4)

	return events.NewEventBus(syncPublisher, asyncPublisher)
}

func startBackgroundWorkers(ctx context.Context, services *appService.Service, appHandlers *handler.Handler) {
	services.StartBackgroundWorkers(ctx, appService.BackgroundWorkerHooks{
		RoomCleanup: func(result roomModel.CleanupRoomsResult) {
			appHandlers.CloseDeletedRoomSockets(result, "room deleted by cleanup")
		},
	})
}
