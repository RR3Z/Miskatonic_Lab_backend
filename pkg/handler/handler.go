package handler

import (
	"github.com/go-chi/chi/v5"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/health", func(r chi.Router) {
		r.Get("/", h.health)
	})

	router.Route("/api", func(r chi.Router) {
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
