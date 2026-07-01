package notes

import "github.com/jackc/pgx/v5/pgtype"

type NoteModel struct {
	ID          pgtype.UUID        `json:"id"`
	CharacterID pgtype.UUID        `json:"character_id"`
	Title       string             `json:"title"`
	Body        string             `json:"body"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}
