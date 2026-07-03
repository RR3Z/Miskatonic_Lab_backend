package listeners

import (
	"context"
	"log/slog"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws"
	wsHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/helpers"
	"github.com/jackc/pgx/v5/pgtype"
)

type DiceRollerRoomListener struct {
	roomService room.IRoom
	hub         *ws.RoomHub
	logger      *slog.Logger
}

func NewDiceRollerRoomListener(roomService room.IRoom, hub *ws.RoomHub) *DiceRollerRoomListener {
	return &DiceRollerRoomListener{
		roomService: roomService,
		hub:         hub,
		logger:      slog.Default(),
	}
}

func (l *DiceRollerRoomListener) Handle(ctx context.Context, event events.Event) {
	e, ok := event.(diceEvents.DiceRollMakeSucceeded)
	if !ok {
		return
	}
	if e.RoomID == nil {
		return
	}

	roomUUID := pgtype.UUID{}
	if err := roomUUID.Scan(*e.RoomID); err != nil {
		l.logger.ErrorContext(ctx, "dice room listener: invalid room id",
			"room_id", *e.RoomID,
			"error", err,
		)
		return
	}

	roomEvent, err := l.roomService.CreateDiceRollRoomEvent(ctx, model.CreateDiceRollRoomEventInput{
		RoomID:      roomUUID,
		ActorID:     e.UserID,
		RollID:      e.RollID,
		CharacterID: e.CharacterID,
		Expression:  e.Expression,
		Result:      e.Result,
		Details:     e.Details,
	})
	if err != nil {
		l.logger.ErrorContext(ctx, "dice room listener: failed to create room event",
			"room_id", *e.RoomID,
			"user_id", e.UserID,
			"error", err,
		)
		return
	}

	l.hub.Broadcast(wsHelpers.EventFromRoomEventModel(roomEvent))
}
