package notesDTO

import "github.com/jackc/pgx/v5/pgtype"

type GetNotesInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetNoteInput struct {
	UserID      string
	CharacterID pgtype.UUID
	NoteID      pgtype.UUID
}

type CreateNoteInput struct {
	Title       string
	Body        string
	UserID      string
	CharacterID pgtype.UUID
}

type UpdateNoteInput struct {
	Title       string
	Body        string
	UserID      string
	CharacterID pgtype.UUID
	NoteID      pgtype.UUID
}

type DeleteNoteInput struct {
	UserID      string
	CharacterID pgtype.UUID
	NoteID      pgtype.UUID
}
