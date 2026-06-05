package character

type CharacterSanityGetSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterSanityGetFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterSanityGetSucceeded) EventName() string {
	return "character.sanity.get_succeeded"
}

func (CharacterSanityGetFailed) EventName() string {
	return "character.sanity.get_failed"
}

type CharacterSanityUpsertSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterSanityUpsertFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterSanityUpsertSucceeded) EventName() string {
	return "character.sanity.upsert_succeeded"
}

func (CharacterSanityUpsertFailed) EventName() string {
	return "character.sanity.upsert_failed"
}

type CharacterSanityDeleteSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterSanityDeleteFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterSanityDeleteSucceeded) EventName() string {
	return "character.sanity.delete_succeeded"
}

func (CharacterSanityDeleteFailed) EventName() string {
	return "character.sanity.delete_failed"
}
