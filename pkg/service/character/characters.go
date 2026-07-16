package character

import (
	"context"
	"errors"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/publishers"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
)

const MaxCharactersPerUser int64 = 30

type CharacterService struct {
	repos         *repository.Repository
	publisher     events.EventPublisher
	portraitStore PortraitStore
}

func NewCharacterService(repos *repository.Repository, store PortraitStore, publisher events.EventPublisher) *CharacterService {
	if publisher == nil {
		publisher = &publishers.NoopPublisher{}
	}
	return &CharacterService{repos: repos, portraitStore: store, publisher: publisher}
}

// Characters
func (s *CharacterService) GetAllCharacters(ctx context.Context, userID string) ([]characterDTO.CharacterSummaryModel, error) {
	rows, err := s.repos.Queries.GetAllUserCharacterCards(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]characterDTO.CharacterSummaryModel, len(rows))
	for i, row := range rows {
		result[i] = characterDTO.ToCharacterSummaryModel(row)
		result[i].PortraitUrl = s.portraitURL(row.PortraitKey)
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

	result := characterDTO.ToCharacterModel(rawData)
	result.PortraitUrl = s.portraitURL(characterGeneralData.PortraitKey)
	return result, nil
}

func (s *CharacterService) CreateCharacter(ctx context.Context, input characterDTO.CreateCharacterInput) (characterDTO.CharacterShortModel, error) {
	if err := validateRequiredString(input.Name, 255, characterErrors.ErrNameRequired, characterErrors.ErrNameTooLong); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}
	if err := validateNonNegative(characterErrors.ErrAgeNegative, input.Age); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}
	if err := validateSex(input.Sex); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return characterDTO.CharacterShortModel{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	if _, err := queries.LockUserForCharacterCreation(ctx, input.UserID); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}

	count, err := queries.CountUserCharacters(ctx, input.UserID)
	if err != nil {
		return characterDTO.CharacterShortModel{}, err
	}
	if count >= MaxCharactersPerUser {
		return characterDTO.CharacterShortModel{}, characterErrors.ErrCharacterLimitReached
	}

	character, err := queries.CreateCharacter(ctx, db.CreateCharacterParams{
		UserID:     input.UserID,
		Name:       input.Name,
		Occupation: input.Occupation,
		Age:        input.Age,
		Sex:        input.Sex,
		Residence:  input.Residence,
		Birthplace: input.Birthplace,
	})
	if err != nil {
		return characterDTO.CharacterShortModel{}, characterErrors.MapCharacterConstraintError(err)
	}
	if err := createDefaultCharacterSkills(ctx, queries, input.UserID, character.ID); err != nil {
		return characterDTO.CharacterShortModel{}, characterErrors.MapCharacterConstraintError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}

	result := characterDTO.ToCharacterShortModel(character)
	result.PortraitUrl = s.portraitURL(character.PortraitKey)
	return result, nil
}

func (s *CharacterService) UpdateCharacter(ctx context.Context, input characterDTO.UpdateCharacterInput) (characterDTO.CharacterShortModel, error) {
	if err := validateRequiredString(input.Name, 255, characterErrors.ErrNameRequired, characterErrors.ErrNameTooLong); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}
	if err := validateNonNegative(characterErrors.ErrAgeNegative, input.Age); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}
	if err := validateSex(input.Sex); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}

	character, err := s.repos.Queries.UpdateCharacter(ctx, db.UpdateCharacterParams{
		UserID:     input.UserID,
		ID:         input.ID,
		Name:       input.Name,
		Occupation: input.Occupation,
		Age:        input.Age,
		Sex:        input.Sex,
		Residence:  input.Residence,
		Birthplace: input.Birthplace,
	})
	if err != nil {
		return characterDTO.CharacterShortModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	result := characterDTO.ToCharacterShortModel(character)
	result.PortraitUrl = s.portraitURL(character.PortraitKey)
	return result, nil
}

func (s *CharacterService) PatchCharacter(ctx context.Context, input characterDTO.PatchCharacterInput) (characterDTO.CharacterShortModel, error) {
	if !input.HasChanges() || input.Name.Set && input.Name.Value == nil {
		return characterDTO.CharacterShortModel{}, characterErrors.ErrPatchInvalid
	}
	if input.Name.Set {
		if err := validateRequiredString(*input.Name.Value, 255, characterErrors.ErrNameRequired, characterErrors.ErrNameTooLong); err != nil {
			return characterDTO.CharacterShortModel{}, err
		}
	}
	if input.Age.Set {
		if err := validateNonNegative(characterErrors.ErrAgeNegative, input.Age.Value); err != nil {
			return characterDTO.CharacterShortModel{}, err
		}
	}
	if input.Sex.Set {
		if err := validateSex(input.Sex.Value); err != nil {
			return characterDTO.CharacterShortModel{}, err
		}
	}

	name := ""
	if input.Name.Value != nil {
		name = *input.Name.Value
	}
	character, err := s.repos.Queries.PatchCharacter(ctx, db.PatchCharacterParams{
		UserID:        input.UserID,
		ID:            input.ID,
		SetName:       input.Name.Set,
		Name:          name,
		SetOccupation: input.Occupation.Set,
		Occupation:    input.Occupation.Value,
		SetAge:        input.Age.Set,
		Age:           input.Age.Value,
		SetSex:        input.Sex.Set,
		Sex:           input.Sex.Value,
		SetResidence:  input.Residence.Set,
		Residence:     input.Residence.Value,
		SetBirthplace: input.Birthplace.Set,
		Birthplace:    input.Birthplace.Value,
	})
	if err != nil {
		return characterDTO.CharacterShortModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	result := characterDTO.ToCharacterShortModel(character)
	result.PortraitUrl = s.portraitURL(character.PortraitKey)
	return result, nil
}

func (s *CharacterService) DeleteCharacter(ctx context.Context, input characterDTO.DeleteCharacterInput) error {
	character, err := s.repos.Queries.DeleteCharacter(ctx, db.DeleteCharacterParams{
		ID:     input.ID,
		UserID: input.UserID,
	})
	if err != nil {
		return err
	}

	if character.PortraitKey != nil {
		s.removePortraitFile(
			context.WithoutCancel(ctx),
			*character.PortraitKey,
			"failed to remove deleted character portrait",
			"character_id", input.ID.String(),
		)
	}
	return nil
}
