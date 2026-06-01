package main

import (
	"log"

	MiskatonicLab "github.com/RR3Z/Miskatonic_Lab_backend"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[ENV] .env file was not loaded, using system environment variables")
	}

	handlers := new(handler.Handler)

	server := new(MiskatonicLab.Server)
	if err := server.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("[HTTP_SERVER] ERROR while running http server: %v", err)
	}
}
