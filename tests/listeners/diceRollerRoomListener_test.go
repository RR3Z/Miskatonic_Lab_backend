package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/listeners"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

type fakeListenerRoomService struct {
	diceEvent  roomModel.RoomEventModel
	diceErr    error
	diceCalls  int
	diceInput  roomModel.CreateDiceRollRoomEventInput
	memberErr  error
	publishErr error
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

func listenerTestUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}

func TestDiceRollerRoomListener_NoRoomID_NoOp(t *testing.T) {
	svc := &fakeListenerRoomService{}
	listener := listeners.NewDiceRollerRoomListener(svc, ws.NewRoomHub())

	listener.Handle(context.Background(), diceEvents.DiceRollMakeSucceeded{
		UserID:      "user_1",
		CharacterID: "11111111-1111-1111-1111-111111111111",
		RollID:      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Expression:  "1d20",
		Result:      15,
		RoomID:      nil,
	})

	require.Zero(t, svc.diceCalls)
}

func TestDiceRollerRoomListener_Success_CreatesRoomEventAndBroadcasts(t *testing.T) {
	roomUUID := listenerTestUUID("22222222-2222-2222-2222-222222222222")
	svc := &fakeListenerRoomService{
		diceEvent: roomModel.RoomEventModel{
			RoomID:  roomUUID,
			ActorID: "user_1",
			Type:    string(roomEvents.EventDiceRoll),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	hub := ws.NewRoomHub()
	go hub.Run(ctx)
	client, events := ws.NewTestClient(hub, roomUUID.String())
	hub.Register <- client
	_ = client

	listener := listeners.NewDiceRollerRoomListener(svc, hub)

	roomIDStr := "22222222-2222-2222-2222-222222222222"
	listener.Handle(ctx, diceEvents.DiceRollMakeSucceeded{
		UserID:      "user_1",
		CharacterID: "11111111-1111-1111-1111-111111111111",
		RollID:      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Expression:  "2d6+1",
		Result:      8,
		Details:     []byte(`[{"type":"dice","sides":6,"rolls":[5,3]},{"type":"modifier","value":1}]`),
		RoomID:      &roomIDStr,
	})

	require.Equal(t, 1, svc.diceCalls)
	require.Equal(t, "user_1", svc.diceInput.ActorID)
	require.Equal(t, "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", svc.diceInput.RollID)
	require.Equal(t, "2d6+1", svc.diceInput.Expression)
	require.Equal(t, int32(8), svc.diceInput.Result)

	select {
	case event := <-events:
		require.Equal(t, string(roomEvents.EventDiceRoll), event.Type)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for broadcast")
	}
}

func TestDiceRollerRoomListener_MembershipError_LogsAndNoBroadcast(t *testing.T) {
	svc := &fakeListenerRoomService{
		diceErr: errors.New("not a member"),
	}
	listener := listeners.NewDiceRollerRoomListener(svc, ws.NewRoomHub())

	roomIDStr := "22222222-2222-2222-2222-222222222222"
	listener.Handle(context.Background(), diceEvents.DiceRollMakeSucceeded{
		UserID:      "user_1",
		CharacterID: "11111111-1111-1111-1111-111111111111",
		RollID:      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Expression:  "1d20",
		Result:      10,
		RoomID:      &roomIDStr,
	})

	require.Equal(t, 1, svc.diceCalls)
}

func TestDiceRollerRoomListener_InvalidRoomID_LogsAndNoOp(t *testing.T) {
	svc := &fakeListenerRoomService{}
	listener := listeners.NewDiceRollerRoomListener(svc, ws.NewRoomHub())

	badRoomID := "not-a-uuid"
	listener.Handle(context.Background(), diceEvents.DiceRollMakeSucceeded{
		UserID:      "user_1",
		CharacterID: "11111111-1111-1111-1111-111111111111",
		RollID:      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Expression:  "1d20",
		Result:      10,
		RoomID:      &badRoomID,
	})

	require.Zero(t, svc.diceCalls)
}

func TestDiceRollerRoomListener_WrongEventType_NoOp(t *testing.T) {
	svc := &fakeListenerRoomService{}
	listener := listeners.NewDiceRollerRoomListener(svc, ws.NewRoomHub())

	listener.Handle(context.Background(), diceEvents.DiceRollMakeFailed{
		UserID:      "user_1",
		CharacterID: "11111111-1111-1111-1111-111111111111",
		Err:         errors.New("fail"),
		RoomID:      nil,
	})

	require.Zero(t, svc.diceCalls)
}
