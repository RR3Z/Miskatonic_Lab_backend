package tests

import (
	"context"

	roomModels "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/jackc/pgx/v5/pgtype"
)

type fakeRoomHandlerService struct {
	err error

	room   roomModels.RoomModel
	member roomModels.RoomMemberModel

	createCalls int
	createInput roomModels.CreateRoomInput

	getCalls int
	getInput roomModels.GetRoomInput

	updateCalls int
	updateInput roomModels.UpdateRoomInput

	transferCalls int
	transferInput roomModels.TransferOwnershipInput

	deleteCalls int
	deleteInput roomModels.DeleteRoomInput

	joinCalls int
	joinInput roomModels.JoinRoomInput

	leaveCalls int
	leaveInput roomModels.LeaveRoomInput

	kickCalls int
	kickInput roomModels.KickMemberInput

	selectCharacterCalls int
	selectCharacterInput roomModels.SelectCharacterInput

	changeRoleCalls int
	changeRoleInput roomModels.ChangeRoleInput

	selectedCharacters          []roomModels.SelectedCharacterModel
	listSelectedCharactersCalls int
	listSelectedCharactersInput roomModels.ListSelectedCharactersInput

	events          []roomModels.RoomEventModel
	listEventsCalls int
	listEventsInput roomModels.ListRoomEventsInput

	chatEvent       roomModels.RoomEventModel
	createChatCalls int
	createChatInput roomModels.CreateChatMessageInput

	touchActivityCalls int
	touchActivityInput roomModels.TouchRoomActivityInput
}

func (f *fakeRoomHandlerService) totalCalls() int {
	return f.createCalls + f.getCalls + f.updateCalls + f.transferCalls + f.deleteCalls + f.joinCalls + f.leaveCalls + f.kickCalls + f.selectCharacterCalls + f.changeRoleCalls + f.listSelectedCharactersCalls + f.listEventsCalls + f.createChatCalls
}

func (f *fakeRoomHandlerService) CreateRoom(_ context.Context, input roomModels.CreateRoomInput) (roomModels.RoomModel, error) {
	f.createCalls++
	f.createInput = input
	return f.room, f.err
}

func (f *fakeRoomHandlerService) GetRoom(_ context.Context, input roomModels.GetRoomInput) (roomModels.RoomModel, error) {
	f.getCalls++
	f.getInput = input
	return f.room, f.err
}

func (f *fakeRoomHandlerService) UpdateRoom(_ context.Context, input roomModels.UpdateRoomInput) (roomModels.RoomModel, error) {
	f.updateCalls++
	f.updateInput = input
	return f.room, f.err
}

func (f *fakeRoomHandlerService) TransferOwnership(_ context.Context, input roomModels.TransferOwnershipInput) (roomModels.RoomModel, error) {
	f.transferCalls++
	f.transferInput = input
	return f.room, f.err
}

func (f *fakeRoomHandlerService) DeleteRoom(_ context.Context, input roomModels.DeleteRoomInput) error {
	f.deleteCalls++
	f.deleteInput = input
	return f.err
}

func (f *fakeRoomHandlerService) JoinRoom(_ context.Context, input roomModels.JoinRoomInput) (roomModels.RoomMemberModel, error) {
	f.joinCalls++
	f.joinInput = input
	return f.member, f.err
}

func (f *fakeRoomHandlerService) LeaveRoom(_ context.Context, input roomModels.LeaveRoomInput) error {
	f.leaveCalls++
	f.leaveInput = input
	return f.err
}

func (f *fakeRoomHandlerService) KickMember(_ context.Context, input roomModels.KickMemberInput) error {
	f.kickCalls++
	f.kickInput = input
	return f.err
}

func (f *fakeRoomHandlerService) SelectCharacter(_ context.Context, input roomModels.SelectCharacterInput) (roomModels.RoomMemberModel, error) {
	f.selectCharacterCalls++
	f.selectCharacterInput = input
	return f.member, f.err
}

func (f *fakeRoomHandlerService) ChangeRole(_ context.Context, input roomModels.ChangeRoleInput) (roomModels.RoomMemberModel, error) {
	f.changeRoleCalls++
	f.changeRoleInput = input
	return f.member, f.err
}

func (f *fakeRoomHandlerService) ListSelectedCharacters(_ context.Context, input roomModels.ListSelectedCharactersInput) ([]roomModels.SelectedCharacterModel, error) {
	f.listSelectedCharactersCalls++
	f.listSelectedCharactersInput = input
	return f.selectedCharacters, f.err
}

func (f *fakeRoomHandlerService) ListRoomEvents(_ context.Context, input roomModels.ListRoomEventsInput) ([]roomModels.RoomEventModel, error) {
	f.listEventsCalls++
	f.listEventsInput = input
	return f.events, f.err
}

func (f *fakeRoomHandlerService) CreateChatMessage(_ context.Context, input roomModels.CreateChatMessageInput) (roomModels.RoomEventModel, error) {
	f.createChatCalls++
	f.createChatInput = input
	return f.chatEvent, f.err
}

func (f *fakeRoomHandlerService) TouchRoomActivity(_ context.Context, input roomModels.TouchRoomActivityInput) error {
	f.touchActivityCalls++
	f.touchActivityInput = input
	return f.err
}

func (f *fakeRoomHandlerService) EnsureMember(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.err
}

func (f *fakeRoomHandlerService) EnsureOwner(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.err
}

func (f *fakeRoomHandlerService) EnsureCanPublishRoomEvent(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.err
}

func (f *fakeRoomHandlerService) CreateDiceRollRoomEvent(_ context.Context, _ roomModels.CreateDiceRollRoomEventInput) (roomModels.RoomEventModel, error) {
	return roomModels.RoomEventModel{}, f.err
}

func (f *fakeRoomHandlerService) CreateCharacterChangedRoomEvents(_ context.Context, _ roomModels.CreateCharacterChangedRoomEventsInput) ([]roomModels.RoomEventModel, error) {
	return nil, f.err
}
