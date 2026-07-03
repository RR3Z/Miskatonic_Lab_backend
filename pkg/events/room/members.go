package room

type RoomMemberJoinSucceeded struct {
	RoomID   string
	UserID   string
	MemberID string
}

type RoomMemberJoinFailed struct {
	RoomID string
	UserID string
	Err    error
}

func (RoomMemberJoinSucceeded) EventName() string { return "room.member.join_succeeded" }
func (RoomMemberJoinFailed) EventName() string    { return "room.member.join_failed" }

type RoomMemberLeaveSucceeded struct {
	RoomID        string
	UserID        string
	DeletedRoomID *string
}

type RoomMemberLeaveFailed struct {
	RoomID string
	UserID string
	Err    error
}

func (RoomMemberLeaveSucceeded) EventName() string { return "room.member.leave_succeeded" }
func (RoomMemberLeaveFailed) EventName() string    { return "room.member.leave_failed" }

type RoomMemberKickSucceeded struct {
	RoomID       string
	ActorUserID  string
	TargetUserID string
}

type RoomMemberKickFailed struct {
	RoomID       string
	ActorUserID  string
	TargetUserID string
	Err          error
}

func (RoomMemberKickSucceeded) EventName() string { return "room.member.kick_succeeded" }
func (RoomMemberKickFailed) EventName() string    { return "room.member.kick_failed" }

type RoomMemberSelectCharacterSucceeded struct {
	RoomID      string
	UserID      string
	CharacterID string
}

type RoomMemberSelectCharacterFailed struct {
	RoomID      string
	UserID      string
	CharacterID string
	Err         error
}

func (RoomMemberSelectCharacterSucceeded) EventName() string {
	return "room.member.select_character_succeeded"
}
func (RoomMemberSelectCharacterFailed) EventName() string {
	return "room.member.select_character_failed"
}

type RoomMemberChangeRoleSucceeded struct {
	RoomID       string
	ActorUserID  string
	TargetUserID string
	Role         string
}

type RoomMemberChangeRoleFailed struct {
	RoomID       string
	ActorUserID  string
	TargetUserID string
	Role         string
	Err          error
}

func (RoomMemberChangeRoleSucceeded) EventName() string {
	return "room.member.change_role_succeeded"
}
func (RoomMemberChangeRoleFailed) EventName() string {
	return "room.member.change_role_failed"
}
