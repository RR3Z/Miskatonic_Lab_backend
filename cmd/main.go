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
	services := appService.NewService(repos, eventBus)
	appHandlers := handler.NewHandler(services)
	appRouter := appHandlers.InitRoutes(authMiddleware)

	startBackgroundWorkers(ctx, services, appHandlers)
	registerEventListeners(eventBus, services, appHandlers)

	return runHTTPServer(appRouter, serverPort())
}

func main() {
	setupLogger()
	os.Exit(run())
}
