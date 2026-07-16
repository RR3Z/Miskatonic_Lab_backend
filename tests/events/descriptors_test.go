package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	"github.com/stretchr/testify/require"
)

func TestDescriptorRegistryDescribesValueAndPointerEvents(t *testing.T) {
	registry := events.NewDescriptorRegistry(
		characterEvents.Descriptors(),
		diceEvents.Descriptors(),
		roomEvents.Descriptors(),
	)

	descriptor, ok := registry.Describe(&diceEvents.DiceRollMakeSucceeded{})
	require.True(t, ok)
	require.Equal(t, "dice", descriptor.Domain)
	require.Equal(t, "dice_roll", descriptor.Resource)
	require.Equal(t, "make", descriptor.Action)
	require.Equal(t, events.OutcomeSucceeded, descriptor.Outcome)

	descriptor, ok = registry.Describe(characterEvents.CharacterHealthUpsertFailed{})
	require.True(t, ok)
	require.Equal(t, "character", descriptor.Domain)
	require.Equal(t, "health", descriptor.Resource)
	require.Equal(t, "upsert", descriptor.Action)
	require.Equal(t, events.OutcomeFailed, descriptor.Outcome)
}

func TestDescriptorRegistryIgnoresUnknownEvents(t *testing.T) {
	registry := events.NewDescriptorRegistry(characterEvents.Descriptors())

	_, ok := registry.Describe(testEvent{name: "unknown.event"})
	require.False(t, ok)
}

func TestCharacterRoomMutationEventsAreDescribed(t *testing.T) {
	registry := events.NewDescriptorRegistry(characterEvents.Descriptors())
	mutations := characterEvents.RoomMutationEvents()

	require.Len(t, mutations, 25)
	requireEventPrototype(t, mutations, characterEvents.CharacterPortraitReplaceSucceeded{})
	requireEventPrototype(t, mutations, characterEvents.CharacterHealthUpsertSucceeded{})
	requireEventPrototype(t, mutations, characterEvents.CharacterBackstoryItemDeleteSucceeded{})

	for _, event := range mutations {
		descriptor, ok := registry.Describe(event)
		require.True(t, ok, event.EventName())
		require.Equal(t, "character", descriptor.Domain)
		require.Equal(t, events.OutcomeSucceeded, descriptor.Outcome)
	}
}

func TestDiceRoomPublishingEventsAreDescribed(t *testing.T) {
	registry := events.NewDescriptorRegistry(diceEvents.Descriptors())
	publishingEvents := diceEvents.RoomPublishingEvents()

	require.Len(t, publishingEvents, 1)
	requireEventPrototype(t, publishingEvents, diceEvents.DiceRollMakeSucceeded{})

	descriptor, ok := registry.Describe(publishingEvents[0])
	require.True(t, ok)
	require.Equal(t, "dice", descriptor.Domain)
	require.Equal(t, "dice_roll", descriptor.Resource)
	require.Equal(t, "make", descriptor.Action)
	require.Equal(t, events.OutcomeSucceeded, descriptor.Outcome)
}

func TestRoomDescriptorsCoverRoomDomainEvents(t *testing.T) {
	registry := events.NewDescriptorRegistry(roomEvents.Descriptors())

	descriptor, ok := registry.Describe(roomEvents.RoomCleanupSucceeded{})
	require.True(t, ok)
	require.Equal(t, "room", descriptor.Domain)
	require.Equal(t, "room_cleanup", descriptor.Resource)
	require.Equal(t, "cleanup", descriptor.Action)
	require.Equal(t, events.OutcomeSucceeded, descriptor.Outcome)

	descriptor, ok = registry.Describe(roomEvents.RoomMemberJoinFailed{})
	require.True(t, ok)
	require.Equal(t, "room_member", descriptor.Resource)
	require.Equal(t, "join", descriptor.Action)
	require.Equal(t, events.OutcomeFailed, descriptor.Outcome)
	require.Len(t, roomEvents.AllEvents(), len(roomEvents.Descriptors()))
}

func requireEventPrototype(t *testing.T, prototypes []events.Event, expected events.Event) {
	t.Helper()

	for _, prototype := range prototypes {
		if events.EventTypeOf(prototype) == events.EventTypeOf(expected) {
			return
		}
	}
	require.Failf(t, "missing event prototype", "expected %s", expected.EventName())
}
