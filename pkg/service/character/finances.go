package character

import (
	"context"

	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
)

// Finances
func (s *CharacterService) GetFinances(ctx context.Context, input financesDTO.GetFinancesInput) (db.Finance, error) {
	finances, err := s.repos.Queries.GetFinances(ctx, db.GetFinancesParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.Finance{}, err
	}

	return finances, nil
}

func (s *CharacterService) UpsertFinances(ctx context.Context, input financesDTO.UpsertFinancesInput) (db.Finance, error) {
	if input.SpendingLimit != nil {
		if err := validateRequiredString(*input.SpendingLimit, 120, nil, characterErrors.ErrFinancesMoneyTooLong); err != nil {
			return db.Finance{}, err
		}
	}
	if input.Cash != nil {
		if err := validateRequiredString(*input.Cash, 120, nil, characterErrors.ErrFinancesMoneyTooLong); err != nil {
			return db.Finance{}, err
		}
	}

	finances, err := s.repos.Queries.UpsertFinances(ctx, db.UpsertFinancesParams{
		UserID:        input.UserID,
		CharacterID:   input.CharacterID,
		SpendingLimit: input.SpendingLimit,
		Cash:          input.Cash,
		Assets:        input.Assets,
	})
	if err != nil {
		return db.Finance{}, characterErrors.MapCharacterConstraintError(err)
	}

	return finances, nil
}

func (s *CharacterService) DeleteFinances(ctx context.Context, input financesDTO.DeleteFinancesInput) error {
	if _, err := s.repos.Queries.DeleteFinances(ctx, db.DeleteFinancesParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}
