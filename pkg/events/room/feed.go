package room

type RoomEventsListSucceeded struct {
	RoomID string
	UserID string
	Count  int
}

type RoomEventsListFailed struct {
	RoomID string
	UserID string
	Err    error
}

func (RoomEventsListSucceeded) EventName() string { return "room.events.list_succeeded" }
func (RoomEventsListFailed) EventName() string    { return "room.events.list_failed" }

type RoomChatMessageCreateSucceeded struct {
	RoomID  string
	ActorID string
	EventID string
}

type RoomChatMessageCreateFailed struct {
	RoomID  string
	ActorID string
	Err     error
}

func (RoomChatMessageCreateSucceeded) EventName() string {
	return "room.event.chat_create_succeeded"
}
func (RoomChatMessageCreateFailed) EventName() string {
	return "room.event.chat_create_failed"
}

type RoomDiceRollEventCreateSucceeded struct {
	RoomID      string
	ActorID     string
	EventID     string
	RollID      string
	CharacterID string
}

type RoomDiceRollEventCreateFailed struct {
	RoomID      string
	ActorID     string
	RollID      string
	CharacterID string
	Err         error
}

func (RoomDiceRollEventCreateSucceeded) EventName() string {
	return "room.event.dice_roll_create_succeeded"
}
func (RoomDiceRollEventCreateFailed) EventName() string {
	return "room.event.dice_roll_create_failed"
}

type RoomCharacterChangedEventsCreateSucceeded struct {
	ActorID     string
	CharacterID string
	Count       int
}

type RoomCharacterChangedEventsCreateFailed struct {
	ActorID     string
	CharacterID string
	Err         error
}

func (RoomCharacterChangedEventsCreateSucceeded) EventName() string {
	return "room.event.character_changed_create_succeeded"
}
func (RoomCharacterChangedEventsCreateFailed) EventName() string {
	return "room.event.character_changed_create_failed"
}
