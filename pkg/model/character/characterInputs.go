package characterDTO

import "github.com/jackc/pgx/v5/pgtype"

type GetCharacterInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type CreateCharacterInput struct {
	UserID      string
	Name        string
	PlayerName  *string
	Occupation  *string
	Age         *int16
	Sex         *string
	Residence   *string
	Birthplace  *string
	PortraitUrl *string
}

type UpdateCharacterInput struct {
	UserID      string
	ID          pgtype.UUID
	Name        string
	PlayerName  *string
	Occupation  *string
	Age         *int16
	Sex         *string
	Residence   *string
	Birthplace  *string
	PortraitUrl *string
}

type DeleteCharacterInput struct {
	UserID string
	ID     pgtype.UUID
}
