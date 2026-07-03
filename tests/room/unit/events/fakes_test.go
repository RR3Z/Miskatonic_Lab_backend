package tests

import (
	"context"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/jackc/pgx/v5/pgtype"
)

type fakeRoomEventPublisher struct {
	events []events.Event
}

func (f *fakeRoomEventPublisher) Publish(_ context.Context, event events.Event) {
	f.events = append(f.events, event)
}

type fakeEventPublishingRoomService struct {
	err                error
	room               model.RoomModel
	member             model.RoomMemberModel
	leaveResult        model.LeaveRoomResult
	selectedCharacters []model.SelectedCharacterModel
	roomEvents         []model.RoomEventModel
	roomEvent          model.RoomEventModel
	cleanupResult      model.CleanupRoomsResult
}

func (f *fakeEventPublishingRoomService) CreateRoom(_ context.Context, _ model.CreateRoomInput) (model.RoomModel, error) {
	return f.room, f.err
}

func (f *fakeEventPublishingRoomService) GetRoom(_ context.Context, _ model.GetRoomInput) (model.RoomModel, error) {
	return f.room, f.err
}

func (f *fakeEventPublishingRoomService) UpdateRoom(_ context.Context, _ model.UpdateRoomInput) (model.RoomModel, error) {
	return f.room, f.err
}

func (f *fakeEventPublishingRoomService) TransferOwnership(_ context.Context, _ model.TransferOwnershipInput) (model.RoomModel, error) {
	return f.room, f.err
}

func (f *fakeEventPublishingRoomService) DeleteRoom(_ context.Context, _ model.DeleteRoomInput) error {
	return f.err
}

func (f *fakeEventPublishingRoomService) JoinRoom(_ context.Context, _ model.JoinRoomInput) (model.RoomMemberModel, error) {
	return f.member, f.err
}

func (f *fakeEventPublishingRoomService) LeaveRoom(_ context.Context, _ model.LeaveRoomInput) (model.LeaveRoomResult, error) {
	return f.leaveResult, f.err
}

func (f *fakeEventPublishingRoomService) KickMember(_ context.Context, _ model.KickMemberInput) error {
	return f.err
}

func (f *fakeEventPublishingRoomService) SelectCharacter(_ context.Context, _ model.SelectCharacterInput) (model.RoomMemberModel, error) {
	return f.member, f.err
}

func (f *fakeEventPublishingRoomService) ChangeRole(_ context.Context, _ model.ChangeRoleInput) (model.RoomMemberModel, error) {
	return f.member, f.err
}

func (f *fakeEventPublishingRoomService) ListSelectedCharacters(_ context.Context, _ model.ListSelectedCharactersInput) ([]model.SelectedCharacterModel, error) {
	return f.selectedCharacters, f.err
}

func (f *fakeEventPublishingRoomService) TouchRoomActivity(_ context.Context, _ model.TouchRoomActivityInput) error {
	return f.err
}

func (f *fakeEventPublishingRoomService) EnsureMember(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.err
}

func (f *fakeEventPublishingRoomService) EnsureOwner(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.err
}

func (f *fakeEventPublishingRoomService) EnsureCanPublishRoomEvent(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.err
}

func (f *fakeEventPublishingRoomService) ListRoomEvents(_ context.Context, _ model.ListRoomEventsInput) ([]model.RoomEventModel, error) {
	return f.roomEvents, f.err
}

func (f *fakeEventPublishingRoomService) CreateChatMessage(_ context.Context, _ model.CreateChatMessageInput) (model.RoomEventModel, error) {
	return f.roomEvent, f.err
}

func (f *fakeEventPublishingRoomService) CreateDiceRollRoomEvent(_ context.Context, _ model.CreateDiceRollRoomEventInput) (model.RoomEventModel, error) {
	return f.roomEvent, f.err
}

func (f *fakeEventPublishingRoomService) CreateCharacterChangedRoomEvents(_ context.Context, _ model.CreateCharacterChangedRoomEventsInput) ([]model.RoomEventModel, error) {
	return f.roomEvents, f.err
}

func (f *fakeEventPublishingRoomService) CleanupRooms(_ context.Context, _ model.CleanupRoomsInput) (model.CleanupRoomsResult, error) {
	return f.cleanupResult, f.err
}

func (f *fakeEventPublishingRoomService) StartCleanupWorker(_ context.Context, _ time.Duration, _ func(model.CleanupRoomsResult)) {
}
