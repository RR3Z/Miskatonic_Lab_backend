package tests

import (
	"context"

	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/jackc/pgx/v5/pgtype"
)

type fakeListenerRoomService struct {
	diceEvent      roomModel.RoomEventModel
	diceErr        error
	diceCalls      int
	diceInput      roomModel.CreateDiceRollRoomEventInput
	characterEvent []roomModel.RoomEventModel
	characterErr   error
	characterCalls int
	characterInput roomModel.CreateCharacterChangedRoomEventsInput
	memberErr      error
	publishErr     error
}

func (f *fakeListenerRoomService) CreateRoom(_ context.Context, _ roomModel.CreateRoomInput) (roomModel.RoomModel, error) {
	return roomModel.RoomModel{}, nil
}
func (f *fakeListenerRoomService) GetRoom(_ context.Context, _ roomModel.GetRoomInput) (roomModel.RoomModel, error) {
	return roomModel.RoomModel{}, nil
}
func (f *fakeListenerRoomService) UpdateRoom(_ context.Context, _ roomModel.UpdateRoomInput) (roomModel.RoomModel, error) {
	return roomModel.RoomModel{}, nil
}
func (f *fakeListenerRoomService) TransferOwnership(_ context.Context, _ roomModel.TransferOwnershipInput) (roomModel.RoomModel, error) {
	return roomModel.RoomModel{}, nil
}
func (f *fakeListenerRoomService) DeleteRoom(_ context.Context, _ roomModel.DeleteRoomInput) error {
	return nil
}
func (f *fakeListenerRoomService) JoinRoom(_ context.Context, _ roomModel.JoinRoomInput) (roomModel.RoomMemberModel, error) {
	return roomModel.RoomMemberModel{}, nil
}
func (f *fakeListenerRoomService) LeaveRoom(_ context.Context, _ roomModel.LeaveRoomInput) error {
	return nil
}
func (f *fakeListenerRoomService) KickMember(_ context.Context, _ roomModel.KickMemberInput) error {
	return nil
}
func (f *fakeListenerRoomService) SelectCharacter(_ context.Context, _ roomModel.SelectCharacterInput) (roomModel.RoomMemberModel, error) {
	return roomModel.RoomMemberModel{}, nil
}
func (f *fakeListenerRoomService) ChangeRole(_ context.Context, _ roomModel.ChangeRoleInput) (roomModel.RoomMemberModel, error) {
	return roomModel.RoomMemberModel{}, nil
}
func (f *fakeListenerRoomService) ListSelectedCharacters(_ context.Context, _ roomModel.ListSelectedCharactersInput) ([]roomModel.SelectedCharacterModel, error) {
	return nil, nil
}
func (f *fakeListenerRoomService) ListRoomEvents(_ context.Context, _ roomModel.ListRoomEventsInput) ([]roomModel.RoomEventModel, error) {
	return nil, nil
}
func (f *fakeListenerRoomService) CreateChatMessage(_ context.Context, _ roomModel.CreateChatMessageInput) (roomModel.RoomEventModel, error) {
	return roomModel.RoomEventModel{}, nil
}
func (f *fakeListenerRoomService) TouchRoomActivity(_ context.Context, _ roomModel.TouchRoomActivityInput) error {
	return nil
}
func (f *fakeListenerRoomService) EnsureMember(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.memberErr
}
func (f *fakeListenerRoomService) EnsureOwner(_ context.Context, _ pgtype.UUID, _ string) error {
	return nil
}
func (f *fakeListenerRoomService) EnsureCanPublishRoomEvent(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.publishErr
}
func (f *fakeListenerRoomService) CreateDiceRollRoomEvent(_ context.Context, input roomModel.CreateDiceRollRoomEventInput) (roomModel.RoomEventModel, error) {
	f.diceCalls++
	f.diceInput = input
	return f.diceEvent, f.diceErr
}
func (f *fakeListenerRoomService) CreateCharacterChangedRoomEvents(_ context.Context, input roomModel.CreateCharacterChangedRoomEventsInput) ([]roomModel.RoomEventModel, error) {
	f.characterCalls++
	f.characterInput = input
	return f.characterEvent, f.characterErr
}

func listenerTestUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}
