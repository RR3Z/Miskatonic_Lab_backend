package room

type RoomCreateSucceeded struct {
	RoomID  string
	OwnerID string
}

type RoomCreateFailed struct {
	OwnerID string
	Err     error
}

func (RoomCreateSucceeded) EventName() string { return "room.create_succeeded" }
func (RoomCreateFailed) EventName() string    { return "room.create_failed" }

type RoomGetSucceeded struct {
	RoomID string
	UserID string
}

type RoomGetFailed struct {
	RoomID string
	UserID string
	Err    error
}

func (RoomGetSucceeded) EventName() string { return "room.get_succeeded" }
func (RoomGetFailed) EventName() string    { return "room.get_failed" }

type RoomUpdateSucceeded struct {
	RoomID  string
	OwnerID string
}

type RoomUpdateFailed struct {
	RoomID  string
	OwnerID string
	Err     error
}

func (RoomUpdateSucceeded) EventName() string { return "room.update_succeeded" }
func (RoomUpdateFailed) EventName() string    { return "room.update_failed" }

type RoomTransferOwnershipSucceeded struct {
	RoomID     string
	OwnerID    string
	NewOwnerID string
}

type RoomTransferOwnershipFailed struct {
	RoomID     string
	OwnerID    string
	NewOwnerID string
	Err        error
}

func (RoomTransferOwnershipSucceeded) EventName() string {
	return "room.transfer_ownership_succeeded"
}
func (RoomTransferOwnershipFailed) EventName() string {
	return "room.transfer_ownership_failed"
}

type RoomDeleteSucceeded struct {
	RoomID  string
	OwnerID string
}

type RoomDeleteFailed struct {
	RoomID  string
	OwnerID string
	Err     error
}

func (RoomDeleteSucceeded) EventName() string { return "room.delete_succeeded" }
func (RoomDeleteFailed) EventName() string    { return "room.delete_failed" }
