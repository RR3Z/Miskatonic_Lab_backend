package room

import (
	"context"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room/helpers"
)

func (s *RoomService) CreateCharacterChangedRoomEvents(ctx context.Context, input model.CreateCharacterChangedRoomEventsInput) ([]model.RoomEventModel, error) {
	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	roomIDs, err := queries.ListRoomIDsBySelectedCharacter(ctx, input.CharacterID)
	if err != nil {
		return nil, err
	}

	events := make([]model.RoomEventModel, 0, len(roomIDs))
	for _, roomID := range roomIDs {
		payload, err := roomHelpers.CharacterChangedPayload(
			input.CharacterID.String(),
			input.Change.Resource,
			input.Change.Action,
			input.Change.ResourceID,
			input.Change.SourceEvent,
		)
		if err != nil {
			return nil, err
		}

		event, err := queries.CreateRoomEvent(ctx, db.CreateRoomEventParams{
			RoomID:    roomID,
			ActorID:   input.ActorID,
			EventType: string(roomEvents.EventCharacterChanged),
			Payload:   payload,
		})
		if err != nil {
			return nil, err
		}

		if _, err := queries.TouchRoomActivity(ctx, roomID); err != nil {
			return nil, err
		}

		events = append(events, model.ToRoomEventModel(event))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return events, nil
}
