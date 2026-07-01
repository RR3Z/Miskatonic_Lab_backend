package notes

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type NoteModel struct {
	ID          pgtype.UUID        `json:"id"`
	CharacterID pgtype.UUID        `json:"character_id"`
	Title       string             `json:"title"`
	Body        string             `json:"body"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

func ToNoteModel(n db.Note) NoteModel {
	return NoteModel{
		ID:          n.ID,
		CharacterID: n.CharacterID,
		Title:       n.Title,
		Body:        n.Body,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
	}
}

func ToNoteModels(notes []db.Note) []NoteModel {
	models := make([]NoteModel, len(notes))
	for i, n := range notes {
		models[i] = ToNoteModel(n)
	}
	return models
}
