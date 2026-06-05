package character

type CharacterLuckGetSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterLuckGetFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterLuckGetSucceeded) EventName() string {
	return "character.luck.get_succeeded"
}

func (CharacterLuckGetFailed) EventName() string {
	return "character.luck.get_failed"
}

type CharacterLuckUpsertSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterLuckUpsertFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterLuckUpsertSucceeded) EventName() string {
	return "character.luck.upsert_succeeded"
}

func (CharacterLuckUpsertFailed) EventName() string {
	return "character.luck.upsert_failed"
}

type CharacterLuckDeleteSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterLuckDeleteFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterLuckDeleteSucceeded) EventName() string {
	return "character.luck.delete_succeeded"
}

func (CharacterLuckDeleteFailed) EventName() string {
	return "character.luck.delete_failed"
}
