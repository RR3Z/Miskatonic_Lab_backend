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
