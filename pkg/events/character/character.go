package character

// Every struct here implement Event

type CharactersListSucceeded struct {
	UserID string
	Count  int
}

type CharactersListFailed struct {
	UserID string
	Err    error
}

func (CharactersListSucceeded) EventName() string {
	return "characters.list_succeeded"
}

func (CharactersListFailed) EventName() string {
	return "characters.list_failed"
}

type CharacterGetSucceeded struct {
	UserID      string
	CharacterID string
	Name        string
}

type CharacterGetFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterGetSucceeded) EventName() string {
	return "character.get_succeeded"
}

func (CharacterGetFailed) EventName() string {
	return "character.get_failed"
}

type CharacterCreateSucceeded struct {
	UserID      string
	CharacterID string
	Name        string
}

type CharacterCreateFailed struct {
	UserID string
	Err    error
}

func (CharacterCreateSucceeded) EventName() string {
	return "character.create_succeeded"
}

func (CharacterCreateFailed) EventName() string {
	return "character.create_failed"
}

type CharacterUpdateSucceeded struct {
	UserID      string
	CharacterID string
	Name        string
}

type CharacterUpdateFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterUpdateSucceeded) EventName() string {
	return "character.update_succeeded"
}

func (CharacterUpdateFailed) EventName() string {
	return "character.update_failed"
}

type CharacterDeleteSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterDeleteFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterDeleteSucceeded) EventName() string {
	return "character.delete_succeeded"
}

func (CharacterDeleteFailed) EventName() string {
	return "character.delete_failed"
}
