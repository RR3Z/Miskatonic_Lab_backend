package room

import (
	"context"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room/helpers"
)

func (s *RoomService) CreateDiceRollRoomEvent(ctx context.Context, input model.CreateDiceRollRoomEventInput) (model.RoomEventModel, error) {
	if err := s.EnsureMember(ctx, input.RoomID, input.ActorID); err != nil {
		return model.RoomEventModel{}, err
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomEventModel{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)

	payload, err := roomHelpers.DiceRollPayload(input.RollID, input.CharacterID, input.Expression, input.Result, input.Details)
	if err != nil {
		return model.RoomEventModel{}, err
	}

	event, err := queries.CreateRoomEvent(ctx, db.CreateRoomEventParams{
		RoomID:    input.RoomID,
		ActorID:   input.ActorID,
		EventType: string(roomEvents.EventDiceRoll),
		Payload:   payload,
	})
	if err != nil {
		return model.RoomEventModel{}, err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomEventModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomEventModel{}, err
	}

	return model.ToRoomEventModel(event), nil
}
