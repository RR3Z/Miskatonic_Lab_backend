package roomDTO

import "github.com/jackc/pgx/v5/pgtype"

type CleanupRoomsResult struct {
	InvalidDeleted        int           `json:"invalid_deleted"`
	InvalidDeletedRoomIDs []pgtype.UUID `json:"invalid_deleted_room_ids"`
	DeletedRoomIDs        []pgtype.UUID `json:"deleted_room_ids"`
}

type StartupPurgeRoomsResult struct {
	DeletedRoomIDs []pgtype.UUID
}
