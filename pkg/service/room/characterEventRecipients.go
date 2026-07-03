package room

import (
	"context"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *RoomService) characterChangedRoomEventModel(
	ctx context.Context,
	queries *db.Queries,
	event db.RoomEvent,
	characterID pgtype.UUID,
) (model.RoomEventModel, error) {
	recipients, err := queries.ListCharacterChangedRoomEventRecipients(ctx, db.ListCharacterChangedRoomEventRecipientsParams{
		RoomID:      event.RoomID,
		CharacterID: characterID,
	})
	if err != nil {
		return model.RoomEventModel{}, err
	}

	eventModel := model.ToRoomEventModel(event)
	eventModel.TargetUserIDs = recipients
	return eventModel, nil
}
