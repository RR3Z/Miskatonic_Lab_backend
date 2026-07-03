package room

import roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"

type targetedEvent struct {
	roomID  string
	userIDs []string
	event   roomEvents.Event
}

func (rh *RoomHub) Broadcast(event roomEvents.Event) {
	rh.broadcast <- event
}

func (rh *RoomHub) SendToUsers(roomID string, userIDs []string, event roomEvents.Event) {
	if len(userIDs) == 0 {
		return
	}

	rh.targeted <- targetedEvent{roomID: roomID, userIDs: userIDs, event: event}
}

func (rh *RoomHub) broadcastToRoom(event roomEvents.Event) {
	for client := range rh.rooms[event.RoomID] {
		rh.sendToClient(client, event)
	}
}

func (rh *RoomHub) sendTargeted(target targetedEvent) {
	recipients := target.recipientSet()
	for client := range rh.rooms[target.roomID] {
		if _, ok := recipients[client.userID]; !ok {
			continue
		}

		rh.sendToClient(client, target.event)
	}
}

func (target targetedEvent) recipientSet() map[string]struct{} {
	recipients := make(map[string]struct{}, len(target.userIDs))
	for _, userID := range target.userIDs {
		recipients[userID] = struct{}{}
	}
	return recipients
}

func (rh *RoomHub) sendToClient(client *Client, event roomEvents.Event) {
	select {
	case client.send <- event:
	default:
		rh.removeClient(client)
	}
}
