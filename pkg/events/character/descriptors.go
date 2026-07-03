package character

import "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"

const domain = "character"

var descriptors = []events.EventDescriptor{
	descriptor(CharactersListSucceeded{}, "characters", "list", events.OutcomeSucceeded),
	descriptor(CharactersListFailed{}, "characters", "list", events.OutcomeFailed),
	descriptor(CharacterGetSucceeded{}, "character", "get", events.OutcomeSucceeded),
	descriptor(CharacterGetFailed{}, "character", "get", events.OutcomeFailed),
	descriptor(CharacterCreateSucceeded{}, "character", "create", events.OutcomeSucceeded),
	descriptor(CharacterCreateFailed{}, "character", "create", events.OutcomeFailed),
	descriptor(CharacterUpdateSucceeded{}, "character", "update", events.OutcomeSucceeded),
	descriptor(CharacterUpdateFailed{}, "character", "update", events.OutcomeFailed),
	descriptor(CharacterDeleteSucceeded{}, "character", "delete", events.OutcomeSucceeded),
	descriptor(CharacterDeleteFailed{}, "character", "delete", events.OutcomeFailed),

	descriptor(CharacterHealthGetSucceeded{}, "health", "get", events.OutcomeSucceeded),
	descriptor(CharacterHealthGetFailed{}, "health", "get", events.OutcomeFailed),
	descriptor(CharacterHealthUpsertSucceeded{}, "health", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterHealthUpsertFailed{}, "health", "upsert", events.OutcomeFailed),
	descriptor(CharacterHealthDeleteSucceeded{}, "health", "delete", events.OutcomeSucceeded),
	descriptor(CharacterHealthDeleteFailed{}, "health", "delete", events.OutcomeFailed),

	descriptor(CharacterSanityGetSucceeded{}, "sanity", "get", events.OutcomeSucceeded),
	descriptor(CharacterSanityGetFailed{}, "sanity", "get", events.OutcomeFailed),
	descriptor(CharacterSanityUpsertSucceeded{}, "sanity", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterSanityUpsertFailed{}, "sanity", "upsert", events.OutcomeFailed),
	descriptor(CharacterSanityDeleteSucceeded{}, "sanity", "delete", events.OutcomeSucceeded),
	descriptor(CharacterSanityDeleteFailed{}, "sanity", "delete", events.OutcomeFailed),

	descriptor(CharacterMagicGetSucceeded{}, "magic", "get", events.OutcomeSucceeded),
	descriptor(CharacterMagicGetFailed{}, "magic", "get", events.OutcomeFailed),
	descriptor(CharacterMagicUpsertSucceeded{}, "magic", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterMagicUpsertFailed{}, "magic", "upsert", events.OutcomeFailed),
	descriptor(CharacterMagicDeleteSucceeded{}, "magic", "delete", events.OutcomeSucceeded),
	descriptor(CharacterMagicDeleteFailed{}, "magic", "delete", events.OutcomeFailed),

	descriptor(CharacterLuckGetSucceeded{}, "luck", "get", events.OutcomeSucceeded),
	descriptor(CharacterLuckGetFailed{}, "luck", "get", events.OutcomeFailed),
	descriptor(CharacterLuckUpsertSucceeded{}, "luck", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterLuckUpsertFailed{}, "luck", "upsert", events.OutcomeFailed),
	descriptor(CharacterLuckDeleteSucceeded{}, "luck", "delete", events.OutcomeSucceeded),
	descriptor(CharacterLuckDeleteFailed{}, "luck", "delete", events.OutcomeFailed),

	descriptor(CharacterFinancesGetSucceeded{}, "finances", "get", events.OutcomeSucceeded),
	descriptor(CharacterFinancesGetFailed{}, "finances", "get", events.OutcomeFailed),
	descriptor(CharacterFinancesUpsertSucceeded{}, "finances", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterFinancesUpsertFailed{}, "finances", "upsert", events.OutcomeFailed),
	descriptor(CharacterFinancesDeleteSucceeded{}, "finances", "delete", events.OutcomeSucceeded),
	descriptor(CharacterFinancesDeleteFailed{}, "finances", "delete", events.OutcomeFailed),

	descriptor(CharacterDerivedStatsGetSucceeded{}, "derived_stats", "get", events.OutcomeSucceeded),
	descriptor(CharacterDerivedStatsGetFailed{}, "derived_stats", "get", events.OutcomeFailed),
	descriptor(CharacterDerivedStatsUpsertSucceeded{}, "derived_stats", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterDerivedStatsUpsertFailed{}, "derived_stats", "upsert", events.OutcomeFailed),
	descriptor(CharacterDerivedStatsAutoRecalculateSucceeded{}, "derived_stats", "auto_recalculate", events.OutcomeSucceeded),
	descriptor(CharacterDerivedStatsAutoRecalculateSkipped{}, "derived_stats", "auto_recalculate", events.OutcomeSkipped),
	descriptor(CharacterDerivedStatsAutoRecalculateFailed{}, "derived_stats", "auto_recalculate", events.OutcomeFailed),
	descriptor(CharacterDerivedStatsDeleteSucceeded{}, "derived_stats", "delete", events.OutcomeSucceeded),
	descriptor(CharacterDerivedStatsDeleteFailed{}, "derived_stats", "delete", events.OutcomeFailed),

	descriptor(CharacterCharacteristicsGetSucceeded{}, "characteristics", "get", events.OutcomeSucceeded),
	descriptor(CharacterCharacteristicsGetFailed{}, "characteristics", "get", events.OutcomeFailed),
	descriptor(CharacterCharacteristicsUpsertSucceeded{}, "characteristics", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterCharacteristicsUpsertFailed{}, "characteristics", "upsert", events.OutcomeFailed),
	descriptor(CharacterCharacteristicsDeleteSucceeded{}, "characteristics", "delete", events.OutcomeSucceeded),
	descriptor(CharacterCharacteristicsDeleteFailed{}, "characteristics", "delete", events.OutcomeFailed),

	descriptor(CharacterBackstoryGetSucceeded{}, "backstory", "get", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryGetFailed{}, "backstory", "get", events.OutcomeFailed),
	descriptor(CharacterBackstoryUpsertSucceeded{}, "backstory", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryUpsertFailed{}, "backstory", "upsert", events.OutcomeFailed),
	descriptor(CharacterBackstoryDeleteSucceeded{}, "backstory", "delete", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryDeleteFailed{}, "backstory", "delete", events.OutcomeFailed),
	descriptor(CharacterBackstoryItemsListSucceeded{}, "backstory_items", "list", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryItemsListFailed{}, "backstory_items", "list", events.OutcomeFailed),
	descriptor(CharacterBackstoryItemGetSucceeded{}, "backstory_item", "get", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryItemGetFailed{}, "backstory_item", "get", events.OutcomeFailed),
	descriptor(CharacterBackstoryItemCreateSucceeded{}, "backstory_item", "create", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryItemCreateFailed{}, "backstory_item", "create", events.OutcomeFailed),
	descriptor(CharacterBackstoryItemUpdateSucceeded{}, "backstory_item", "update", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryItemUpdateFailed{}, "backstory_item", "update", events.OutcomeFailed),
	descriptor(CharacterBackstoryItemDeleteSucceeded{}, "backstory_item", "delete", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryItemDeleteFailed{}, "backstory_item", "delete", events.OutcomeFailed),

	descriptor(CharacterSkillsListSucceeded{}, "skills", "list", events.OutcomeSucceeded),
	descriptor(CharacterSkillsListFailed{}, "skills", "list", events.OutcomeFailed),
	descriptor(CharacterSkillGetSucceeded{}, "skill", "get", events.OutcomeSucceeded),
	descriptor(CharacterSkillGetFailed{}, "skill", "get", events.OutcomeFailed),
	descriptor(CharacterSkillCreateSucceeded{}, "skill", "create", events.OutcomeSucceeded),
	descriptor(CharacterSkillCreateFailed{}, "skill", "create", events.OutcomeFailed),
	descriptor(CharacterSkillUpdateSucceeded{}, "skill", "update", events.OutcomeSucceeded),
	descriptor(CharacterSkillUpdateFailed{}, "skill", "update", events.OutcomeFailed),
	descriptor(CharacterSkillDeleteSucceeded{}, "skill", "delete", events.OutcomeSucceeded),
	descriptor(CharacterSkillDeleteFailed{}, "skill", "delete", events.OutcomeFailed),

	descriptor(CharacterNotesListSucceeded{}, "notes", "list", events.OutcomeSucceeded),
	descriptor(CharacterNotesListFailed{}, "notes", "list", events.OutcomeFailed),
	descriptor(CharacterNoteGetSucceeded{}, "note", "get", events.OutcomeSucceeded),
	descriptor(CharacterNoteGetFailed{}, "note", "get", events.OutcomeFailed),
	descriptor(CharacterNoteCreateSucceeded{}, "note", "create", events.OutcomeSucceeded),
	descriptor(CharacterNoteCreateFailed{}, "note", "create", events.OutcomeFailed),
	descriptor(CharacterNoteUpdateSucceeded{}, "note", "update", events.OutcomeSucceeded),
	descriptor(CharacterNoteUpdateFailed{}, "note", "update", events.OutcomeFailed),
	descriptor(CharacterNoteDeleteSucceeded{}, "note", "delete", events.OutcomeSucceeded),
	descriptor(CharacterNoteDeleteFailed{}, "note", "delete", events.OutcomeFailed),
}

var roomMutationDescriptors = []events.EventDescriptor{
	descriptor(CharacterUpdateSucceeded{}, "character", "update", events.OutcomeSucceeded),
	descriptor(CharacterHealthUpsertSucceeded{}, "health", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterHealthDeleteSucceeded{}, "health", "delete", events.OutcomeSucceeded),
	descriptor(CharacterSanityUpsertSucceeded{}, "sanity", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterSanityDeleteSucceeded{}, "sanity", "delete", events.OutcomeSucceeded),
	descriptor(CharacterMagicUpsertSucceeded{}, "magic", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterMagicDeleteSucceeded{}, "magic", "delete", events.OutcomeSucceeded),
	descriptor(CharacterLuckUpsertSucceeded{}, "luck", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterLuckDeleteSucceeded{}, "luck", "delete", events.OutcomeSucceeded),
	descriptor(CharacterCharacteristicsUpsertSucceeded{}, "characteristics", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterCharacteristicsDeleteSucceeded{}, "characteristics", "delete", events.OutcomeSucceeded),
	descriptor(CharacterDerivedStatsUpsertSucceeded{}, "derived_stats", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterDerivedStatsDeleteSucceeded{}, "derived_stats", "delete", events.OutcomeSucceeded),
	descriptor(CharacterFinancesUpsertSucceeded{}, "finances", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterFinancesDeleteSucceeded{}, "finances", "delete", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryUpsertSucceeded{}, "backstory", "upsert", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryDeleteSucceeded{}, "backstory", "delete", events.OutcomeSucceeded),
	descriptor(CharacterSkillCreateSucceeded{}, "skill", "create", events.OutcomeSucceeded),
	descriptor(CharacterSkillUpdateSucceeded{}, "skill", "update", events.OutcomeSucceeded),
	descriptor(CharacterSkillDeleteSucceeded{}, "skill", "delete", events.OutcomeSucceeded),
	descriptor(CharacterNoteCreateSucceeded{}, "note", "create", events.OutcomeSucceeded),
	descriptor(CharacterNoteUpdateSucceeded{}, "note", "update", events.OutcomeSucceeded),
	descriptor(CharacterNoteDeleteSucceeded{}, "note", "delete", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryItemCreateSucceeded{}, "backstory_item", "create", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryItemUpdateSucceeded{}, "backstory_item", "update", events.OutcomeSucceeded),
	descriptor(CharacterBackstoryItemDeleteSucceeded{}, "backstory_item", "delete", events.OutcomeSucceeded),
}

func Descriptors() []events.EventDescriptor {
	return append([]events.EventDescriptor(nil), descriptors...)
}

func RoomMutationEvents() []events.Event {
	return events.EventPrototypes(roomMutationDescriptors)
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
