package sanityDTO

import "github.com/jackc/pgx/v5/pgtype"

type GetSanityInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertSanityInput struct {
	UserID        string
	CharacterID   pgtype.UUID
	MaxSanity     *int16
	CurrentSanity *int16
	TempInsanity  *bool
	IndefInsanity *bool
}

type DeleteSanityInput struct {
	UserID      string
	CharacterID pgtype.UUID
}
