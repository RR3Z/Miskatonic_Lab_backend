package characterDTO

import (
	"io"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetCharacterInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type CreateCharacterInput struct {
	UserID     string
	Name       string
	Occupation *string
	Age        *int16
	Sex        *string
	Residence  *string
	Birthplace *string
}

type UpdateCharacterInput struct {
	UserID     string
	ID         pgtype.UUID
	Name       string
	Occupation *string
	Age        *int16
	Sex        *string
	Residence  *string
	Birthplace *string
}

type PatchCharacterInput struct {
	UserID string
	ID     pgtype.UUID

	Name       PatchValue[string]
	Occupation PatchValue[string]
	Age        PatchValue[int16]
	Sex        PatchValue[string]
	Residence  PatchValue[string]
	Birthplace PatchValue[string]
}

func (i PatchCharacterInput) HasChanges() bool {
	return i.Name.Set ||
		i.Occupation.Set ||
		i.Age.Set ||
		i.Sex.Set ||
		i.Residence.Set ||
		i.Birthplace.Set
}

type ReplacePortraitInput struct {
	UserID      string
	CharacterID pgtype.UUID
	File        io.Reader
}

type DeleteCharacterInput struct {
	UserID string
	ID     pgtype.UUID
}
