package room

import (
	"context"
	"errors"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *RoomService) ListSelectedCharacters(ctx context.Context, input model.ListSelectedCharactersInput) ([]model.SelectedCharacterModel, error) {
	if err := s.EnsureMember(ctx, input.RoomID, input.UserID); err != nil {
		return nil, err
	}

	members, err := s.repos.Queries.ListVisibleSelectedRoomMembers(ctx, db.ListVisibleSelectedRoomMembersParams{
		RoomID: input.RoomID,
		UserID: input.UserID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.SelectedCharacterModel, 0, len(members))
	for _, member := range members {
		if !member.CharacterID.Valid {
			continue
		}

		character, err := s.loadSelectedCharacter(ctx, member.UserID, member.CharacterID)
		if err != nil {
			return nil, err
		}

		result = append(result, model.ToSelectedCharacterModel(member, character))
	}

	return result, nil
}

func (s *RoomService) loadSelectedCharacter(ctx context.Context, userID string, characterID pgtype.UUID) (characterDTO.CharacterModel, error) {
	var rawData characterDTO.CharacterDBData

	character, err := s.repos.Queries.GetCharacter(ctx, db.GetCharacterParams{
		UserID: userID,
		ID:     characterID,
	})
	if err != nil {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Character = character

	characteristics, err := s.repos.Queries.GetCharacteristics(ctx, db.GetCharacteristicsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Characteristics = characteristics

	derivedStats, err := s.repos.Queries.GetDerivedStats(ctx, db.GetDerivedStatsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.DerivedStats = derivedStats

	hp, err := s.repos.Queries.GetHealthState(ctx, db.GetHealthStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.HP = hp

	sanity, err := s.repos.Queries.GetSanityState(ctx, db.GetSanityStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Sanity = sanity

	mp, err := s.repos.Queries.GetMagicState(ctx, db.GetMagicStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.MP = mp

	luck, err := s.repos.Queries.GetLuckState(ctx, db.GetLuckStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Luck = luck

	skills, err := s.repos.Queries.GetSkills(ctx, character.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Skills = skills

	notes, err := s.repos.Queries.GetNotes(ctx, db.GetNotesParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Notes = notes

	backstory, err := s.repos.Queries.GetBackstory(ctx, character.ID)
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
		UserID:      userID,
		CharacterID: characterID,
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
