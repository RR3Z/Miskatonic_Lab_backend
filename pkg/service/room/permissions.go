package room

import (
	"context"
	"errors"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *RoomService) EnsureMember(ctx context.Context, roomID pgtype.UUID, userID string) error {
	_, err := s.repos.Queries.GetMember(ctx, db.GetMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotMember
		}
		return err
	}

	return nil
}

func (s *RoomService) EnsureOwner(ctx context.Context, roomID pgtype.UUID, userID string) error {
	room, err := s.repos.Queries.GetRoomByID(ctx, db.GetRoomByIDParams{
		ID:     roomID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotMember
		}
		return err
	}

	if room.OwnerID != userID {
		return ErrNotOwner
	}

	return nil
}

func (s *RoomService) EnsureCanPublishRoomEvent(ctx context.Context, roomID pgtype.UUID, userID string) error {
	return s.EnsureMember(ctx, roomID, userID)
}
