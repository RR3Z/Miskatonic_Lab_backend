package character

type CharacterFinancesGetSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterFinancesGetFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterFinancesGetSucceeded) EventName() string {
	return "character.finances.get_succeeded"
}

func (CharacterFinancesGetFailed) EventName() string {
	return "character.finances.get_failed"
}

type CharacterFinancesUpsertSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterFinancesUpsertFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterFinancesUpsertSucceeded) EventName() string {
	return "character.finances.upsert_succeeded"
}

func (CharacterFinancesUpsertFailed) EventName() string {
	return "character.finances.upsert_failed"
}

type CharacterFinancesDeleteSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterFinancesDeleteFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterFinancesDeleteSucceeded) EventName() string {
	return "character.finances.delete_succeeded"
}

func (CharacterFinancesDeleteFailed) EventName() string {
	return "character.finances.delete_failed"
}
