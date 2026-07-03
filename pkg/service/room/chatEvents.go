package room

import (
	"context"
	"errors"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room/helpers"
	"github.com/jackc/pgx/v5"
	"strings"
)

func (s *RoomService) CreateChatMessage(ctx context.Context, input model.CreateChatMessageInput) (model.RoomEventModel, error) {
	if err := validateChatMessage(input.Text); err != nil {
		return model.RoomEventModel{}, err
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomEventModel{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	if _, err := queries.GetMember(ctx, db.GetMemberParams{
		RoomID: input.RoomID,
		UserID: input.ActorID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomEventModel{}, ErrNotMember
		}
		return model.RoomEventModel{}, err
	}

	payload, err := roomHelpers.ChatMessagePayload(strings.TrimSpace(input.Text))
	if err != nil {
		return model.RoomEventModel{}, err
	}

	event, err := queries.CreateRoomEvent(ctx, db.CreateRoomEventParams{
		RoomID:    input.RoomID,
		ActorID:   input.ActorID,
		EventType: string(model.EventChatMessage),
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
