package ws

import (
	"context"
	"encoding/json"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomEventService interface {
	CreateChatMessage(ctx context.Context, input roomModel.CreateChatMessageInput) (roomModel.RoomEventModel, error)
}

type Client struct {
	roomID   string
	roomUUID pgtype.UUID
	userID   string
	conn     *websocket.Conn
	send     chan roomEvents.Event
	hub      *RoomHub
	service  RoomEventService
}

func NewClient(hub *RoomHub, service RoomEventService, roomID pgtype.UUID, userID string, conn *websocket.Conn) *Client {
	return &Client{
		roomID:   roomID.String(),
		roomUUID: roomID,
		userID:   userID,
		conn:     conn,
		send:     make(chan roomEvents.Event, 32),
		hub:      hub,
		service:  service,
	}
}

func (c *Client) ReadLoop(ctx context.Context) {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close(websocket.StatusNormalClosure, "client disconnected")
	}()

	for {
		var command incomingEvent

		if err := wsjson.Read(ctx, c.conn, &command); err != nil {
			return
		}

		event, err := c.handleCommand(ctx, command)
		if err != nil {
			return
		}
		c.hub.Broadcast(event)
	}
}

func (c *Client) WriteLoop(ctx context.Context) {
	defer c.conn.Close(websocket.StatusNormalClosure, "write loop stopped")

	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-c.send:
			if !ok {
				return
			}

			if err := wsjson.Write(ctx, c.conn, event); err != nil {
				return
			}
		}
	}
}

type incomingEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (c *Client) handleCommand(ctx context.Context, command incomingEvent) (roomEvents.Event, error) {
	switch command.Type {
	case string(roomEvents.EventChatMessage):
		var payload roomEvents.ChatMessagePayload
		if err := json.Unmarshal(command.Payload, &payload); err != nil {
			return roomEvents.Event{}, err
		}

		event, err := c.service.CreateChatMessage(ctx, roomModel.CreateChatMessageInput{
			RoomID:  c.roomUUID,
			ActorID: c.userID,
			Text:    payload.Text,
		})
		if err != nil {
			return roomEvents.Event{}, err
		}

		return eventFromModel(event), nil
	default:
		return roomEvents.Event{}, websocket.CloseError{
			Code:   websocket.StatusUnsupportedData,
			Reason: "unsupported room event type",
		}
	}
}

func eventFromModel(event roomModel.RoomEventModel) roomEvents.Event {
	return roomEvents.Event{
		Type:    event.Type,
		RoomID:  event.RoomID.String(),
		ActorID: event.ActorID,
		Payload: event.Payload,
	}
}
