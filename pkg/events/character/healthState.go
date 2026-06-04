package character

type CharacterHealthGetSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterHealthGetFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterHealthGetSucceeded) EventName() string {
	return "character.health.get_succeeded"
}

func (CharacterHealthGetFailed) EventName() string {
	return "character.health.get_failed"
}

type CharacterHealthUpsertSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterHealthUpsertFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterHealthUpsertSucceeded) EventName() string {
	return "character.health.upsert_succeeded"
}

func (CharacterHealthUpsertFailed) EventName() string {
	return "character.health.upsert_failed"
}

type CharacterHealthDeleteSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterHealthDeleteFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterHealthDeleteSucceeded) EventName() string {
	return "character.health.delete_succeeded"
}

func (CharacterHealthDeleteFailed) EventName() string {
	return "character.health.delete_failed"
}
