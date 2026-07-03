package ws

import (
	"context"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
)

type RoomHub struct {
	rooms      map[string]map[*Client]struct{}
	Register   chan *Client
	Unregister chan *Client
	broadcast  chan roomEvents.Event
}

func NewRoomHub() *RoomHub {
	return &RoomHub{
		rooms:      make(map[string]map[*Client]struct{}),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan roomEvents.Event),
	}
}

func (rh *RoomHub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-rh.Register:
			if rh.rooms[client.roomID] == nil {
				rh.rooms[client.roomID] = make(map[*Client]struct{})
			}
			rh.rooms[client.roomID][client] = struct{}{}

		case client := <-rh.Unregister:
			rh.removeClient(client)

		case event := <-rh.broadcast:
			rh.broadcastToRoom(event)

		}
	}
}

func (rh *RoomHub) removeClient(client *Client) {
	if clients := rh.rooms[client.roomID]; clients != nil {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.send)
		}

		if len(clients) == 0 {
			delete(rh.rooms, client.roomID)
		}
	}
}
