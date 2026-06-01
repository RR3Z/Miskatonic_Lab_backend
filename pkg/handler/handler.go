package handler

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Handler struct {
	AllowedOrigins []string
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   h.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/health", func(r chi.Router) {
		r.Get("/", h.health)
	})

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
