package room

import "github.com/jackc/pgx/v5/pgtype"

type CreateRoomInput struct {
	OwnerID    string
	MaxPlayers *int32
}

type GetRoomInput struct {
	RoomID pgtype.UUID
	UserID string
}

type UpdateRoomInput struct {
	RoomID     pgtype.UUID
	OwnerID    string
	MaxPlayers int32
}

type TransferOwnershipInput struct {
	RoomID     pgtype.UUID
	OwnerID    string
	NewOwnerID string
}

type DeleteRoomInput struct {
	RoomID  pgtype.UUID
	OwnerID string
}

type JoinRoomInput struct {
	RoomID      pgtype.UUID
	UserID      string
	InviteToken string
}

type LeaveRoomInput struct {
	RoomID pgtype.UUID
	UserID string
}

type KickMemberInput struct {
	RoomID       pgtype.UUID
	ActorUserID  string
	TargetUserID string
}

type SelectCharacterInput struct {
	RoomID      pgtype.UUID
	UserID      string
	CharacterID pgtype.UUID
}

type ChangeRoleInput struct {
	RoomID       pgtype.UUID
	ActorUserID  string
	TargetUserID string
	Role         string
}
