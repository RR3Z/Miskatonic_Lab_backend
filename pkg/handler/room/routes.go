package room

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/httpadapter"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	rooms roomService.IRoom
}

func New(rooms roomService.IRoom) *Handler {
	return &Handler{rooms: rooms}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/", httpadapter.AppHandler(h.createRoom).ServeHTTP)

	r.Route("/{roomID}", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getRoom).ServeHTTP)
		r.Put("/", httpadapter.AppHandler(h.updateRoom).ServeHTTP)
		r.Delete("/", httpadapter.AppHandler(h.deleteRoom).ServeHTTP)
		r.Put("/owner", httpadapter.AppHandler(h.transferRoomOwnership).ServeHTTP)

		r.Post("/join", httpadapter.AppHandler(h.joinRoom).ServeHTTP)
		r.Delete("/leave", httpadapter.AppHandler(h.leaveRoom).ServeHTTP)
		r.Delete("/kick/{userID}", httpadapter.AppHandler(h.kickMember).ServeHTTP)

		r.Put("/character", httpadapter.AppHandler(h.selectCharacter).ServeHTTP)
		r.Put("/members/{userID}/role", httpadapter.AppHandler(h.changeRole).ServeHTTP)
	})
}
