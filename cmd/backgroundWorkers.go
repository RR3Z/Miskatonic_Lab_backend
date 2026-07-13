package main

import (
	"context"
	"log/slog"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	portraitMaintenance "github.com/RR3Z/Miskatonic_Lab_backend/pkg/maintenance/portrait"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	appService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
)

type backgroundWorkerDependencies struct {
	Services           *appService.Service
	Handlers           *handler.Handler
	PortraitReconciler portraitMaintenance.ReconciliationRunner
}

func startBackgroundWorkers(ctx context.Context, dependencies backgroundWorkerDependencies) {
	dependencies.Services.StartBackgroundWorkers(ctx, appService.BackgroundWorkerHooks{
		RoomCleanup: func(result roomModel.CleanupRoomsResult) {
			dependencies.Handlers.CloseDeletedRoomSockets(result, "room deleted by cleanup")
		},
	})

	portraitMaintenance.NewWorker(
		dependencies.PortraitReconciler,
		portraitMaintenance.DefaultReconciliationInterval,
		slog.Default(),
	).Start(ctx)
}
