package luck

import "github.com/jackc/pgx/v5/pgtype"

type GetLuckInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertLuckInput struct {
	UserID       string
	CharacterID  pgtype.UUID
	StartingLuck *int16
	CurrentLuck  *int16
}

type DeleteLuckInput struct {
	UserID      string
	CharacterID pgtype.UUID
}
