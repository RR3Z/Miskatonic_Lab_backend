package character

type CharacterMagicGetSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterMagicGetFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterMagicGetSucceeded) EventName() string {
	return "character.magic.get_succeeded"
}

func (CharacterMagicGetFailed) EventName() string {
	return "character.magic.get_failed"
}

type CharacterMagicUpsertSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterMagicUpsertFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterMagicUpsertSucceeded) EventName() string {
	return "character.magic.upsert_succeeded"
}

func (CharacterMagicUpsertFailed) EventName() string {
	return "character.magic.upsert_failed"
}

type CharacterMagicDeleteSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterMagicDeleteFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterMagicDeleteSucceeded) EventName() string {
	return "character.magic.delete_succeeded"
}

func (CharacterMagicDeleteFailed) EventName() string {
	return "character.magic.delete_failed"
}
