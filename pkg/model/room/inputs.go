package roomDTO

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateRoomInput struct {
	OwnerID    string
	Name       string
	MaxPlayers *int32
	Password   string
}

type ListRoomsInput struct {
	UserID string
}

type GetRoomInput struct {
	RoomID pgtype.UUID
	UserID string
}

type UpdateRoomInput struct {
	RoomID     pgtype.UUID
	OwnerID    string
	Name       *string
	MaxPlayers int32
	Password   *string
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
	Password    string
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

type CleanupRoomsInput struct {
	Now time.Time
}

type ListRoomEventsInput struct {
	RoomID pgtype.UUID
	UserID string
	Limit  int32
}

type TouchRoomActivityInput struct {
	RoomID pgtype.UUID
	UserID string
}
