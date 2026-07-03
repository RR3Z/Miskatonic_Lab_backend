package listeners

import (
	"context"
	"log/slog"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	listenerHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/listeners/helpers"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/jackc/pgx/v5/pgtype"
)

type CharacterRoomListener struct {
	roomService room.IRoom
	logger      *slog.Logger
}

func NewCharacterRoomListener(roomService room.IRoom) *CharacterRoomListener {
	return &CharacterRoomListener{
		roomService: roomService,
		logger:      slog.Default(),
	}
}

func (l *CharacterRoomListener) Handle(ctx context.Context, event events.Event) {
	actorID, characterID, change, ok := listenerHelpers.CharacterChangedRoomEventInput(event)
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

	if _, err := l.roomService.CreateCharacterChangedRoomEvents(ctx, model.CreateCharacterChangedRoomEventsInput{
		CharacterID: characterUUID,
		ActorID:     actorID,
		Change:      change,
	}); err != nil {
		l.logger.ErrorContext(ctx, "character room listener: failed to create room events",
			"character_id", characterID,
			"user_id", actorID,
			"event", event.EventName(),
			"error", err,
		)
	}
}
