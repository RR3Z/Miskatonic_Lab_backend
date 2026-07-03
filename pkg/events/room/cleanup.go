package room

type RoomCleanupSucceeded struct {
	InactiveDeleted int
	InvalidDeleted  int
	DeletedCount    int
}

type RoomCleanupFailed struct {
	Err error
}

func (RoomCleanupSucceeded) EventName() string { return "room.cleanup_succeeded" }
func (RoomCleanupFailed) EventName() string    { return "room.cleanup_failed" }
