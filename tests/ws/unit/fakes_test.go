package ws_test

import (
	"context"
	"encoding/json"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"testing"
)

type fakeRoomEventService struct {
	inputs chan roomModel.CreateChatMessageInput
	err    error
}

func newFakeRoomEventService() *fakeRoomEventService {
	return &fakeRoomEventService{
		inputs: make(chan roomModel.CreateChatMessageInput, 8),
	}
}

func (f *fakeRoomEventService) CreateChatMessage(_ context.Context, input roomModel.CreateChatMessageInput) (roomModel.RoomEventModel, error) {
	f.inputs <- input
	if f.err != nil {
		return roomModel.RoomEventModel{}, f.err
	}

	payload, err := json.Marshal(roomModel.ChatMessagePayload{Text: input.Text})
	if err != nil {
		return roomModel.RoomEventModel{}, err
	}

	return roomModel.RoomEventModel{
		RoomID:  input.RoomID,
		ActorID: input.ActorID,
		Type:    string(roomModel.EventChatMessage),
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
