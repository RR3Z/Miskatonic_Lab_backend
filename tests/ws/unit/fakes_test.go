package ws_test

import (
	"context"
	"encoding/json"
	"testing"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
)

type fakeRoomEventService struct {
	inputs chan roomModel.CreateChatMessageInput
}

func newFakeRoomEventService() *fakeRoomEventService {
	return &fakeRoomEventService{
		inputs: make(chan roomModel.CreateChatMessageInput, 8),
	}
}

func (f *fakeRoomEventService) CreateChatMessage(_ context.Context, input roomModel.CreateChatMessageInput) (roomModel.RoomEventModel, error) {
	f.inputs <- input

	payload, err := json.Marshal(roomEvents.ChatMessagePayload{Text: input.Text})
	if err != nil {
		return roomModel.RoomEventModel{}, err
	}

	return roomModel.RoomEventModel{
		RoomID:  input.RoomID,
		ActorID: input.ActorID,
		Type:    string(roomEvents.EventChatMessage),
		Payload: payload,
	}, nil
}

func (f *fakeRoomEventService) waitForCreateChatInput(t *testing.T, ctx context.Context) roomModel.CreateChatMessageInput {
	t.Helper()

	select {
	case input := <-f.inputs:
		return input
	case <-ctx.Done():
		t.Fatal("timed out waiting for chat message persistence")
		return roomModel.CreateChatMessageInput{}
	}
}
