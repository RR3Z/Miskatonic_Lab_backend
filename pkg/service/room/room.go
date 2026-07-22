package room

import (
	"context"
	"time"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/jackc/pgx/v5/pgtype"
)

type IRoom interface {
	CreateRoom(ctx context.Context, input model.CreateRoomInput) (model.RoomMutationResult[model.RoomModel], error)
	ListRooms(ctx context.Context, input model.ListRoomsInput) ([]model.RoomSummaryModel, error)
	GetRoom(ctx context.Context, input model.GetRoomInput) (model.RoomModel, error)
	UpdateRoom(ctx context.Context, input model.UpdateRoomInput) (model.RoomMutationResult[model.RoomModel], error)
	TransferOwnership(ctx context.Context, input model.TransferOwnershipInput) (model.RoomMutationResult[model.RoomModel], error)
	DeleteRoom(ctx context.Context, input model.DeleteRoomInput) (model.RoomMutationResult[struct{}], error)

	JoinRoom(ctx context.Context, input model.JoinRoomInput) (model.RoomMutationResult[model.RoomMemberModel], error)
	LeaveRoom(ctx context.Context, input model.LeaveRoomInput) (model.RoomMutationResult[model.LeaveRoomResult], error)
	KickMember(ctx context.Context, input model.KickMemberInput) (model.RoomMutationResult[struct{}], error)

	SelectCharacter(ctx context.Context, input model.SelectCharacterInput) (model.RoomMutationResult[model.RoomMemberModel], error)
	ChangeRole(ctx context.Context, input model.ChangeRoleInput) (model.RoomMutationResult[model.RoomMemberModel], error)
	ListSelectedCharacters(ctx context.Context, input model.ListSelectedCharactersInput) ([]model.SelectedCharacterModel, error)
	TouchRoomActivity(ctx context.Context, input model.TouchRoomActivityInput) error

	EnsureMember(ctx context.Context, roomID pgtype.UUID, userID string) error
	EnsureOwner(ctx context.Context, roomID pgtype.UUID, userID string) error
	EnsureCanPublishRoomEvent(ctx context.Context, roomID pgtype.UUID, userID string) error

	ListRoomEvents(ctx context.Context, input model.ListRoomEventsInput) ([]model.RoomEventModel, error)
	CreateChatMessage(ctx context.Context, input model.CreateChatMessageInput) (model.RoomEventModel, error)
	CreateDiceRollRoomEvent(ctx context.Context, input model.CreateDiceRollRoomEventInput) (model.RoomEventModel, error)
	CreateCharacterChangedRoomEvents(ctx context.Context, input model.CreateCharacterChangedRoomEventsInput) ([]model.RoomEventModel, error)
}

type IRoomMaintenance interface {
	PurgeEphemeralRooms(ctx context.Context) (model.StartupPurgeRoomsResult, error)
	CleanupRooms(ctx context.Context) (model.CleanupRoomsResult, error)
	StartCleanupWorker(ctx context.Context, interval time.Duration, afterCleanup func(model.CleanupRoomsResult))
}
