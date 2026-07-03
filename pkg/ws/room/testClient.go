package room

import roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"

func NewTestClient(hub *RoomHub, roomID string) (*Client, <-chan roomEvents.Event) {
	return NewTestClientWithUser(hub, roomID, "")
}

func NewTestClientWithUser(hub *RoomHub, roomID string, userID string) (*Client, <-chan roomEvents.Event) {
	ch := make(chan roomEvents.Event, 256)
	return &Client{
		roomID: roomID,
		userID: userID,
		send:   ch,
		hub:    hub,
	}, ch
}
