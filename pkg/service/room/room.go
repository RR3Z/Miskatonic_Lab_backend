package room

import (
	"context"
	"errors"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

var (
	ErrRoomNotFound      = errors.New("room not found")
	ErrRoomFull          = errors.New("room is full")
	ErrAlreadyMember     = errors.New("already a member of this room")
	ErrNotMember         = errors.New("not a member of this room")
	ErrNotOwner          = errors.New("only the room owner can perform this action")
	ErrCannotKickOwner   = errors.New("cannot kick the room owner")
	ErrCharacterNotOwned = errors.New("character does not belong to you")
)

type IRoom interface {
	CreateRoom(ctx context.Context, params db.CreateRoomParams) (model.RoomModel, error)
	GetRoom(ctx context.Context, params db.GetRoomByIDParams) (model.RoomModel, error)
	UpdateRoom(ctx context.Context, params db.UpdateRoomParams) (model.RoomModel, error)
	TransferOwnership(ctx context.Context, params db.TransferRoomOwnershipParams) (model.RoomModel, error)
	DeleteRoom(ctx context.Context, params db.DeleteRoomParams) error

	JoinRoom(ctx context.Context, meta db.GetRoomMetaDataParams, member db.GetMemberParams) (model.RoomMemberModel, error)
	LeaveRoom(ctx context.Context, params db.RemoveMemberParams) error
	KickMember(ctx context.Context, actor db.GetRoomByIDParams, target db.RemoveMemberParams) error

	SelectCharacter(ctx context.Context, params db.UpdateMemberCharacterParams) (model.RoomMemberModel, error)
	ChangeRole(ctx context.Context, actor db.GetRoomByIDParams, target db.UpdateMemberRoleParams) (model.RoomMemberModel, error)
}
