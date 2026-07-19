package room

import (
	"context"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func createMutationEvent(
	ctx context.Context,
	queries *db.Queries,
	roomID pgtype.UUID,
	actorID string,
	eventType model.EventType,
	payload []byte,
) (model.RoomEventModel, error) {
	event, err := queries.CreateRoomEvent(ctx, db.CreateRoomEventParams{
		RoomID:    roomID,
		ActorID:   actorID,
		EventType: string(eventType),
		Payload:   payload,
	})
	if err != nil {
		return model.RoomEventModel{}, err
	}

	return model.ToRoomEventModel(event), nil
}
