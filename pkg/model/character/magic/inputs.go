package magicDTO

import "github.com/jackc/pgx/v5/pgtype"

type GetMagicInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertMagicInput struct {
	UserID      string
	CharacterID pgtype.UUID
	MaxMp       *int16
	CurrentMp   *int16
}

type DeleteMagicInput struct {
	UserID      string
	CharacterID pgtype.UUID
}
