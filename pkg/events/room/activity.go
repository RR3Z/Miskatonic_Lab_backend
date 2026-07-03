package room

type RoomActivityTouchSucceeded struct {
	RoomID string
	UserID string
}

type RoomActivityTouchFailed struct {
	RoomID string
	UserID string
	Err    error
}

func (RoomActivityTouchSucceeded) EventName() string { return "room.activity.touch_succeeded" }
func (RoomActivityTouchFailed) EventName() string    { return "room.activity.touch_failed" }
