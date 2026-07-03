package roomlisteners

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
)

func MutationCharacterEvents() []events.Event {
	return characterEvents.RoomMutationEvents()
}

func CharacterChangedRoomEventInput(event events.Event) (string, string, model.CharacterChangedRoomEventChange, bool) {
	sourceEvent := event.EventName()

	switch e := event.(type) {
	case characterEvents.CharacterUpdateSucceeded:
		return e.UserID, e.CharacterID, change("character", "update", nil, &sourceEvent), true
	case characterEvents.CharacterHealthUpsertSucceeded:
		return e.UserID, e.CharacterID, change("health", "upsert", nil, &sourceEvent), true
	case characterEvents.CharacterHealthDeleteSucceeded:
		return e.UserID, e.CharacterID, change("health", "delete", nil, &sourceEvent), true
	case characterEvents.CharacterSanityUpsertSucceeded:
		return e.UserID, e.CharacterID, change("sanity", "upsert", nil, &sourceEvent), true
	case characterEvents.CharacterSanityDeleteSucceeded:
		return e.UserID, e.CharacterID, change("sanity", "delete", nil, &sourceEvent), true
	case characterEvents.CharacterMagicUpsertSucceeded:
		return e.UserID, e.CharacterID, change("magic", "upsert", nil, &sourceEvent), true
	case characterEvents.CharacterMagicDeleteSucceeded:
		return e.UserID, e.CharacterID, change("magic", "delete", nil, &sourceEvent), true
	case characterEvents.CharacterLuckUpsertSucceeded:
		return e.UserID, e.CharacterID, change("luck", "upsert", nil, &sourceEvent), true
	case characterEvents.CharacterLuckDeleteSucceeded:
		return e.UserID, e.CharacterID, change("luck", "delete", nil, &sourceEvent), true
	case characterEvents.CharacterCharacteristicsUpsertSucceeded:
		return e.UserID, e.CharacterID, change("characteristics", "upsert", nil, &sourceEvent), true
	case characterEvents.CharacterCharacteristicsDeleteSucceeded:
		return e.UserID, e.CharacterID, change("characteristics", "delete", nil, &sourceEvent), true
	case characterEvents.CharacterDerivedStatsUpsertSucceeded:
		return e.UserID, e.CharacterID, change("derived_stats", "upsert", nil, &sourceEvent), true
	case characterEvents.CharacterDerivedStatsDeleteSucceeded:
		return e.UserID, e.CharacterID, change("derived_stats", "delete", nil, &sourceEvent), true
	case characterEvents.CharacterFinancesUpsertSucceeded:
		return e.UserID, e.CharacterID, change("finances", "upsert", nil, &sourceEvent), true
	case characterEvents.CharacterFinancesDeleteSucceeded:
		return e.UserID, e.CharacterID, change("finances", "delete", nil, &sourceEvent), true
	case characterEvents.CharacterBackstoryUpsertSucceeded:
		return e.UserID, e.CharacterID, change("backstory", "upsert", nil, &sourceEvent), true
	case characterEvents.CharacterBackstoryDeleteSucceeded:
		return e.UserID, e.CharacterID, change("backstory", "delete", nil, &sourceEvent), true
	case characterEvents.CharacterSkillCreateSucceeded:
		return e.UserID, e.CharacterID, change("skill", "create", &e.SkillID, &sourceEvent), true
	case characterEvents.CharacterSkillUpdateSucceeded:
		return e.UserID, e.CharacterID, change("skill", "update", &e.SkillID, &sourceEvent), true
	case characterEvents.CharacterSkillDeleteSucceeded:
		return e.UserID, e.CharacterID, change("skill", "delete", &e.SkillID, &sourceEvent), true
	case characterEvents.CharacterNoteCreateSucceeded:
		return e.UserID, e.CharacterID, change("note", "create", &e.NoteID, &sourceEvent), true
	case characterEvents.CharacterNoteUpdateSucceeded:
		return e.UserID, e.CharacterID, change("note", "update", &e.NoteID, &sourceEvent), true
	case characterEvents.CharacterNoteDeleteSucceeded:
		return e.UserID, e.CharacterID, change("note", "delete", &e.NoteID, &sourceEvent), true
	case characterEvents.CharacterBackstoryItemCreateSucceeded:
		return e.UserID, e.CharacterID, change("backstory_item", "create", &e.BackstoryItemID, &sourceEvent), true
	case characterEvents.CharacterBackstoryItemUpdateSucceeded:
		return e.UserID, e.CharacterID, change("backstory_item", "update", &e.BackstoryItemID, &sourceEvent), true
	case characterEvents.CharacterBackstoryItemDeleteSucceeded:
		return e.UserID, e.CharacterID, change("backstory_item", "delete", &e.BackstoryItemID, &sourceEvent), true
	default:
		return "", "", model.CharacterChangedRoomEventChange{}, false
	}
}

func change(resource string, action string, resourceID *string, sourceEvent *string) model.CharacterChangedRoomEventChange {
	return model.CharacterChangedRoomEventChange{
		Resource:    resource,
		Action:      action,
		ResourceID:  resourceID,
		SourceEvent: sourceEvent,
	}
}
