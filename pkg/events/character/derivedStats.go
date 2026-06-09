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

type CharacterDerivedStatsAutoRecalculateSucceeded struct {
	UserID      string
	CharacterID string
	Source      string
}

type CharacterDerivedStatsAutoRecalculateSkipped struct {
	UserID      string
	CharacterID string
	Source      string
	Reason      string
}

type CharacterDerivedStatsAutoRecalculateFailed struct {
	UserID      string
	CharacterID string
	Source      string
	Err         error
}

func (CharacterDerivedStatsAutoRecalculateSucceeded) EventName() string {
	return "character.derived_stats.auto_recalculate_succeeded"
}

func (CharacterDerivedStatsAutoRecalculateSkipped) EventName() string {
	return "character.derived_stats.auto_recalculate_skipped"
}

func (CharacterDerivedStatsAutoRecalculateFailed) EventName() string {
	return "character.derived_stats.auto_recalculate_failed"
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
