package backstoriesDTO

import "github.com/jackc/pgx/v5/pgtype"

type GetBackstoryInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type UpsertBackstoryInput struct {
	UserID              string
	CharacterID         pgtype.UUID
	PersonalDescription *string
}

type DeleteBackstoryInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetBackstoryItemsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetBackstoryItemInput struct {
	UserID          string
	CharacterID     pgtype.UUID
	BackstoryItemID pgtype.UUID
}

type CreateBackstoryItemInput struct {
	Section     string
	Title       string
	Text        string
	UserID      string
	CharacterID pgtype.UUID
}

type UpdateBackstoryItemInput struct {
	Section         string
	Title           string
	Text            string
	UserID          string
	CharacterID     pgtype.UUID
	BackstoryItemID pgtype.UUID
}

type DeleteBackstoryItemInput struct {
	UserID          string
	CharacterID     pgtype.UUID
	BackstoryItemID pgtype.UUID
}
