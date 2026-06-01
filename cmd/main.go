package main

import (
	"log"
	"os"

	MiskatonicLab "github.com/RR3Z/Miskatonic_Lab_backend"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[ENV] .env file was not loaded, using system environment variables")
	}

	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		log.Fatal("[CLERK] CLERK_SECRET_KEY is not set")
	}
	clerk.SetKey(clerkSecretKey)

	handlers := new(handler.Handler)

	server := new(MiskatonicLab.Server)
	if err := server.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("[HTTP_SERVER] ERROR while running http server: %v", err)
	}
}
