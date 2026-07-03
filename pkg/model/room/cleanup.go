package roomDTO

import "github.com/jackc/pgx/v5/pgtype"

type CleanupRoomsResult struct {
	InactiveDeleted        int           `json:"inactive_deleted"`
	InvalidDeleted         int           `json:"invalid_deleted"`
	InactiveDeletedRoomIDs []pgtype.UUID `json:"inactive_deleted_room_ids"`
	InvalidDeletedRoomIDs  []pgtype.UUID `json:"invalid_deleted_room_ids"`
	DeletedRoomIDs         []pgtype.UUID `json:"deleted_room_ids"`
}
