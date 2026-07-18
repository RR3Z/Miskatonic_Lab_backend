package inventoryDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type InventoryItemModel struct {
	ID          pgtype.UUID        `json:"id"`
	CharacterID pgtype.UUID        `json:"character_id"`
	Name        string             `json:"name"`
	Quantity    *int32             `json:"quantity"`
	Category    *string            `json:"category"`
	Description *string            `json:"description"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

func ToInventoryItemModel(item db.CharacterInventoryItem) InventoryItemModel {
	return InventoryItemModel{
		ID:          item.ID,
		CharacterID: item.CharacterID,
		Name:        item.Name,
		Quantity:    item.Quantity,
		Category:    item.Category,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func ToInventoryItemModels(items []db.CharacterInventoryItem) []InventoryItemModel {
	models := make([]InventoryItemModel, len(items))
	for i, item := range items {
		models[i] = ToInventoryItemModel(item)
	}
	return models
}
