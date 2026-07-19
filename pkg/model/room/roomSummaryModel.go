package roomDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomSummaryModel struct {
	ID          pgtype.UUID        `json:"id"`
	Name        string             `json:"name"`
	MaxPlayers  int32              `json:"max_players"`
	MemberCount int32              `json:"member_count"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	IsMember    bool               `json:"is_member"`
}

func ToRoomSummaryModels(rooms []db.ListRoomsRow) []RoomSummaryModel {
	models := make([]RoomSummaryModel, len(rooms))
	for index, room := range rooms {
		models[index] = RoomSummaryModel{
			ID:          room.ID,
			Name:        room.Name,
			MaxPlayers:  room.MaxPlayers,
			MemberCount: room.MemberCount,
			CreatedAt:   room.CreatedAt,
			IsMember:    room.IsMember,
		}
	}

	return models
}
