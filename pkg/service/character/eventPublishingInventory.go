package character

import (
	"context"

	inventoryDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/inventory"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *EventPublishingCharacterService) GetInventoryItems(ctx context.Context, input inventoryDTO.GetInventoryItemsInput) ([]db.CharacterInventoryItem, error) {
	return s.next.GetInventoryItems(ctx, input)
}

func (s *EventPublishingCharacterService) GetInventoryItem(ctx context.Context, input inventoryDTO.GetInventoryItemInput) (db.CharacterInventoryItem, error) {
	return s.next.GetInventoryItem(ctx, input)
}

func (s *EventPublishingCharacterService) CreateInventoryItem(ctx context.Context, input inventoryDTO.CreateInventoryItemInput) (db.CharacterInventoryItem, error) {
	return s.next.CreateInventoryItem(ctx, input)
}

func (s *EventPublishingCharacterService) UpdateInventoryItem(ctx context.Context, input inventoryDTO.UpdateInventoryItemInput) (db.CharacterInventoryItem, error) {
	return s.next.UpdateInventoryItem(ctx, input)
}

func (s *EventPublishingCharacterService) DeleteInventoryItem(ctx context.Context, input inventoryDTO.DeleteInventoryItemInput) error {
	return s.next.DeleteInventoryItem(ctx, input)
}
