package room

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomModel struct {
	ID          pgtype.UUID        `json:"id"`
	OwnerID     string             `json:"owner_id"`
	MaxPlayers  int32              `json:"max_players"`
	InviteToken string             `json:"invite_token"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	Members     []RoomMemberModel  `json:"members"`
}

func ToRoomModel(r db.Room, members []db.RoomMember) RoomModel {
	m := make([]RoomMemberModel, len(members))
	for i, mb := range members {
		m[i] = ToRoomMemberModel(mb)
	}

	return RoomModel{
		ID:          r.ID,
		OwnerID:     r.OwnerID,
		MaxPlayers:  r.MaxPlayers,
		InviteToken: r.InviteToken,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		Members:     m,
	}
}
