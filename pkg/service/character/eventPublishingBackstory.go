package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
)

func (s *EventPublishingCharacterService) GetBackstory(ctx context.Context, input backstoriesDTO.GetBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	backstory, err := s.next.GetBackstory(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return backstoriesDTO.BackstoryModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return backstory, nil
}

func (s *EventPublishingCharacterService) UpsertBackstory(ctx context.Context, input backstoriesDTO.UpsertBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	backstory, err := s.next.UpsertBackstory(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return backstoriesDTO.BackstoryModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return backstory, nil
}

func (s *EventPublishingCharacterService) DeleteBackstory(ctx context.Context, input backstoriesDTO.DeleteBackstoryInput) error {
	err := s.next.DeleteBackstory(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

func (s *EventPublishingCharacterService) GetBackstoryItems(ctx context.Context, input backstoriesDTO.GetBackstoryItemsInput) ([]backstoriesDTO.BackstoryItemModel, error) {
	items, err := s.next.GetBackstoryItems(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemsListFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return nil, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemsListSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		Count:       len(items),
	})

	return items, nil
}

func (s *EventPublishingCharacterService) GetBackstoryItem(ctx context.Context, input backstoriesDTO.GetBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	item, err := s.next.GetBackstoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemGetFailed{
			UserID:          input.UserID,
			CharacterID:     input.CharacterID.String(),
			BackstoryItemID: input.BackstoryItemID.String(),
			Err:             err,
		})
		return backstoriesDTO.BackstoryItemModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemGetSucceeded{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID.String(),
		BackstoryItemID: input.BackstoryItemID.String(),
		Section:         item.Section,
		Title:           item.Title,
	})

	return item, nil
}

func (s *EventPublishingCharacterService) CreateBackstoryItem(ctx context.Context, input backstoriesDTO.CreateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	item, err := s.next.CreateBackstoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemCreateFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return backstoriesDTO.BackstoryItemModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemCreateSucceeded{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID.String(),
		BackstoryItemID: item.ID.String(),
		Section:         item.Section,
		Title:           item.Title,
	})

	return item, nil
}

func (s *EventPublishingCharacterService) UpdateBackstoryItem(ctx context.Context, input backstoriesDTO.UpdateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	item, err := s.next.UpdateBackstoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemUpdateFailed{
			UserID:          input.UserID,
			CharacterID:     input.CharacterID.String(),
			BackstoryItemID: input.BackstoryItemID.String(),
			Err:             err,
		})
		return backstoriesDTO.BackstoryItemModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemUpdateSucceeded{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID.String(),
		BackstoryItemID: input.BackstoryItemID.String(),
		Section:         item.Section,
		Title:           item.Title,
	})

	return item, nil
}

func (s *EventPublishingCharacterService) DeleteBackstoryItem(ctx context.Context, input backstoriesDTO.DeleteBackstoryItemInput) error {
	err := s.next.DeleteBackstoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemDeleteFailed{
			UserID:          input.UserID,
			CharacterID:     input.CharacterID.String(),
			BackstoryItemID: input.BackstoryItemID.String(),
			Err:             err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemDeleteSucceeded{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID.String(),
		BackstoryItemID: input.BackstoryItemID.String(),
	})

	return nil
}
