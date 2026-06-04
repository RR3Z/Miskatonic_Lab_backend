package character

type CharacterNotesListSucceeded struct {
	UserID      string
	CharacterID string
	Count       int
}

type CharacterNotesListFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterNotesListSucceeded) EventName() string {
	return "character.notes.list_succeeded"
}

func (CharacterNotesListFailed) EventName() string {
	return "character.notes.list_failed"
}

type CharacterNoteGetSucceeded struct {
	UserID      string
	CharacterID string
	NoteID      string
	Title       string
}

type CharacterNoteGetFailed struct {
	UserID      string
	CharacterID string
	NoteID      string
	Err         error
}

func (CharacterNoteGetSucceeded) EventName() string {
	return "character.note.get_succeeded"
}

func (CharacterNoteGetFailed) EventName() string {
	return "character.note.get_failed"
}

type CharacterNoteCreateSucceeded struct {
	UserID      string
	CharacterID string
	NoteID      string
	Title       string
}

type CharacterNoteCreateFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterNoteCreateSucceeded) EventName() string {
	return "character.note.create_succeeded"
}

func (CharacterNoteCreateFailed) EventName() string {
	return "character.note.create_failed"
}

type CharacterNoteUpdateSucceeded struct {
	UserID      string
	CharacterID string
	NoteID      string
	Title       string
}

type CharacterNoteUpdateFailed struct {
	UserID      string
	CharacterID string
	NoteID      string
	Err         error
}

func (CharacterNoteUpdateSucceeded) EventName() string {
	return "character.note.update_succeeded"
}

func (CharacterNoteUpdateFailed) EventName() string {
	return "character.note.update_failed"
}

type CharacterNoteDeleteSucceeded struct {
	UserID      string
	CharacterID string
	NoteID      string
}

type CharacterNoteDeleteFailed struct {
	UserID      string
	CharacterID string
	NoteID      string
	Err         error
}

func (CharacterNoteDeleteSucceeded) EventName() string {
	return "character.note.delete_succeeded"
}

func (CharacterNoteDeleteFailed) EventName() string {
	return "character.note.delete_failed"
}
