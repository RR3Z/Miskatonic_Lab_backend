package handler

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	characterHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character"
	diceRollerHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/diceRoller"
	roomHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/room"
	userHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/user"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	services          *service.Service
	auxiliaryHandlers *AuxiliaryHandlers
}

type AuxiliaryHandlers struct {
	characterHandler  *characterHandler.CharacterHandler
	diceRollerHandler *diceRollerHandler.DiceRollerHandler
	roomHandler       *roomHandler.RoomHandler
	userHandler       *userHandler.UserHandler
}

func NewHandler(services *service.Service) *Handler {
	roomHub := ws.NewRoomHub()
	go roomHub.Run(context.Background())

	var diceHandler *diceRollerHandler.DiceRollerHandler
	if services.Room != nil {
		diceHandler = diceRollerHandler.NewWithRoomChecker(services.DiceRoller, services.Room)
	} else {
		diceHandler = diceRollerHandler.New(services.DiceRoller)
	}

	return &Handler{
		services: services,
		auxiliaryHandlers: new(AuxiliaryHandlers{
			characterHandler:  characterHandler.New(services.Character),
			diceRollerHandler: diceHandler,
			roomHandler:       roomHandler.NewWithHub(services.Room, roomHub),
			userHandler:       userHandler.New(services.User),
		}),
	}
}

func (h *Handler) RoomHub() *ws.RoomHub {
	return h.auxiliaryHandlers.roomHandler.Hub()
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

	h.auxiliaryHandlers.userHandler.RegisterPublicRoutes(router)

	router.Route("/api", func(r chi.Router) {
		if authMiddleware != nil {
			r.Use(authMiddleware)
		}

		h.auxiliaryHandlers.userHandler.RegisterProtectedRoutes(r)

		r.Route("/characters", func(r chi.Router) {
			h.auxiliaryHandlers.characterHandler.RegisterRoutes(r)
		})

		r.Route("/dice-roll", func(r chi.Router) {
			h.auxiliaryHandlers.diceRollerHandler.RegisterRoutes(r)
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
