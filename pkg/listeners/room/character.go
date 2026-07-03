package roomlisteners

import (
	"context"
	"log/slog"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	wsHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/helpers"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
	"github.com/jackc/pgx/v5/pgtype"
)

type CharacterRoomListener struct {
	roomService roomService.IRoom
	hub         *ws.RoomHub
	logger      *slog.Logger
}

func NewCharacterRoomListener(roomService roomService.IRoom, hub *ws.RoomHub) *CharacterRoomListener {
	return &CharacterRoomListener{
		roomService: roomService,
		hub:         hub,
		logger:      slog.Default(),
	}
}

func (l *CharacterRoomListener) Handle(ctx context.Context, event events.Event) {
	actorID, characterID, change, ok := CharacterChangedRoomEventInput(event)
	if !ok {
		return
	}

	characterUUID := pgtype.UUID{}
	if err := characterUUID.Scan(characterID); err != nil {
		l.logger.ErrorContext(ctx, "character room listener: invalid character id",
			"character_id", characterID,
			"event", event.EventName(),
			"error", err,
		)
		return
	}

	roomEvents, err := l.roomService.CreateCharacterChangedRoomEvents(ctx, model.CreateCharacterChangedRoomEventsInput{
		CharacterID: characterUUID,
		ActorID:     actorID,
		Change:      change,
	})
	if err != nil {
		l.logger.ErrorContext(ctx, "character room listener: failed to create room events",
			"character_id", characterID,
			"user_id", actorID,
			"event", event.EventName(),
			"error", err,
		)
		return
	}

	for _, roomEvent := range roomEvents {
		l.hub.SendToUsers(
			roomEvent.RoomID.String(),
			roomEvent.TargetUserIDs,
			wsHelpers.EventFromRoomEventModel(roomEvent),
		)
	}
}
