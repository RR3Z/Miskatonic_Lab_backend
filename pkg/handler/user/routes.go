package user

import (
	httpAdapter "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/httpadapter"
	userService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service userService.IUser
}

func New(service userService.IUser) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterPublicRoutes(r chi.Router) {
	r.Post("/webhooks/clerk/user", httpAdapter.AppHandler(h.handleClerkWebhook).ServeHTTP)
}

func (h *UserHandler) RegisterProtectedRoutes(r chi.Router) {
	r.Get("/me", httpAdapter.AppHandler(h.getMe).ServeHTTP)
}
