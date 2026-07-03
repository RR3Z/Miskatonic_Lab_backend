package tests

import (
	"context"
	"testing"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/listeners"
	"github.com/stretchr/testify/require"
)

func TestCharacterRoomListener_HealthUpsert_CreatesCharacterChangedRoomEvent(t *testing.T) {
	svc := &fakeListenerRoomService{}
	listener := listeners.NewCharacterRoomListener(svc)

	listener.Handle(context.Background(), characterEvents.CharacterHealthUpsertSucceeded{
		UserID:      "user_1",
		CharacterID: "11111111-1111-1111-1111-111111111111",
	})

	require.Equal(t, 1, svc.characterCalls)
	require.Equal(t, "user_1", svc.characterInput.ActorID)
	require.Equal(t, listenerTestUUID("11111111-1111-1111-1111-111111111111"), svc.characterInput.CharacterID)
	require.Equal(t, "health", svc.characterInput.Change.Resource)
	require.Equal(t, "upsert", svc.characterInput.Change.Action)
	require.Nil(t, svc.characterInput.Change.ResourceID)
	require.NotNil(t, svc.characterInput.Change.SourceEvent)
	require.Equal(t, "character.health.upsert_succeeded", *svc.characterInput.Change.SourceEvent)
}

func TestCharacterRoomListener_SkillUpdate_CreatesCharacterChangedRoomEvent(t *testing.T) {
	svc := &fakeListenerRoomService{}
	listener := listeners.NewCharacterRoomListener(svc)

	listener.Handle(context.Background(), characterEvents.CharacterSkillUpdateSucceeded{
		UserID:      "user_1",
		CharacterID: "11111111-1111-1111-1111-111111111111",
		SkillID:     "22222222-2222-2222-2222-222222222222",
		Name:        "Library Use",
	})

	require.Equal(t, 1, svc.characterCalls)
	require.Equal(t, "skill", svc.characterInput.Change.Resource)
	require.Equal(t, "update", svc.characterInput.Change.Action)
	require.NotNil(t, svc.characterInput.Change.ResourceID)
	require.Equal(t, "22222222-2222-2222-2222-222222222222", *svc.characterInput.Change.ResourceID)
	require.NotNil(t, svc.characterInput.Change.SourceEvent)
	require.Equal(t, "character.skill.update_succeeded", *svc.characterInput.Change.SourceEvent)
}

func TestCharacterRoomListener_ReadAndListEventsIgnored(t *testing.T) {
	svc := &fakeListenerRoomService{}
	listener := listeners.NewCharacterRoomListener(svc)

	listener.Handle(context.Background(), characterEvents.CharacterGetSucceeded{
		UserID:      "user_1",
		CharacterID: "11111111-1111-1111-1111-111111111111",
		Name:        "Investigator",
	})
	listener.Handle(context.Background(), characterEvents.CharacterSkillsListSucceeded{
		UserID:      "user_1",
		CharacterID: "11111111-1111-1111-1111-111111111111",
		Count:       1,
	})
	listener.Handle(context.Background(), characterEvents.CharactersListSucceeded{
		UserID: "user_1",
		Count:  1,
	})

	require.Zero(t, svc.characterCalls)
}

func TestCharacterRoomListener_CharacterDeleteIgnored(t *testing.T) {
	svc := &fakeListenerRoomService{}
	listener := listeners.NewCharacterRoomListener(svc)

	listener.Handle(context.Background(), characterEvents.CharacterDeleteSucceeded{
		UserID:      "user_1",
		CharacterID: "11111111-1111-1111-1111-111111111111",
	})

	require.Zero(t, svc.characterCalls)
}
