package character

import (
	"context"

	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
)

// Characteristics
func (s *CharacterService) GetCharacteristics(ctx context.Context, input characteristicsDTO.GetCharacteristicsInput) (db.Characteristic, error) {
	characteristics, err := s.repos.Queries.GetCharacteristics(ctx, db.GetCharacteristicsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.Characteristic{}, err
	}

	return characteristics, nil
}

func (s *CharacterService) UpsertCharacteristics(ctx context.Context, input characteristicsDTO.UpsertCharacteristicsInput) (db.Characteristic, error) {
	if err := validateNonNegative(characterErrors.ErrCharacteristicsNegative, input.Strength, input.Constitution, input.Size, input.Dexterity, input.Appearance, input.Intelligence, input.Power, input.Education); err != nil {
		return db.Characteristic{}, err
	}

	characteristics, err := s.repos.Queries.UpsertCharacteristics(ctx, db.UpsertCharacteristicsParams{
		Strength:     input.Strength,
		Constitution: input.Constitution,
		Size:         input.Size,
		Dexterity:    input.Dexterity,
		Appearance:   input.Appearance,
		Intelligence: input.Intelligence,
		Power:        input.Power,
		Education:    input.Education,
		UserID:       input.UserID,
		CharacterID:  input.CharacterID,
	})
	if err != nil {
		return db.Characteristic{}, err
	}

	character, err := s.repos.Queries.GetCharacter(ctx, db.GetCharacterParams{
		UserID: input.UserID,
		ID:     input.CharacterID,
	})
	if err == nil {
		s.recalculateDerivedStats(ctx, input.UserID, input.CharacterID, character.Age, characteristics, "characteristics_upsert")
	}

	return characteristics, nil
}

func (s *CharacterService) DeleteCharacteristics(ctx context.Context, input characteristicsDTO.DeleteCharacteristicsInput) error {
	if _, err := s.repos.Queries.DeleteCharacteristics(ctx, db.DeleteCharacteristicsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}
