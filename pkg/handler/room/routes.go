package room

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/httpAdapter"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws"
	"github.com/go-chi/chi/v5"
)

type RoomHandler struct {
	service    roomService.IRoom
	hub        *ws.RoomHub
	dispatcher *ws.CommandDispatcher
}

func New(service roomService.IRoom) *RoomHandler {
	hub := ws.NewRoomHub()
	go hub.Run(context.Background())
	return &RoomHandler{
		service:    service,
		hub:        hub,
		dispatcher: ws.NewCommandDispatcher(service),
	}
}

func (h *RoomHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", httpAdapter.AppHandler(h.createRoom).ServeHTTP)

	r.Route("/{roomID}", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getRoom).ServeHTTP)
		r.Put("/", httpAdapter.AppHandler(h.updateRoom).ServeHTTP)
		r.Delete("/", httpAdapter.AppHandler(h.deleteRoom).ServeHTTP)
		r.Get("/events", httpAdapter.AppHandler(h.listRoomEvents).ServeHTTP)
		r.Get("/ws", h.serveRoomWS)
		r.Put("/owner", httpAdapter.AppHandler(h.transferRoomOwnership).ServeHTTP)

		r.Post("/join", httpAdapter.AppHandler(h.joinRoom).ServeHTTP)
		r.Delete("/leave", httpAdapter.AppHandler(h.leaveRoom).ServeHTTP)
		r.Delete("/kick/{userID}", httpAdapter.AppHandler(h.kickMember).ServeHTTP)

		r.Put("/character", httpAdapter.AppHandler(h.selectCharacter).ServeHTTP)
		r.Put("/members/{userID}/role", httpAdapter.AppHandler(h.changeRole).ServeHTTP)
	})
}
