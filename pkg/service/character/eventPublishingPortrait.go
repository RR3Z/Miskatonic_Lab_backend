package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
)

func (s *EventPublishingCharacterService) ReplacePortrait(ctx context.Context, input characterDTO.ReplacePortraitInput) (characterDTO.CharacterShortModel, error) {
	character, err := s.next.ReplacePortrait(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterPortraitReplaceFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return characterDTO.CharacterShortModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterPortraitReplaceSucceeded{
		UserID:      input.UserID,
		CharacterID: character.ID.String(),
	})

	return character, nil
}
