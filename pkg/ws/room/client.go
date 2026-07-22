package room

import (
	"context"
	"time"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	wsCommands "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/commands"
	wsHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/helpers"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	clientPingInterval = 10 * time.Second
	clientPingTimeout  = 5 * time.Second
)

type Client struct {
	roomID     string
	roomUUID   pgtype.UUID
	userID     string
	conn       *websocket.Conn
	send       chan roomEvents.Event
	hub        *RoomHub
	dispatcher *wsCommands.CommandDispatcher
}

func NewClient(hub *RoomHub, dispatcher *wsCommands.CommandDispatcher, roomID pgtype.UUID, userID string, conn *websocket.Conn) *Client {
	return &Client{
		roomID:     roomID.String(),
		roomUUID:   roomID,
		userID:     userID,
		conn:       conn,
		send:       make(chan roomEvents.Event, 32),
		hub:        hub,
		dispatcher: dispatcher,
	}
}

func (c *Client) ReadLoop(ctx context.Context) {
	closeCode := websocket.StatusNormalClosure
	closeReason := "client disconnected"

	defer func() {
		c.hub.Unregister <- c
		c.conn.Close(closeCode, closeReason)
	}()

	for {
		var command wsCommands.Envelope

		if err := wsjson.Read(ctx, c.conn, &command); err != nil {
			return
		}

		result, err := c.dispatcher.Dispatch(ctx, command, wsCommands.Context{
			RoomID:  c.roomUUID,
			ActorID: c.userID,
		})
		if err != nil {
			closeCode, closeReason = wsHelpers.CloseStatusForCommandError(err)
			return
		}

		if result.Reply != nil {
			if ok := c.sendDirect(ctx, *result.Reply); !ok {
				return
			}
			continue
		}

		if result.Broadcast != nil {
			c.hub.Broadcast(*result.Broadcast)
		}
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

func (c *Client) PingLoop(ctx context.Context) {
	if c.conn == nil {
		return
	}

	ticker := time.NewTicker(clientPingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, clientPingTimeout)
			err := c.conn.Ping(pingCtx)
			cancel()
			if err != nil {
				c.conn.Close(websocket.StatusGoingAway, "websocket ping failed")
				return
			}
		}
	}
}

func (c *Client) sendDirect(ctx context.Context, event roomEvents.Event) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	select {
	case <-ctx.Done():
		return false
	case c.send <- event:
		return true
	}
}
