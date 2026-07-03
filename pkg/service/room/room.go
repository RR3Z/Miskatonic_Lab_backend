package room

import (
	"context"
	"time"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/jackc/pgx/v5/pgtype"
)

type IRoom interface {
	CreateRoom(ctx context.Context, input model.CreateRoomInput) (model.RoomModel, error)
	GetRoom(ctx context.Context, input model.GetRoomInput) (model.RoomModel, error)
	UpdateRoom(ctx context.Context, input model.UpdateRoomInput) (model.RoomModel, error)
	TransferOwnership(ctx context.Context, input model.TransferOwnershipInput) (model.RoomModel, error)
	DeleteRoom(ctx context.Context, input model.DeleteRoomInput) error

	JoinRoom(ctx context.Context, input model.JoinRoomInput) (model.RoomMemberModel, error)
	LeaveRoom(ctx context.Context, input model.LeaveRoomInput) error
	KickMember(ctx context.Context, input model.KickMemberInput) error

	SelectCharacter(ctx context.Context, input model.SelectCharacterInput) (model.RoomMemberModel, error)
	ChangeRole(ctx context.Context, input model.ChangeRoleInput) (model.RoomMemberModel, error)
	ListRoomEvents(ctx context.Context, input model.ListRoomEventsInput) ([]model.RoomEventModel, error)
	CreateChatMessage(ctx context.Context, input model.CreateChatMessageInput) (model.RoomEventModel, error)
	CreateDiceRollRoomEvent(ctx context.Context, input model.CreateDiceRollRoomEventInput) (model.RoomEventModel, error)
	TouchRoomActivity(ctx context.Context, input model.TouchRoomActivityInput) error

	EnsureMember(ctx context.Context, roomID pgtype.UUID, userID string) error
	EnsureOwner(ctx context.Context, roomID pgtype.UUID, userID string) error
	EnsureCanPublishRoomEvent(ctx context.Context, roomID pgtype.UUID, userID string) error
}

type IRoomMaintenance interface {
	CleanupRooms(ctx context.Context, input model.CleanupRoomsInput) (model.CleanupRoomsResult, error)
	StartCleanupWorker(ctx context.Context, interval time.Duration)
}
