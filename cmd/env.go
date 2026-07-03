package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Warn(
			"env file was not loaded",
			"component", "main",
			"file", ".env",
			"fallback", "system environment variables",
			"error", err,
		)
	}
}

func serverPort() string {
	port := os.Getenv("PORT")
	if port != "" {
		return port
	}

	slog.Warn(
		"server port is not set",
		"component", "main",
		"env", "PORT",
		"default", "8000",
	)
	return "8000"
}
