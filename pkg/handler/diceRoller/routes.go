package diceRoller

import (
	"context"

	httpAdapter "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/httpadapter"
	diceRollerService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomAccessChecker interface {
	EnsureCanPublishRoomEvent(ctx context.Context, roomID pgtype.UUID, userID string) error
}

type DiceRollerHandler struct {
	service     diceRollerService.IDiceRoller
	roomChecker RoomAccessChecker
}

func New(service diceRollerService.IDiceRoller) *DiceRollerHandler {
	return &DiceRollerHandler{service: service}
}

func NewWithRoomChecker(service diceRollerService.IDiceRoller, checker RoomAccessChecker) *DiceRollerHandler {
	return &DiceRollerHandler{service: service, roomChecker: checker}
}

func (h *DiceRollerHandler) RegisterRoutes(r chi.Router) {
	r.Route("/{characterID}", func(r chi.Router) {
		r.Post("/", httpAdapter.AppHandler(h.makeRoll).ServeHTTP)
		r.Get("/lasts", httpAdapter.AppHandler(h.getLastDiceRolls).ServeHTTP)
	})
}
