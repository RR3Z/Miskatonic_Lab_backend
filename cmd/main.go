package main

import (
	"log"
	"os"

	MiskatonicLab "github.com/RR3Z/Miskatonic_Lab_backend"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load ENVs
	if err := godotenv.Load(); err != nil {
		log.Println("[ENV] ERROR .env file was not loaded, using system environment variables")
	}

	// Connect Clerk SDK
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		log.Fatal("[CLERK] ERROR because CLERK_SECRET_KEY is not set")
	}
	clerk.SetKey(clerkSecretKey)

	// Configure CORS
	allowedOrigins := config.ParseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if len(allowedOrigins) == 0 {
		log.Fatalf("[main -> CORS] ERROR because CORS_ALLOWED_ORIGINS is not set in .env")
	}
	corsConfig := middleware.CORSConfig{
		AllowedOrigins: allowedOrigins,
	}

	// Launch Server
	handlers := handler.NewHandler(corsConfig)

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		log.Fatal("[HTTP_SERVER] ERROR because PORT is not set (the default port will be used - 8000)")
		serverPort = "8000"
	}

	server := new(MiskatonicLab.Server)

	if err := server.Run(serverPort, handlers.InitRoutes()); err != nil {
		log.Fatalf("[HTTP_SERVER] ERROR while running http server: %v", err)
	}
}
