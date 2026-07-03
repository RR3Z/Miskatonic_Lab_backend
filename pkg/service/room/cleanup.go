package room

import (
	"context"
	"log/slog"
	"time"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room/helpers"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	ROOM_INACTIVITY_TTL           = 12 * time.Hour
	DEFAULT_ROOM_CLEANUP_INTERVAL = 15 * time.Minute
)

func (s *RoomService) CleanupRooms(ctx context.Context, input model.CleanupRoomsInput) (model.CleanupRoomsResult, error) {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	inactiveRooms, err := s.repos.Queries.DeleteInactiveRooms(ctx, pgtype.Timestamptz{
		Time:  now.Add(-ROOM_INACTIVITY_TTL),
		Valid: true,
	})
	if err != nil {
		return model.CleanupRoomsResult{}, err
	}

	invalidRooms, err := s.repos.Queries.DeleteInvalidRooms(ctx)
	if err != nil {
		return model.CleanupRoomsResult{}, err
	}

	inactiveRoomIDs := roomHelpers.RoomIDsFromRooms(inactiveRooms)
	invalidRoomIDs := roomHelpers.RoomIDsFromRooms(invalidRooms)

	return model.CleanupRoomsResult{
		InactiveDeleted:        len(inactiveRooms),
		InvalidDeleted:         len(invalidRooms),
		InactiveDeletedRoomIDs: inactiveRoomIDs,
		InvalidDeletedRoomIDs:  invalidRoomIDs,
		DeletedRoomIDs:         roomHelpers.AppendRoomIDs(inactiveRoomIDs, invalidRoomIDs),
	}, nil
}

func (s *RoomService) StartCleanupWorker(ctx context.Context, interval time.Duration, afterCleanup func(model.CleanupRoomsResult)) {
	if interval <= 0 {
		interval = DEFAULT_ROOM_CLEANUP_INTERVAL
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				result, err := s.CleanupRooms(ctx, model.CleanupRoomsInput{})
				if err != nil {
					slog.Warn("room cleanup failed", "component", "room_cleanup", "error", err)
					continue
				}
				if result.InactiveDeleted > 0 || result.InvalidDeleted > 0 {
					slog.Info(
						"room cleanup deleted rooms",
						"component", "room_cleanup",
						"inactive_deleted", result.InactiveDeleted,
						"invalid_deleted", result.InvalidDeleted,
					)
				}
				if afterCleanup != nil {
					afterCleanup(result)
				}
			}
		}
	}()
}
