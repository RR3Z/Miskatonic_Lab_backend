package room

import (
	"context"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
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

func (s *RoomService) TouchRoomActivity(ctx context.Context, input model.TouchRoomActivityInput) error {
	if err := s.EnsureMember(ctx, input.RoomID, input.UserID); err != nil {
		return err
	}

	_, err := s.repos.Queries.TouchRoomActivity(ctx, input.RoomID)
	return err
}
