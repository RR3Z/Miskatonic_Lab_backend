package main

import (
	"context"
	"log/slog"
	"os"

	MiskatonicLab "github.com/RR3Z/Miskatonic_Lab_backend"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	EventsLogging "github.com/RR3Z/Miskatonic_Lab_backend/pkg/observability/logging"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

func run() int {
	// Load ENV
	if err := godotenv.Load(); err != nil {
		slog.Warn(
			"env file was not loaded",
			"component", "main",
			"file", ".env",
			"fallback", "system environment variables",
			"error", err,
		)
	}

	// Connect Postgres
	ctx := context.Background()
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
		return 1
	}
	defer dbConnection.Close()

	// Connect Clerk SDK
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		slog.Error(
			"clerk secret key is not set",
			"component", "main",
			"env", "CLERK_SECRET_KEY",
		)
		return 1
	}
	clerk.SetKey(clerkSecretKey)

	// Configure CORS
	allowedOrigins := config.ParseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if len(allowedOrigins) == 0 {
		slog.Error(
			"cors allowed origins are not set",
			"component", "main",
			"env", "CORS_ALLOWED_ORIGINS",
		)
		return 1
	}
	corsConfig := middleware.CORSConfig{
		AllowedOrigins: allowedOrigins,
	}

	// Logging
	publisher := events.NewSyncPublisher()
	publisher.Subscribe(EventsLogging.NewCharacterEventLogger(slog.Default()))

	// Launch Server
	repos := repository.NewRepository(dbConnection)
	service := service.NewService(repos, publisher)
	handlers := handler.NewHandler(service, corsConfig)

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		slog.Warn(
			"server port is not set",
			"component", "main",
			"env", "PORT",
			"default", "8000",
		)
		serverPort = "8000"
	}

	server := new(MiskatonicLab.Server)

	slog.Info(
		"http server starting",
		"component", "main",
		"port", serverPort,
	)

	if err := server.Run(serverPort, handlers.InitRoutes()); err != nil {
		slog.Error(
			"http server stopped with error",
			"component", "http_server",
			"port", serverPort,
			"error", err,
		)
		return 1
	}

	return 0
}

func setupLogger() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: "15:04:05",
	}))

	slog.SetDefault(logger)
}

func main() {
	setupLogger()
	os.Exit(run())
}
