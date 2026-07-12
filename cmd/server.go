package main

import (
	"log/slog"
	"net/http"

	MiskatonicLab "github.com/RR3Z/Miskatonic_Lab_backend"
)

func runHTTPServer(app http.Handler, port string) int {
	server := new(MiskatonicLab.Server)

	slog.Info(
		"http server starting",
		"component", "main",
		"port", port,
	)

	if err := server.Run(port, app); err != nil {
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
