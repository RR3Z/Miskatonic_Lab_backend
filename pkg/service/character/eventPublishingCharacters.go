package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
)

func (s *EventPublishingCharacterService) GetAllCharacters(ctx context.Context, userID string) ([]characterDTO.CharacterShortModel, error) {
	characters, err := s.next.GetAllCharacters(ctx, userID)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharactersListFailed{
			UserID: userID,
			Err:    err,
		})

		return nil, err
	}

	s.publisher.Publish(ctx, characterEvents.CharactersListSucceeded{
		UserID: userID,
		Count:  len(characters),
	})

	return characters, nil
}

func (s *EventPublishingCharacterService) GetCharacter(ctx context.Context, input characterDTO.GetCharacterInput) (characterDTO.CharacterModel, error) {
	character, err := s.next.GetCharacter(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return characterDTO.CharacterModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		Name:        character.Name,
	})

	return character, nil
}

func (s *EventPublishingCharacterService) CreateCharacter(ctx context.Context, input characterDTO.CreateCharacterInput) (characterDTO.CharacterShortModel, error) {
	character, err := s.next.CreateCharacter(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCreateFailed{
			UserID: input.UserID,
			Err:    err,
		})
		return characterDTO.CharacterShortModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCreateSucceeded{
		UserID:      input.UserID,
		CharacterID: character.ID.String(),
		Name:        character.Name,
	})

	return character, nil
}

func (s *EventPublishingCharacterService) UpdateCharacter(ctx context.Context, input characterDTO.UpdateCharacterInput) (characterDTO.CharacterShortModel, error) {
	character, err := s.next.UpdateCharacter(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterUpdateFailed{
			UserID:      input.UserID,
			CharacterID: input.ID.String(),
			Err:         err,
		})
		return characterDTO.CharacterShortModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterUpdateSucceeded{
		UserID:      input.UserID,
		CharacterID: character.ID.String(),
		Name:        character.Name,
	})

	return character, nil
}

func (s *EventPublishingCharacterService) DeleteCharacter(ctx context.Context, input characterDTO.DeleteCharacterInput) error {
	if err := s.next.DeleteCharacter(ctx, input); err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.ID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.ID.String(),
	})

	return nil
}
