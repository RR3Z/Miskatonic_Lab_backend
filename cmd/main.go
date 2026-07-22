package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	appService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
)

func run() int {
	loadEnv()

	ctx := context.Background()
	dbConnection, err := connectPostgres(ctx)
	if err != nil {
		return 1
	}
	defer dbConnection.Close()

	authMiddleware, clerkConfigured := configureClerk(ctx)
	if !clerkConfigured {
		return 1
	}

	eventBus := newEventBus(ctx)
	repos := repository.NewRepository(dbConnection)

	portraitModule, err := configurePortraitModule(repos.Queries)
	if err != nil {
		return 1
	}

	services := appService.NewService(repos, eventBus, portraitModule.Store)
	purgeResult, err := services.PurgeEphemeralRooms(ctx)
	if err != nil {
		slog.Error("ephemeral room purge failed", "component", "room_startup", "error", err)
		return 1
	}
	slog.Info("ephemeral rooms purged", "component", "room_startup", "deleted_rooms", len(purgeResult.DeletedRoomIDs))

	appHandlers := handler.NewHandler(handler.Dependencies{
		Services:           services,
		PortraitFileServer: portraitModule.FileServer,
	})
	appRouter := appHandlers.InitRoutes(authMiddleware)

	startBackgroundWorkers(ctx, backgroundWorkerDependencies{
		Services:           services,
		Handlers:           appHandlers,
		PortraitReconciler: portraitModule.Reconciler,
	})
	registerEventListeners(eventBus, services, appHandlers)

	return runHTTPServer(appRouter, serverPort())
}

func main() {
	setupLogger()
	os.Exit(run())
}
