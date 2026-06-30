package room

import "github.com/jackc/pgx/v5/pgtype"

type CreateRoomRequest struct {
	MaxPlayers *int32 `json:"max_players"`
}

type UpdateRoomRequest struct {
	MaxPlayers int32 `json:"max_players"`
}

type SelectCharacterRequest struct {
	CharacterID pgtype.UUID `json:"character_id"`
}

type JoinRoomRequest struct {
	InviteToken string `json:"invite_token"`
}

type ChangeRoleRequest struct {
	Role string `json:"role"`
}

type TransferRoomOwnershipRequest struct {
	UserID string `json:"user_id"`
}
