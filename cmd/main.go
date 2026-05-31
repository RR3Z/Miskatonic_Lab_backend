package main

import (
	"log"

	MiskatonicLab "github.com/RR3Z/Miskatonic_Lab_backend"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)

	server := new(MiskatonicLab.Server)
	if err := server.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("[HTTP_SERVER] ERROR while running http server: %v", err)
	}
}
