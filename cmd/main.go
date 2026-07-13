package main

import (
	"context"
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
