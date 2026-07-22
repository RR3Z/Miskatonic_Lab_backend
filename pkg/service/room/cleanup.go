package room

import (
	"context"
	"log/slog"
	"time"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room/helpers"
)

const (
	DEFAULT_ROOM_CLEANUP_INTERVAL = 15 * time.Minute
)

func (s *RoomService) PurgeEphemeralRooms(ctx context.Context) (model.StartupPurgeRoomsResult, error) {
	rooms, err := s.repos.Queries.DeleteAllRooms(ctx)
	if err != nil {
		return model.StartupPurgeRoomsResult{}, err
	}

	return model.StartupPurgeRoomsResult{
		DeletedRoomIDs: roomHelpers.RoomIDsFromRooms(rooms),
	}, nil
}

func (s *RoomService) CleanupRooms(ctx context.Context) (model.CleanupRoomsResult, error) {
	invalidRooms, err := s.repos.Queries.DeleteInvalidRooms(ctx)
	if err != nil {
		return model.CleanupRoomsResult{}, err
	}

	invalidRoomIDs := roomHelpers.RoomIDsFromRooms(invalidRooms)

	return model.CleanupRoomsResult{
		InvalidDeleted:        len(invalidRooms),
		InvalidDeletedRoomIDs: invalidRoomIDs,
		DeletedRoomIDs:        invalidRoomIDs,
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
				result, err := s.CleanupRooms(ctx)
				if err != nil {
					slog.Warn("room cleanup failed", "component", "room_cleanup", "error", err)
					continue
				}
				if result.InvalidDeleted > 0 {
					slog.Info(
						"room cleanup deleted rooms",
						"component", "room_cleanup",
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
