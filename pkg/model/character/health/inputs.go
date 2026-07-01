package healthDTO

import "github.com/jackc/pgx/v5/pgtype"

type GetHealthInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertHealthInput struct {
	UserID      string
	CharacterID pgtype.UUID
	MaxHp       *int16
	CurrentHp   *int16
	MajorWound  *bool
	Unconscious *bool
	Dying       *bool
	Dead        *bool
}

type DeleteHealthInput struct {
	UserID      string
	CharacterID pgtype.UUID
}
