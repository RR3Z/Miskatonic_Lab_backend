package room

import (
	"context"
	"errors"
	"strings"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room/helpers"
	"github.com/jackc/pgx/v5"
)

func (s *RoomService) ListRoomEvents(ctx context.Context, input model.ListRoomEventsInput) ([]model.RoomEventModel, error) {
	if err := s.EnsureMember(ctx, input.RoomID, input.UserID); err != nil {
		return nil, err
	}

	events, err := s.repos.Queries.ListRoomEvents(ctx, db.ListRoomEventsParams{
		RoomID:     input.RoomID,
		UserID:     input.UserID,
		LimitCount: normalizeRoomEventsLimit(input.Limit),
	})
	if err != nil {
		return nil, err
	}

	return model.ToRoomEventModels(events), nil
}

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
		EventType: string(roomEvents.EventChatMessage),
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

func (s *RoomService) TouchRoomActivity(ctx context.Context, input model.TouchRoomActivityInput) error {
	if err := s.EnsureMember(ctx, input.RoomID, input.UserID); err != nil {
		return err
	}

	_, err := s.repos.Queries.TouchRoomActivity(ctx, input.RoomID)
	return err
}
