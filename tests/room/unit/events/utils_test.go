package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

const (
	roomEventTestRoomID      = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	roomEventTestMemberID    = "cccccccc-cccc-cccc-cccc-cccccccccccc"
	roomEventTestCharacterID = "dddddddd-dddd-dddd-dddd-dddddddddddd"
	roomEventTestEventID     = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	roomEventTestUserID      = "user_room_test"
	roomEventTestOwnerID     = "owner_room_test"
	roomEventTestTargetID    = "target_room_test"
)

func newRoomEventPublishingTestSubject() (*fakeEventPublishingRoomService, *fakeRoomEventPublisher, *roomService.EventPublishingRoomService) {
	next := &fakeEventPublishingRoomService{
		room: model.RoomModel{
			ID:      roomTestUUID(roomEventTestRoomID),
			OwnerID: roomEventTestOwnerID,
		},
		member: model.RoomMemberModel{
			ID:          roomTestUUID(roomEventTestMemberID),
			RoomID:      roomTestUUID(roomEventTestRoomID),
			UserID:      roomEventTestUserID,
			CharacterID: roomTestUUID(roomEventTestCharacterID),
			Role:        "player",
		},
		roomEvent: model.RoomEventModel{
			ID:      roomTestUUID(roomEventTestEventID),
			RoomID:  roomTestUUID(roomEventTestRoomID),
			ActorID: roomEventTestUserID,
		},
		roomEvents: []model.RoomEventModel{
			{ID: roomTestUUID(roomEventTestEventID), RoomID: roomTestUUID(roomEventTestRoomID), ActorID: roomEventTestUserID},
			{ID: roomTestUUID(roomEventTestMemberID), RoomID: roomTestUUID(roomEventTestRoomID), ActorID: roomEventTestOwnerID},
		},
		cleanupResult: model.CleanupRoomsResult{
			InactiveDeleted: 1,
			InvalidDeleted:  1,
			DeletedRoomIDs: []pgtype.UUID{
				roomTestUUID(roomEventTestRoomID),
				roomTestUUID(roomEventTestMemberID),
			},
		},
	}
	publisher := &fakeRoomEventPublisher{}
	return next, publisher, roomService.NewEventPublishingRoomService(next, next, publisher)
}

func roomTestUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	if err := uuid.Scan(value); err != nil {
		panic(err)
	}
	return uuid
}

func requireRoomPublishedEvent(t *testing.T, publisher *fakeRoomEventPublisher, expected events.Event) {
	t.Helper()

	require.Len(t, publisher.events, 1)
	require.Equal(t, expected, publisher.events[0])
}
