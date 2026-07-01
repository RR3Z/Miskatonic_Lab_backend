package character

import (
	"context"

	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
)

// Backstory
func (s *CharacterService) GetBackstory(ctx context.Context, input backstoriesDTO.GetBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	backstory, err := s.repos.Queries.GetBackstoryByCharacter(ctx, db.GetBackstoryByCharacterParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryModel{}, err
	}

	items, err := s.repos.Queries.GetBackstoryItems(ctx, db.GetBackstoryItemsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryModel{}, err
	}

	return backstoriesDTO.ToBackstoryModel(backstory, items), nil
}

func (s *CharacterService) UpsertBackstory(ctx context.Context, input backstoriesDTO.UpsertBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	backstory, err := s.repos.Queries.UpsertBackstory(ctx, db.UpsertBackstoryParams{
		UserID:              input.UserID,
		CharacterID:         input.CharacterID,
		PersonalDescription: input.PersonalDescription,
	})
	if err != nil {
		return backstoriesDTO.BackstoryModel{}, err
	}

	items, err := s.repos.Queries.GetBackstoryItems(ctx, db.GetBackstoryItemsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryModel{}, err
	}

	return backstoriesDTO.ToBackstoryModel(backstory, items), nil
}

func (s *CharacterService) DeleteBackstory(ctx context.Context, input backstoriesDTO.DeleteBackstoryInput) error {
	if _, err := s.repos.Queries.DeleteBackstory(ctx, db.DeleteBackstoryParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *CharacterService) GetBackstoryItems(ctx context.Context, input backstoriesDTO.GetBackstoryItemsInput) ([]backstoriesDTO.BackstoryItemModel, error) {
	items, err := s.repos.Queries.GetBackstoryItems(ctx, db.GetBackstoryItemsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return nil, err
	}

	return backstoriesDTO.ToBackstoryItemModels(items), nil
}

func (s *CharacterService) GetBackstoryItem(ctx context.Context, input backstoriesDTO.GetBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	item, err := s.repos.Queries.GetBackstoryItem(ctx, db.GetBackstoryItemParams{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID,
		BackstoryItemID: input.BackstoryItemID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryItemModel{}, err
	}

	return backstoriesDTO.ToBackstoryItemModel(item), nil
}

func (s *CharacterService) CreateBackstoryItem(ctx context.Context, input backstoriesDTO.CreateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	if err := validateBackstoryItemInput(input.Section, input.Title, input.Text); err != nil {
		return backstoriesDTO.BackstoryItemModel{}, err
	}

	item, err := s.repos.Queries.CreateBackstoryItem(ctx, db.CreateBackstoryItemParams{
		Section:     input.Section,
		Title:       input.Title,
		Text:        input.Text,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryItemModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	return backstoriesDTO.ToBackstoryItemModel(item), nil
}

func (s *CharacterService) UpdateBackstoryItem(ctx context.Context, input backstoriesDTO.UpdateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	if err := validateBackstoryItemInput(input.Section, input.Title, input.Text); err != nil {
		return backstoriesDTO.BackstoryItemModel{}, err
	}

	item, err := s.repos.Queries.UpdateBackstoryItem(ctx, db.UpdateBackstoryItemParams{
		Section:         input.Section,
		Title:           input.Title,
		Text:            input.Text,
		UserID:          input.UserID,
		CharacterID:     input.CharacterID,
		BackstoryItemID: input.BackstoryItemID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryItemModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	return backstoriesDTO.ToBackstoryItemModel(item), nil
}

func (s *CharacterService) DeleteBackstoryItem(ctx context.Context, input backstoriesDTO.DeleteBackstoryItemInput) error {
	_, err := s.repos.Queries.DeleteBackstoryItem(ctx, db.DeleteBackstoryItemParams{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID,
		BackstoryItemID: input.BackstoryItemID,
	})
	return err
}
