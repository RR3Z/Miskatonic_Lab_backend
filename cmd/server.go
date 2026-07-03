package main

import (
	"log/slog"

	MiskatonicLab "github.com/RR3Z/Miskatonic_Lab_backend"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
)

func runHTTPServer(appHandlers *handler.Handler, port string) int {
	server := new(MiskatonicLab.Server)

	slog.Info(
		"http server starting",
		"component", "main",
		"port", port,
	)

	if err := server.Run(port, appHandlers.InitRoutes()); err != nil {
		slog.Error(
			"http server stopped with error",
			"component", "http_server",
			"port", port,
			"error", err,
		)
		return 1
	}

	return 0
}
