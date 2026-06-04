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
	DeleteCharacter(ctx context.Context, input db.DeleteCharacterParams) error

	GetCharacteristics(ctx context.Context, input db.GetCharacteristicsParams) (db.Characteristic, error)
	UpsertCharacteristics(ctx context.Context, input db.UpsertCharacteristicsParams) (db.Characteristic, error)
	DeleteCharacteristics(ctx context.Context, input db.DeleteCharacteristicsParams) error

	GetNotes(ctx context.Context, input db.GetNotesParams) ([]db.Note, error)
	GetNote(ctx context.Context, input db.GetNoteParams) (db.Note, error)
	CreateNote(ctx context.Context, input db.CreateNoteParams) (db.Note, error)
	UpdateNote(ctx context.Context, input db.UpdateNoteParams) (db.Note, error)
	DeleteNote(ctx context.Context, input db.DeleteNoteParams) error
}

type CharacterService struct {
	repos *repository.Repository
}

func NewCharacterService(repos *repository.Repository) *CharacterService {
	return &CharacterService{repos: repos}
}

// Characters
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

	characteristics, err := s.repos.Queries.GetCharacteristics(ctx, db.GetCharacteristicsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
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

	notes, err := s.repos.Queries.GetNotes(ctx, db.GetNotesParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
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
		rawData.BackstoryItems, err = s.repos.Queries.GetBackstoryItemsByBackstoryID(ctx, backstory.ID)
		if err != nil {
			return model.CharacterModel{}, err
		}
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

func (s *CharacterService) DeleteCharacter(ctx context.Context, input db.DeleteCharacterParams) error {
	_, err := s.repos.Queries.DeleteCharacter(ctx, db.DeleteCharacterParams{
		ID:     input.ID,
		UserID: input.UserID,
	})
	return err
}

// Characteristics
func (s *CharacterService) GetCharacteristics(ctx context.Context, input db.GetCharacteristicsParams) (db.Characteristic, error) {
	characteristics, err := s.repos.Queries.GetCharacteristics(ctx, input)
	if err != nil {
		return db.Characteristic{}, err
	}

	return characteristics, nil
}

func (s *CharacterService) UpsertCharacteristics(ctx context.Context, input db.UpsertCharacteristicsParams) (db.Characteristic, error) {
	characteristics, err := s.repos.Queries.UpsertCharacteristics(ctx, input)
	if err != nil {
		return db.Characteristic{}, err
	}

	return characteristics, nil
}

func (s *CharacterService) DeleteCharacteristics(ctx context.Context, input db.DeleteCharacteristicsParams) error {
	if _, err := s.repos.Queries.DeleteCharacteristics(ctx, input); err != nil {
		return err
	}

	return nil
}

// Notes
func (s *CharacterService) GetNotes(ctx context.Context, input db.GetNotesParams) ([]db.Note, error) {
	notes, err := s.repos.Queries.GetNotes(ctx, input)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *CharacterService) GetNote(ctx context.Context, input db.GetNoteParams) (db.Note, error) {
	note, err := s.repos.Queries.GetNote(ctx, input)
	if err != nil {
		return db.Note{}, err
	}

	return note, nil
}

func (s *CharacterService) CreateNote(ctx context.Context, input db.CreateNoteParams) (db.Note, error) {
	note, err := s.repos.Queries.CreateNote(ctx, input)
	if err != nil {
		return db.Note{}, err
	}

	return note, nil
}

func (s *CharacterService) UpdateNote(ctx context.Context, input db.UpdateNoteParams) (db.Note, error) {
	note, err := s.repos.Queries.UpdateNote(ctx, input)
	if err != nil {
		return db.Note{}, err
	}

	return note, nil
}

func (s *CharacterService) DeleteNote(ctx context.Context, input db.DeleteNoteParams) error {
	_, err := s.repos.Queries.DeleteNote(ctx, input)
	return err
}
