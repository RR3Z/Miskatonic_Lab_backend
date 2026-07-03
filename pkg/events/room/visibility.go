package room

type RoomSelectedCharactersListSucceeded struct {
	RoomID string
	UserID string
	Count  int
}

type RoomSelectedCharactersListFailed struct {
	RoomID string
	UserID string
	Err    error
}

func (RoomSelectedCharactersListSucceeded) EventName() string {
	return "room.selected_characters.list_succeeded"
}
func (RoomSelectedCharactersListFailed) EventName() string {
	return "room.selected_characters.list_failed"
}
