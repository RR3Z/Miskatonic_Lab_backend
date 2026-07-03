package roomHelpers

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func RoomIDsFromRooms(rooms []db.Room) []pgtype.UUID {
	ids := make([]pgtype.UUID, 0, len(rooms))
	for _, room := range rooms {
		ids = append(ids, room.ID)
	}
	return ids
}

func AppendRoomIDs(groups ...[]pgtype.UUID) []pgtype.UUID {
	count := 0
	for _, group := range groups {
		count += len(group)
	}

	ids := make([]pgtype.UUID, 0, count)
	for _, group := range groups {
		ids = append(ids, group...)
	}
	return ids
}
