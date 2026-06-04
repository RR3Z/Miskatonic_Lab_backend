package character

type CharacterCharacteristicsGetSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterCharacteristicsGetFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterCharacteristicsGetSucceeded) EventName() string {
	return "character.characteristics.get_succeeded"
}

func (CharacterCharacteristicsGetFailed) EventName() string {
	return "character.characteristics.get_failed"
}

type CharacterCharacteristicsUpsertSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterCharacteristicsUpsertFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterCharacteristicsUpsertSucceeded) EventName() string {
	return "character.characteristics.upsert_succeeded"
}

func (CharacterCharacteristicsUpsertFailed) EventName() string {
	return "character.characteristics.upsert_failed"
}

type CharacterCharacteristicsDeleteSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterCharacteristicsDeleteFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterCharacteristicsDeleteSucceeded) EventName() string {
	return "character.characteristics.delete_succeeded"
}

func (CharacterCharacteristicsDeleteFailed) EventName() string {
	return "character.characteristics.delete_failed"
}
