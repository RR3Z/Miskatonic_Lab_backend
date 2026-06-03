package service

import (
	"context"
	"errors"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
)

type ICharacter interface {
	GetAllCharacters(ctx context.Context, userID string) ([]model.CharacterModel, error)
	GetCharacter(ctx context.Context, input model.GetCharacterInput) (model.CharacterModel, error)
	CreateCharacter(ctx context.Context, input db.CreateCharacterParams) (model.CharacterModel, error)
	UpdateCharacter(ctx context.Context, input db.UpdateCharacterParams) (model.CharacterModel, error)
	DeleteCharacter(ctx context.Context, input model.DeleteCharacterInput) error
}

type CharacterService struct {
	repos *repository.Repository
}

func NewCharacterService(repos *repository.Repository) *CharacterService {
	return &CharacterService{repos: repos}
}

func (s *CharacterService) GetAllCharacters(ctx context.Context, userID string) ([]model.CharacterModel, error) {
	characters, err := s.repos.Queries.GetAllUserCharacters(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]model.CharacterModel, len(characters))
	for i, c := range characters {
		result[i] = model.ToShortCharacterModel(c)
	}

	return result, nil
}

func (s *CharacterService) GetCharacter(ctx context.Context, input model.GetCharacterInput) (model.CharacterModel, error) {
	var rawData model.CharacterDBData

	characterGeneralData, err := s.repos.Queries.GetCharacter(ctx, db.GetCharacterParams{
		UserID: input.UserID,
		ID:     input.CharacterID,
	})
	if err != nil {
		return model.CharacterModel{}, err
	}
	rawData.Character = characterGeneralData

	characteristics, err := s.repos.Queries.GetCharacteristics(ctx, characterGeneralData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.Characteristics = characteristics

	derivedStats, err := s.repos.Queries.GetDerivedStats(ctx, characterGeneralData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.DerivedStats = derivedStats

	hp, err := s.repos.Queries.GetHealthState(ctx, characterGeneralData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.HP = hp

	sanity, err := s.repos.Queries.GetSanityState(ctx, characterGeneralData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.Sanity = sanity

	mp, err := s.repos.Queries.GetMagicState(ctx, characterGeneralData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.MP = mp

	luck, err := s.repos.Queries.GetLuckState(ctx, characterGeneralData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.Luck = luck

	skills, err := s.repos.Queries.GetSkills(ctx, characterGeneralData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.Skills = skills

	notes, err := s.repos.Queries.GetNotes(ctx, characterGeneralData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.Notes = notes

	backstory, err := s.repos.Queries.GetBackstory(ctx, characterGeneralData.ID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return model.CharacterModel{}, err
		}
	} else {
		rawData.Backstory = &backstory
		rawData.BackstoryItems, _ = s.repos.Queries.GetBackstoryItemsByBackstoryID(ctx, backstory.ID)
	}

	finances, err := s.repos.Queries.GetFinances(ctx, characterGeneralData.ID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return model.CharacterModel{}, err
		}
	} else {
		rawData.Finances = &finances
	}

	return model.ToFullCharacterModel(rawData), nil
}

func (s *CharacterService) CreateCharacter(ctx context.Context, input db.CreateCharacterParams) (model.CharacterModel, error) {
	character, err := s.repos.Queries.CreateCharacter(ctx, input)
	if err != nil {
		return model.CharacterModel{}, err
	}

	return model.ToShortCharacterModel(character), nil
}

func (s *CharacterService) UpdateCharacter(ctx context.Context, input db.UpdateCharacterParams) (model.CharacterModel, error) {
	character, err := s.repos.Queries.UpdateCharacter(ctx, input)
	if err != nil {
		return model.CharacterModel{}, err
	}

	return model.ToShortCharacterModel(character), nil
}

func (s *CharacterService) DeleteCharacter(ctx context.Context, input model.DeleteCharacterInput) error {
	return s.repos.Queries.DeleteCharacter(ctx, db.DeleteCharacterParams{
		ID:     input.CharacterID,
		UserID: input.UserID,
	})
}
