package tests

import (
	"context"
	"errors"
	"testing"

	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/stretchr/testify/require"
)

var errRoomEventTest = errors.New("room service failed")

func TestEventPublishingRoomService_CreateRoom_Success(t *testing.T) {
	_, publisher, service := newRoomEventPublishingTestSubject()

	room, err := service.CreateRoom(context.Background(), model.CreateRoomInput{OwnerID: roomEventTestOwnerID})
	require.NoError(t, err)
	require.Equal(t, roomEventTestRoomID, room.ID.String())
	requireRoomPublishedEvent(t, publisher, roomEvents.RoomCreateSucceeded{
		RoomID:  roomEventTestRoomID,
		OwnerID: roomEventTestOwnerID,
	})
}

func TestEventPublishingRoomService_UpdateRoom_Failure(t *testing.T) {
	next, publisher, service := newRoomEventPublishingTestSubject()
	next.err = errRoomEventTest

	_, err := service.UpdateRoom(context.Background(), model.UpdateRoomInput{
		RoomID:  roomTestUUID(roomEventTestRoomID),
		OwnerID: roomEventTestOwnerID,
	})
	require.ErrorIs(t, err, errRoomEventTest)
	requireRoomPublishedEvent(t, publisher, roomEvents.RoomUpdateFailed{
		RoomID:  roomEventTestRoomID,
		OwnerID: roomEventTestOwnerID,
		Err:     errRoomEventTest,
	})
}

func TestEventPublishingRoomService_JoinRoom_Success(t *testing.T) {
	_, publisher, service := newRoomEventPublishingTestSubject()

	member, err := service.JoinRoom(context.Background(), model.JoinRoomInput{
		RoomID: roomTestUUID(roomEventTestRoomID),
		UserID: roomEventTestUserID,
	})
	require.NoError(t, err)
	require.Equal(t, roomEventTestMemberID, member.ID.String())
	requireRoomPublishedEvent(t, publisher, roomEvents.RoomMemberJoinSucceeded{
		RoomID:   roomEventTestRoomID,
		UserID:   roomEventTestUserID,
		MemberID: roomEventTestMemberID,
	})
}

func TestEventPublishingRoomService_LeaveRoom_WithDeletedRoomID(t *testing.T) {
	next, publisher, service := newRoomEventPublishingTestSubject()
	deletedRoomID := roomTestUUID(roomEventTestRoomID)
	next.leaveResult = model.LeaveRoomResult{DeletedRoomID: &deletedRoomID}

	result, err := service.LeaveRoom(context.Background(), model.LeaveRoomInput{
		RoomID: roomTestUUID(roomEventTestRoomID),
		UserID: roomEventTestUserID,
	})
	require.NoError(t, err)
	require.NotNil(t, result.DeletedRoomID)

	expectedDeletedRoomID := roomEventTestRoomID
	requireRoomPublishedEvent(t, publisher, roomEvents.RoomMemberLeaveSucceeded{
		RoomID:        roomEventTestRoomID,
		UserID:        roomEventTestUserID,
		DeletedRoomID: &expectedDeletedRoomID,
	})
}

func TestEventPublishingRoomService_EnsureCanPublishRoomEvent_Failure(t *testing.T) {
	next, publisher, service := newRoomEventPublishingTestSubject()
	next.err = errRoomEventTest

	err := service.EnsureCanPublishRoomEvent(context.Background(), roomTestUUID(roomEventTestRoomID), roomEventTestUserID)
	require.ErrorIs(t, err, errRoomEventTest)
	requireRoomPublishedEvent(t, publisher, roomEvents.RoomEnsureCanPublishEventFailed{
		RoomID: roomEventTestRoomID,
		UserID: roomEventTestUserID,
		Err:    errRoomEventTest,
	})
}

func TestEventPublishingRoomService_ListRoomEvents_Success(t *testing.T) {
	_, publisher, service := newRoomEventPublishingTestSubject()

	events, err := service.ListRoomEvents(context.Background(), model.ListRoomEventsInput{
		RoomID: roomTestUUID(roomEventTestRoomID),
		UserID: roomEventTestUserID,
	})
	require.NoError(t, err)
	require.Len(t, events, 2)
	requireRoomPublishedEvent(t, publisher, roomEvents.RoomEventsListSucceeded{
		RoomID: roomEventTestRoomID,
		UserID: roomEventTestUserID,
		Count:  2,
	})
}

func TestEventPublishingRoomService_CreateDiceRollRoomEvent_Success(t *testing.T) {
	_, publisher, service := newRoomEventPublishingTestSubject()

	event, err := service.CreateDiceRollRoomEvent(context.Background(), model.CreateDiceRollRoomEventInput{
		RoomID:      roomTestUUID(roomEventTestRoomID),
		ActorID:     roomEventTestUserID,
		RollID:      "roll-1",
		CharacterID: roomEventTestCharacterID,
	})
	require.NoError(t, err)
	require.Equal(t, roomEventTestEventID, event.ID.String())
	requireRoomPublishedEvent(t, publisher, roomEvents.RoomDiceRollEventCreateSucceeded{
		RoomID:      roomEventTestRoomID,
		ActorID:     roomEventTestUserID,
		EventID:     roomEventTestEventID,
		RollID:      "roll-1",
		CharacterID: roomEventTestCharacterID,
	})
}

func TestEventPublishingRoomService_CreateCharacterChangedRoomEvents_Failure(t *testing.T) {
	next, publisher, service := newRoomEventPublishingTestSubject()
	next.err = errRoomEventTest

	_, err := service.CreateCharacterChangedRoomEvents(context.Background(), model.CreateCharacterChangedRoomEventsInput{
		ActorID:     roomEventTestUserID,
		CharacterID: roomTestUUID(roomEventTestCharacterID),
	})
	require.ErrorIs(t, err, errRoomEventTest)
	requireRoomPublishedEvent(t, publisher, roomEvents.RoomCharacterChangedEventsCreateFailed{
		ActorID:     roomEventTestUserID,
		CharacterID: roomEventTestCharacterID,
		Err:         errRoomEventTest,
	})
}

func TestEventPublishingRoomService_CleanupRooms_Success(t *testing.T) {
	_, publisher, service := newRoomEventPublishingTestSubject()

	result, err := service.CleanupRooms(context.Background(), model.CleanupRoomsInput{})
	require.NoError(t, err)
	require.Equal(t, 2, len(result.DeletedRoomIDs))
	requireRoomPublishedEvent(t, publisher, roomEvents.RoomCleanupSucceeded{
		InactiveDeleted: 1,
		InvalidDeleted:  1,
		DeletedCount:    2,
	})
}
