package ws

import (
	"context"
	"encoding/json"
	"errors"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	wsHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/helpers"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomEventService interface {
	CreateChatMessage(ctx context.Context, input roomModel.CreateChatMessageInput) (roomModel.RoomEventModel, error)
}

type commandEnvelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type commandContext struct {
	roomID  pgtype.UUID
	actorID string
}

type DispatchResult struct {
	Broadcast *roomEvents.Event
	Reply     *roomEvents.Event
}

type commandHandler func(ctx context.Context, command commandEnvelope, backend commandContext) (DispatchResult, error)

type CommandDispatcher struct {
	handlers map[string]commandHandler
}

func NewCommandDispatcher(service RoomEventService) *CommandDispatcher {
	dispatcher := &CommandDispatcher{
		handlers: make(map[string]commandHandler),
	}

	dispatcher.handlers[string(roomEvents.EventChatMessage)] = chatMessageHandler(service)

	return dispatcher
}

func (d *CommandDispatcher) Dispatch(ctx context.Context, command commandEnvelope, backend commandContext) (DispatchResult, error) {
	handler, ok := d.handlers[command.Type]
	if !ok {
		reply := wsHelpers.UnsupportedCommandTypeEvent(backend.roomID.String(), backend.actorID)
		return DispatchResult{Reply: &reply}, nil
	}

	return handler(ctx, command, backend)
}

func chatMessageHandler(service RoomEventService) commandHandler {
	return func(ctx context.Context, command commandEnvelope, backend commandContext) (DispatchResult, error) {
		var payload roomEvents.ChatMessagePayload
		if err := json.Unmarshal(command.Payload, &payload); err != nil {
			reply := wsHelpers.InvalidCommandPayloadEvent(backend.roomID.String(), backend.actorID)
			return DispatchResult{Reply: &reply}, nil
		}

		event, err := service.CreateChatMessage(ctx, roomModel.CreateChatMessageInput{
			RoomID:  backend.roomID,
			ActorID: backend.actorID,
			Text:    payload.Text,
		})
		if err != nil {
			if errors.Is(err, roomService.ErrInvalidInput) {
				reply := wsHelpers.InvalidCommandPayloadEvent(backend.roomID.String(), backend.actorID)
				return DispatchResult{Reply: &reply}, nil
			}
			return DispatchResult{}, err
		}

		broadcast := wsHelpers.EventFromRoomEventModel(event)
		return DispatchResult{Broadcast: &broadcast}, nil
	}
}
