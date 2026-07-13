package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/middleware"
	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) InitRoutes(authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()

	allowedOrigins := config.ParseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	router.Use(middleware.CORSMiddleware(middleware.CORSConfig{
		AllowedOrigins: allowedOrigins,
	}))
	router.Use(middleware.RequestLoggingMiddleware(slog.Default()))

	h.registerPublicRoutes(router)
	h.registerProtectedRoutes(router, authMiddleware)

	return router
}

func (h *Handler) registerPublicRoutes(router chi.Router) {
	h.domainHandlers.userHandler.RegisterPublicRoutes(router)
	if h.portraitFileServer == nil {
		return
	}

	router.Method(http.MethodGet, portraitStorage.PublicPathPrefix+"*", h.portraitFileServer)
	router.Method(http.MethodHead, portraitStorage.PublicPathPrefix+"*", h.portraitFileServer)
}

func (h *Handler) registerProtectedRoutes(router chi.Router, authMiddleware func(http.Handler) http.Handler) {
	router.Route("/api", func(r chi.Router) {
		r.Use(authMiddleware)

		h.domainHandlers.userHandler.RegisterProtectedRoutes(r)

		r.Route("/characters", func(r chi.Router) {
			h.domainHandlers.characterHandler.RegisterRoutes(r)
		})

		r.Route("/dice-roll", func(r chi.Router) {
			h.domainHandlers.diceRollerHandler.RegisterRoutes(r)
		})

		r.Route("/rooms", func(r chi.Router) {
			h.domainHandlers.roomHandler.RegisterRoutes(r)
		})
	})
}
