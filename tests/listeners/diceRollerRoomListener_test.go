package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/listeners"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
	"github.com/stretchr/testify/require"
)

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
			Type:    string(roomModel.EventDiceRoll),
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
		require.Equal(t, string(roomModel.EventDiceRoll), event.Type)
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
