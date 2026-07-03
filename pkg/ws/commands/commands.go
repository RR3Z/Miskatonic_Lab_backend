package commands

import (
	"context"
	"encoding/json"
	"errors"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	wsHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/helpers"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomEventService interface {
	CreateChatMessage(ctx context.Context, input roomModel.CreateChatMessageInput) (roomModel.RoomEventModel, error)
}

type Envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Context struct {
	RoomID  pgtype.UUID
	ActorID string
}

type DispatchResult struct {
	Broadcast *roomModel.Event
	Reply     *roomModel.Event
}

type commandHandler func(ctx context.Context, command Envelope, backend Context) (DispatchResult, error)

type CommandDispatcher struct {
	handlers map[string]commandHandler
}

func NewCommandDispatcher(service RoomEventService) *CommandDispatcher {
	dispatcher := &CommandDispatcher{
		handlers: make(map[string]commandHandler),
	}

	dispatcher.handlers[string(roomModel.EventChatMessage)] = chatMessageHandler(service)

	return dispatcher
}

func (d *CommandDispatcher) Dispatch(ctx context.Context, command Envelope, backend Context) (DispatchResult, error) {
	handler, ok := d.handlers[command.Type]
	if !ok {
		reply := wsHelpers.UnsupportedCommandTypeEvent(backend.RoomID.String(), backend.ActorID)
		return DispatchResult{Reply: &reply}, nil
	}

	return handler(ctx, command, backend)
}

func chatMessageHandler(service RoomEventService) commandHandler {
	return func(ctx context.Context, command Envelope, backend Context) (DispatchResult, error) {
		var payload roomModel.ChatMessagePayload
		if err := json.Unmarshal(command.Payload, &payload); err != nil {
			reply := wsHelpers.InvalidCommandPayloadEvent(backend.RoomID.String(), backend.ActorID)
			return DispatchResult{Reply: &reply}, nil
		}

		event, err := service.CreateChatMessage(ctx, roomModel.CreateChatMessageInput{
			RoomID:  backend.RoomID,
			ActorID: backend.ActorID,
			Text:    payload.Text,
		})
		if err != nil {
			if errors.Is(err, roomService.ErrInvalidInput) {
				reply := wsHelpers.InvalidCommandPayloadEvent(backend.RoomID.String(), backend.ActorID)
				return DispatchResult{Reply: &reply}, nil
			}
			return DispatchResult{}, err
		}

		broadcast := wsHelpers.EventFromRoomEventModel(event)
		return DispatchResult{Broadcast: &broadcast}, nil
	}
}
