package ws

import "github.com/coder/websocket"

type closeRoomCommand struct {
	roomID string
	reason string
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

func (rh *RoomHub) closeRoomClients(command closeRoomCommand) {
	for client := range rh.rooms[command.roomID] {
		rh.closeClient(client, command.reason)
	}
	delete(rh.rooms, command.roomID)
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
