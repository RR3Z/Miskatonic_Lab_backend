package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	characterHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character"
	roomHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	services          *service.Service
	auxiliaryHandlers *AuxiliaryHandlers
}

type AuxiliaryHandlers struct {
	characterHandler *characterHandler.CharacterHandler
	roomHandler      *roomHandler.RoomHandler
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
		auxiliaryHandlers: new(AuxiliaryHandlers{
			characterHandler: characterHandler.New(services.Character),
			roomHandler:      roomHandler.New(services.Room),
		}),
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	return h.initRoutes(middleware.AuthMiddleware)
}

func (h *Handler) initRoutes(authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()

	allowedOrigins := config.ParseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	router.Use(middleware.CORSMiddleware(middleware.CORSConfig{
		AllowedOrigins: allowedOrigins,
	}))

	router.Use(middleware.RequestLoggingMiddleware(slog.Default()))

	router.Post("/webhooks/clerk/user", AppHandler(h.handleUserClerkWebhook).ServeHTTP)

	router.Route("/api", func(r chi.Router) {
		if authMiddleware != nil {
			r.Use(authMiddleware)
		}

		r.Get("/me", AppHandler(h.getUserByID).ServeHTTP)

		r.Route("/characters", func(r chi.Router) {
			h.auxiliaryHandlers.characterHandler.RegisterRoutes(r)
		})

		r.Route("/dice-roll/{characterID}", func(r chi.Router) {
			r.Post("/", AppHandler(h.makeRoll).ServeHTTP)
			r.Get("/lasts", AppHandler(h.getLastDiceRolls).ServeHTTP)
		})

		r.Route("/rooms", func(r chi.Router) {
			h.auxiliaryHandlers.roomHandler.RegisterRoutes(r)
		})
	})

	return router
}

// FOR TESTS
func (h *Handler) InitRoutesWithAuth(authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	return h.initRoutes(authMiddleware)
}
