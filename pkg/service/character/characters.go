package character

import (
	"context"
	"errors"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
)

type CharacterService struct {
	repos     *repository.Repository
	publisher events.EventPublisher
}

func NewCharacterService(repos *repository.Repository, publisher ...events.EventPublisher) *CharacterService {
	service := &CharacterService{repos: repos}
	if len(publisher) > 0 {
		service.publisher = publisher[0]
	}

	return service
}

// Characters
func (s *CharacterService) GetAllCharacters(ctx context.Context, userID string) ([]characterDTO.CharacterShortModel, error) {
	dbCharacters, err := s.repos.Queries.GetAllUserCharacters(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]characterDTO.CharacterShortModel, len(dbCharacters))
	for i, c := range dbCharacters {
		result[i] = characterDTO.ToCharacterShortModel(c)
	}

	return result, nil
}

func (s *CharacterService) GetCharacter(ctx context.Context, input characterDTO.GetCharacterInput) (characterDTO.CharacterModel, error) {
	var rawData characterDTO.CharacterDBData

	characterGeneralData, err := s.repos.Queries.GetCharacter(ctx, db.GetCharacterParams{
		UserID: input.UserID,
		ID:     input.CharacterID,
	})
	if err != nil {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Character = characterGeneralData

	characteristics, err := s.repos.Queries.GetCharacteristics(ctx, db.GetCharacteristicsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Characteristics = characteristics

	derivedStats, err := s.repos.Queries.GetDerivedStats(ctx, db.GetDerivedStatsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.DerivedStats = derivedStats

	hp, err := s.repos.Queries.GetHealthState(ctx, db.GetHealthStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.HP = hp

	sanity, err := s.repos.Queries.GetSanityState(ctx, db.GetSanityStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Sanity = sanity

	mp, err := s.repos.Queries.GetMagicState(ctx, db.GetMagicStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.MP = mp

	luck, err := s.repos.Queries.GetLuckState(ctx, db.GetLuckStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Luck = luck

	skills, err := s.repos.Queries.GetSkills(ctx, characterGeneralData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Skills = skills

	notes, err := s.repos.Queries.GetNotes(ctx, db.GetNotesParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Notes = notes

	backstory, err := s.repos.Queries.GetBackstory(ctx, characterGeneralData.ID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return characterDTO.CharacterModel{}, err
		}
	} else {
		rawData.Backstory = &backstory
		rawData.BackstoryItems, err = s.repos.Queries.GetBackstoryItemsByBackstoryID(ctx, backstory.ID)
		if err != nil {
			return characterDTO.CharacterModel{}, err
		}
	}

	finances, err := s.repos.Queries.GetFinances(ctx, db.GetFinancesParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return characterDTO.CharacterModel{}, err
		}
	} else {
		rawData.Finances = &finances
	}

	return characterDTO.ToCharacterModel(rawData), nil
}

func (s *CharacterService) CreateCharacter(ctx context.Context, input characterDTO.CreateCharacterInput) (characterDTO.CharacterShortModel, error) {
	if err := validateRequiredString(input.Name, 255, characterErrors.ErrNameRequired, characterErrors.ErrNameTooLong); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}

	character, err := s.repos.Queries.CreateCharacter(ctx, db.CreateCharacterParams{
		UserID:     input.UserID,
		Name:       input.Name,
		PlayerName: input.PlayerName,
		Occupation: input.Occupation,
		Age:        input.Age,
		Sex:        input.Sex,
		Residence:  input.Residence,
		Birthplace: input.Birthplace,
	})
	if err != nil {
		return characterDTO.CharacterShortModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	return characterDTO.ToCharacterShortModel(character), nil
}

func (s *CharacterService) UpdateCharacter(ctx context.Context, input characterDTO.UpdateCharacterInput) (characterDTO.CharacterShortModel, error) {
	if err := validateRequiredString(input.Name, 255, characterErrors.ErrNameRequired, characterErrors.ErrNameTooLong); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}

	character, err := s.repos.Queries.UpdateCharacter(ctx, db.UpdateCharacterParams{
		UserID:     input.UserID,
		ID:         input.ID,
		Name:       input.Name,
		PlayerName: input.PlayerName,
		Occupation: input.Occupation,
		Age:        input.Age,
		Sex:        input.Sex,
		Residence:  input.Residence,
		Birthplace: input.Birthplace,
	})
	if err != nil {
		return characterDTO.CharacterShortModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	characteristics, shouldRecalculate := s.getCharacteristicsForDerivedStatsRecalculation(ctx, input.UserID, input.ID)
	if shouldRecalculate {
		s.recalculateDerivedStats(ctx, character.UserID, character.ID, character.Age, characteristics, "character_update")
	}

	return characterDTO.ToCharacterShortModel(character), nil
}

func (s *CharacterService) DeleteCharacter(ctx context.Context, input characterDTO.DeleteCharacterInput) error {
	_, err := s.repos.Queries.DeleteCharacter(ctx, db.DeleteCharacterParams{
		ID:     input.ID,
		UserID: input.UserID,
	})
	return err
}
