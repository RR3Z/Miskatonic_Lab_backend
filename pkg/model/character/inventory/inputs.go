package inventoryDTO

import "github.com/jackc/pgx/v5/pgtype"

type GetInventoryItemsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetInventoryItemInput struct {
	UserID      string
	CharacterID pgtype.UUID
	ItemID      pgtype.UUID
}

type CreateInventoryItemInput struct {
	Name        string
	Quantity    *int32
	Category    *string
	Description *string
	UserID      string
	CharacterID pgtype.UUID
}

type UpdateInventoryItemInput struct {
	Name        string
	Quantity    *int32
	Category    *string
	Description *string
	UserID      string
	CharacterID pgtype.UUID
	ItemID      pgtype.UUID
}

type DeleteInventoryItemInput struct {
	UserID      string
	CharacterID pgtype.UUID
	ItemID      pgtype.UUID
}
