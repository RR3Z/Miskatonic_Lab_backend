package roomDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomModel struct {
	ID             pgtype.UUID        `json:"id"`
	OwnerID        string             `json:"owner_id"`
	Name           string             `json:"name"`
	MaxPlayers     int32              `json:"max_players"`
	InviteToken    string             `json:"invite_token,omitempty"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
	LastActivityAt pgtype.Timestamptz `json:"last_activity_at"`
	Members        []RoomMemberModel  `json:"members"`
}

func ToRoomModel(r db.Room, members []db.RoomMember, viewerID string) RoomModel {
	m := make([]RoomMemberModel, len(members))
	for i, mb := range members {
		m[i] = ToRoomMemberModel(mb)
	}

	return toRoomModel(r, m, viewerID)
}

func ToRoomModelWithUsernames(r db.Room, members []db.ListMembersByRoomIDRow, viewerID string) RoomModel {
	return toRoomModel(r, ToRoomMemberModels(members), viewerID)
}

func toRoomModel(r db.Room, members []RoomMemberModel, viewerID string) RoomModel {
	model := RoomModel{
		ID:             r.ID,
		OwnerID:        r.OwnerID,
		Name:           r.Name,
		MaxPlayers:     r.MaxPlayers,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
		LastActivityAt: r.LastActivityAt,
		Members:        members,
	}
	if viewerID == r.OwnerID {
		model.InviteToken = r.InviteToken
	}

	return model
}
