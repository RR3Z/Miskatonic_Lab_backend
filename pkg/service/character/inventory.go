package character

import (
	"context"

	inventoryDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/inventory"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *CharacterService) GetInventoryItems(ctx context.Context, input inventoryDTO.GetInventoryItemsInput) ([]db.CharacterInventoryItem, error) {
	return s.repos.Queries.GetInventoryItems(ctx, db.GetInventoryItemsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
}

func (s *CharacterService) GetInventoryItem(ctx context.Context, input inventoryDTO.GetInventoryItemInput) (db.CharacterInventoryItem, error) {
	return s.repos.Queries.GetInventoryItem(ctx, db.GetInventoryItemParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		ItemID:      input.ItemID,
	})
}

func (s *CharacterService) CreateInventoryItem(ctx context.Context, input inventoryDTO.CreateInventoryItemInput) (db.CharacterInventoryItem, error) {
	name, quantity, category, description, err := normalizeInventoryItemInput(input.Name, input.Quantity, input.Category, input.Description)
	if err != nil {
		return db.CharacterInventoryItem{}, err
	}

	return s.repos.Queries.CreateInventoryItem(ctx, db.CreateInventoryItemParams{
		Name:        name,
		Quantity:    quantity,
		Category:    category,
		Description: description,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
}

func (s *CharacterService) UpdateInventoryItem(ctx context.Context, input inventoryDTO.UpdateInventoryItemInput) (db.CharacterInventoryItem, error) {
	name, quantity, category, description, err := normalizeInventoryItemInput(input.Name, input.Quantity, input.Category, input.Description)
	if err != nil {
		return db.CharacterInventoryItem{}, err
	}

	return s.repos.Queries.UpdateInventoryItem(ctx, db.UpdateInventoryItemParams{
		Name:        name,
		Quantity:    quantity,
		Category:    category,
		Description: description,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		ItemID:      input.ItemID,
	})
}

func (s *CharacterService) DeleteInventoryItem(ctx context.Context, input inventoryDTO.DeleteInventoryItemInput) error {
	_, err := s.repos.Queries.DeleteInventoryItem(ctx, db.DeleteInventoryItemParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		ItemID:      input.ItemID,
	})
	return err
}
