package room

import "github.com/coder/websocket"

type closeRoomCommand struct {
	roomID string
	reason string
}

type closeUserCommand struct {
	roomID string
	userID string
	reason string
	done   chan struct{}
}

func (rh *RoomHub) CloseRoom(roomID string, reason string) {
	if roomID == "" {
		return
	}

	command := closeRoomCommand{roomID: roomID, reason: reason}
	select {
	case rh.closeRoom <- command:
	default:
		go func() {
			rh.closeRoom <- command
		}()
	}
}

func (rh *RoomHub) CloseUser(roomID string, userID string, reason string) {
	if roomID == "" || userID == "" {
		return
	}

	command := closeUserCommand{
		roomID: roomID,
		userID: userID,
		reason: reason,
		done:   make(chan struct{}),
	}
	select {
	case rh.closeUser <- command:
	default:
		go func() {
			rh.closeUser <- command
		}()
	}
	<-command.done
}

func (rh *RoomHub) closeRoomClients(command closeRoomCommand) {
	for client := range rh.rooms[command.roomID] {
		rh.closeClient(client, command.reason)
	}
	delete(rh.rooms, command.roomID)
}

func (rh *RoomHub) closeUserClients(command closeUserCommand) {
	for client := range rh.rooms[command.roomID] {
		if client.userID == command.userID {
			rh.closeClient(client, command.reason)
		}
	}
	close(command.done)
}

func (rh *RoomHub) closeClient(client *Client, reason string) {
	if clients := rh.rooms[client.roomID]; clients != nil {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.send)
		}

		if len(clients) == 0 {
			delete(rh.rooms, client.roomID)
		}
	}

	if client.conn != nil {
		go client.conn.Close(websocket.StatusNormalClosure, reason)
	}
}
