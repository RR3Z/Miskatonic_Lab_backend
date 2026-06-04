package handler

import (
	"log/slog"
	"os"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	corsConfig middleware.CORSConfig
	services   *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	allowedOrigins := config.ParseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	router.Use(middleware.CORSMiddleware(middleware.CORSConfig{
		AllowedOrigins: allowedOrigins,
	}))

	router.Use(middleware.RequestLoggingMiddleware(slog.Default()))

	router.Post("/webhooks/clerk/user", AppHandler(h.handleUserClerkWebhook).ServeHTTP)

	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/me", AppHandler(h.getUserByID).ServeHTTP)

		r.Route("/characters", func(r chi.Router) {
			r.Post("/", AppHandler(h.createCharacter).ServeHTTP)
			r.Get("/", AppHandler(h.getAllCharacters).ServeHTTP)

			r.Route("/{characterID}", func(r chi.Router) {
				r.Get("/", AppHandler(h.getCharacter).ServeHTTP)
				r.Put("/", AppHandler(h.updateCharacter).ServeHTTP)
				r.Delete("/", AppHandler(h.deleteCharacter).ServeHTTP)

				r.Route("/characteristics", func(r chi.Router) {
					r.Get("/", AppHandler(h.getCharacteristics).ServeHTTP)
					r.Put("/", AppHandler(h.upsertCharacteristics).ServeHTTP)
					r.Delete("/", AppHandler(h.deleteCharacteristics).ServeHTTP)
				})


				r.Route("/health", func(r chi.Router) {
					r.Get("/", AppHandler(h.getHealth).ServeHTTP)
					r.Put("/", AppHandler(h.upsertHealth).ServeHTTP)
					r.Delete("/", AppHandler(h.deleteHealth).ServeHTTP)
				})

				r.Route("/notes", func(r chi.Router) {
					r.Get("/", AppHandler(h.getNotes).ServeHTTP)
					r.Post("/", AppHandler(h.createNote).ServeHTTP)

					r.Route("/{noteID}", func(r chi.Router) {
						r.Get("/", AppHandler(h.getNote).ServeHTTP)
						r.Put("/", AppHandler(h.updateNote).ServeHTTP)
						r.Delete("/", AppHandler(h.deleteNote).ServeHTTP)
					})
				})
			})
		})
	})

	return router
}
