package characteristics

import "github.com/jackc/pgx/v5/pgtype"

type GetCharacteristicsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertCharacteristicsInput struct {
	Strength     *int16
	Constitution *int16
	Size         *int16
	Dexterity    *int16
	Appearance   *int16
	Intelligence *int16
	Power        *int16
	Education    *int16
	UserID       string
	CharacterID  pgtype.UUID
}

type DeleteCharacteristicsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}
