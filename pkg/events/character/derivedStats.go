package character

type CharacterDerivedStatsGetSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterDerivedStatsGetFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterDerivedStatsGetSucceeded) EventName() string {
	return "character.derived_stats.get_succeeded"
}

func (CharacterDerivedStatsGetFailed) EventName() string {
	return "character.derived_stats.get_failed"
}

type CharacterDerivedStatsUpsertSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterDerivedStatsUpsertFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterDerivedStatsUpsertSucceeded) EventName() string {
	return "character.derived_stats.upsert_succeeded"
}

func (CharacterDerivedStatsUpsertFailed) EventName() string {
	return "character.derived_stats.upsert_failed"
}

type CharacterDerivedStatsDeleteSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterDerivedStatsDeleteFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterDerivedStatsDeleteSucceeded) EventName() string {
	return "character.derived_stats.delete_succeeded"
}

func (CharacterDerivedStatsDeleteFailed) EventName() string {
	return "character.derived_stats.delete_failed"
}
