package room

type RoomEnsureMemberSucceeded struct {
	RoomID string
	UserID string
}

type RoomEnsureMemberFailed struct {
	RoomID string
	UserID string
	Err    error
}

func (RoomEnsureMemberSucceeded) EventName() string { return "room.permission.ensure_member_succeeded" }
func (RoomEnsureMemberFailed) EventName() string    { return "room.permission.ensure_member_failed" }

type RoomEnsureOwnerSucceeded struct {
	RoomID string
	UserID string
}

type RoomEnsureOwnerFailed struct {
	RoomID string
	UserID string
	Err    error
}

func (RoomEnsureOwnerSucceeded) EventName() string { return "room.permission.ensure_owner_succeeded" }
func (RoomEnsureOwnerFailed) EventName() string    { return "room.permission.ensure_owner_failed" }

type RoomEnsureCanPublishEventSucceeded struct {
	RoomID string
	UserID string
}

type RoomEnsureCanPublishEventFailed struct {
	RoomID string
	UserID string
	Err    error
}

func (RoomEnsureCanPublishEventSucceeded) EventName() string {
	return "room.event.ensure_publish_succeeded"
}
func (RoomEnsureCanPublishEventFailed) EventName() string {
	return "room.event.ensure_publish_failed"
}
