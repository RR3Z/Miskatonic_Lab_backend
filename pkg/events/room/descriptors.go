package room

import "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"

const domain = "room"

var descriptors = []events.EventDescriptor{
	descriptor(RoomCreateSucceeded{}, "room", "create", events.OutcomeSucceeded),
	descriptor(RoomCreateFailed{}, "room", "create", events.OutcomeFailed),
	descriptor(RoomGetSucceeded{}, "room", "get", events.OutcomeSucceeded),
	descriptor(RoomGetFailed{}, "room", "get", events.OutcomeFailed),
	descriptor(RoomUpdateSucceeded{}, "room", "update", events.OutcomeSucceeded),
	descriptor(RoomUpdateFailed{}, "room", "update", events.OutcomeFailed),
	descriptor(RoomTransferOwnershipSucceeded{}, "room", "transfer_ownership", events.OutcomeSucceeded),
	descriptor(RoomTransferOwnershipFailed{}, "room", "transfer_ownership", events.OutcomeFailed),
	descriptor(RoomDeleteSucceeded{}, "room", "delete", events.OutcomeSucceeded),
	descriptor(RoomDeleteFailed{}, "room", "delete", events.OutcomeFailed),

	descriptor(RoomMemberJoinSucceeded{}, "room_member", "join", events.OutcomeSucceeded),
	descriptor(RoomMemberJoinFailed{}, "room_member", "join", events.OutcomeFailed),
	descriptor(RoomMemberLeaveSucceeded{}, "room_member", "leave", events.OutcomeSucceeded),
	descriptor(RoomMemberLeaveFailed{}, "room_member", "leave", events.OutcomeFailed),
	descriptor(RoomMemberKickSucceeded{}, "room_member", "kick", events.OutcomeSucceeded),
	descriptor(RoomMemberKickFailed{}, "room_member", "kick", events.OutcomeFailed),
	descriptor(RoomMemberSelectCharacterSucceeded{}, "room_member", "select_character", events.OutcomeSucceeded),
	descriptor(RoomMemberSelectCharacterFailed{}, "room_member", "select_character", events.OutcomeFailed),
	descriptor(RoomMemberChangeRoleSucceeded{}, "room_member", "change_role", events.OutcomeSucceeded),
	descriptor(RoomMemberChangeRoleFailed{}, "room_member", "change_role", events.OutcomeFailed),

	descriptor(RoomSelectedCharactersListSucceeded{}, "selected_characters", "list", events.OutcomeSucceeded),
	descriptor(RoomSelectedCharactersListFailed{}, "selected_characters", "list", events.OutcomeFailed),
	descriptor(RoomActivityTouchSucceeded{}, "room_activity", "touch", events.OutcomeSucceeded),
	descriptor(RoomActivityTouchFailed{}, "room_activity", "touch", events.OutcomeFailed),
	descriptor(RoomEnsureMemberSucceeded{}, "room_permission", "ensure_member", events.OutcomeSucceeded),
	descriptor(RoomEnsureMemberFailed{}, "room_permission", "ensure_member", events.OutcomeFailed),
	descriptor(RoomEnsureOwnerSucceeded{}, "room_permission", "ensure_owner", events.OutcomeSucceeded),
	descriptor(RoomEnsureOwnerFailed{}, "room_permission", "ensure_owner", events.OutcomeFailed),
	descriptor(RoomEnsureCanPublishEventSucceeded{}, "room_event_permission", "ensure_publish", events.OutcomeSucceeded),
	descriptor(RoomEnsureCanPublishEventFailed{}, "room_event_permission", "ensure_publish", events.OutcomeFailed),

	descriptor(RoomEventsListSucceeded{}, "room_events", "list", events.OutcomeSucceeded),
	descriptor(RoomEventsListFailed{}, "room_events", "list", events.OutcomeFailed),
	descriptor(RoomChatMessageCreateSucceeded{}, "room_event", "chat_create", events.OutcomeSucceeded),
	descriptor(RoomChatMessageCreateFailed{}, "room_event", "chat_create", events.OutcomeFailed),
	descriptor(RoomDiceRollEventCreateSucceeded{}, "room_event", "dice_roll_create", events.OutcomeSucceeded),
	descriptor(RoomDiceRollEventCreateFailed{}, "room_event", "dice_roll_create", events.OutcomeFailed),
	descriptor(RoomCharacterChangedEventsCreateSucceeded{}, "room_event", "character_changed_create", events.OutcomeSucceeded),
	descriptor(RoomCharacterChangedEventsCreateFailed{}, "room_event", "character_changed_create", events.OutcomeFailed),

	descriptor(RoomCleanupSucceeded{}, "room_cleanup", "cleanup", events.OutcomeSucceeded),
	descriptor(RoomCleanupFailed{}, "room_cleanup", "cleanup", events.OutcomeFailed),
}

func Descriptors() []events.EventDescriptor {
	return append([]events.EventDescriptor(nil), descriptors...)
}

func AllEvents() []events.Event {
	return events.EventPrototypes(descriptors)
}

func descriptor(event events.Event, resource string, action string, outcome events.EventOutcome) events.EventDescriptor {
	return events.EventDescriptor{
		Event:    event,
		Domain:   domain,
		Resource: resource,
		Action:   action,
		Outcome:  outcome,
	}
}
