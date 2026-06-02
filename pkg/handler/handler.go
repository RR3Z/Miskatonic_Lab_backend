package handler

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	corsConfig middleware.CORSConfig
}

func NewHandler(corsConfig middleware.CORSConfig) *Handler {
	return &Handler{
		corsConfig: corsConfig,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.CORSMiddleware(h.corsConfig))

	router.Post("/webhooks/clerk/user", h.handleUserClerkWebhook)

	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Route("/characters", func(r chi.Router) {
			r.Post("/", h.createCharacter)
			r.Get("/", h.getAllCharacters)

			r.Route("/{characterID}", func(r chi.Router) {
				r.Get("/", h.getCharacter)
				r.Put("/", h.updateCharacter)
				r.Delete("/", h.deleteCharacter)
			})
		})
	})

	return router
}
