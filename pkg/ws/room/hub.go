package room

import (
	"context"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
)

type RoomHub struct {
	rooms      map[string]map[*Client]struct{}
	Register   chan *Client
	Unregister chan *Client
	broadcast  chan roomEvents.Event
	targeted   chan targetedEvent
	closeRoom  chan closeRoomCommand
	closeUser  chan closeUserCommand
}

func NewRoomHub() *RoomHub {
	return &RoomHub{
		rooms:      make(map[string]map[*Client]struct{}),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan roomEvents.Event),
		targeted:   make(chan targetedEvent),
		closeRoom:  make(chan closeRoomCommand, 64),
		closeUser:  make(chan closeUserCommand, 64),
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

		case target := <-rh.targeted:
			rh.sendTargeted(target)

		case command := <-rh.closeRoom:
			rh.closeRoomClients(command)

		case command := <-rh.closeUser:
			rh.closeUserClients(command)
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
