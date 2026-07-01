package backstories

import "github.com/jackc/pgx/v5/pgtype"

type BackstoryItemModel struct {
	ID pgtype.UUID `json:"id"`

	Section string `json:"section"`
	Title   string `json:"title"`
	Text    string `json:"text"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type BackstoryModel struct {
	ID          pgtype.UUID `json:"id"`
	CharacterID pgtype.UUID `json:"character_id"`

	PersonalDescription *string              `json:"personal_description"`
	Items               []BackstoryItemModel `json:"items"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}
