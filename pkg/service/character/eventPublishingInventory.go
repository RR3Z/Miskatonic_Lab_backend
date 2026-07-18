package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	inventoryDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/inventory"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *EventPublishingCharacterService) GetInventoryItems(ctx context.Context, input inventoryDTO.GetInventoryItemsInput) ([]db.CharacterInventoryItem, error) {
	items, err := s.next.GetInventoryItems(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterInventoryItemsListFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return nil, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterInventoryItemsListSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		Count:       len(items),
	})
	return items, nil
}

func (s *EventPublishingCharacterService) GetInventoryItem(ctx context.Context, input inventoryDTO.GetInventoryItemInput) (db.CharacterInventoryItem, error) {
	item, err := s.next.GetInventoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterInventoryItemGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			InventoryID: input.ItemID.String(),
			Err:         err,
		})
		return db.CharacterInventoryItem{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterInventoryItemGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		InventoryID: item.ID.String(),
		Name:        item.Name,
	})
	return item, nil
}

func (s *EventPublishingCharacterService) CreateInventoryItem(ctx context.Context, input inventoryDTO.CreateInventoryItemInput) (db.CharacterInventoryItem, error) {
	item, err := s.next.CreateInventoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterInventoryItemCreateFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.CharacterInventoryItem{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterInventoryItemCreateSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		InventoryID: item.ID.String(),
		Name:        item.Name,
	})
	return item, nil
}

func (s *EventPublishingCharacterService) UpdateInventoryItem(ctx context.Context, input inventoryDTO.UpdateInventoryItemInput) (db.CharacterInventoryItem, error) {
	item, err := s.next.UpdateInventoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterInventoryItemUpdateFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			InventoryID: input.ItemID.String(),
			Err:         err,
		})
		return db.CharacterInventoryItem{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterInventoryItemUpdateSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		InventoryID: item.ID.String(),
		Name:        item.Name,
	})
	return item, nil
}

func (s *EventPublishingCharacterService) DeleteInventoryItem(ctx context.Context, input inventoryDTO.DeleteInventoryItemInput) error {
	err := s.next.DeleteInventoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterInventoryItemDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			InventoryID: input.ItemID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterInventoryItemDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		InventoryID: input.ItemID.String(),
	})
	return nil
}
