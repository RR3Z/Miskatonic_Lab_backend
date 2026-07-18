package roomHelpers

import (
	"context"
	"errors"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func LoadCharacterSheet(ctx context.Context, queries *db.Queries, userID string, characterID pgtype.UUID) (characterDTO.CharacterModel, error) {
	var rawData characterDTO.CharacterDBData

	character, err := queries.GetCharacter(ctx, db.GetCharacterParams{
		UserID: userID,
		ID:     characterID,
	})
	if err != nil {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Character = character

	characteristics, err := queries.GetCharacteristics(ctx, db.GetCharacteristicsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Characteristics = characteristics

	derivedStats, err := queries.GetDerivedStats(ctx, db.GetDerivedStatsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.DerivedStats = derivedStats

	hp, err := queries.GetHealthState(ctx, db.GetHealthStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.HP = hp

	sanity, err := queries.GetSanityState(ctx, db.GetSanityStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Sanity = sanity

	mp, err := queries.GetMagicState(ctx, db.GetMagicStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.MP = mp

	luck, err := queries.GetLuckState(ctx, db.GetLuckStateParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Luck = luck

	skills, err := queries.GetSkills(ctx, character.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Skills = skills

	notes, err := queries.GetNotes(ctx, db.GetNotesParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.Notes = notes

	inventoryItems, err := queries.GetInventoryItems(ctx, db.GetInventoryItemsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return characterDTO.CharacterModel{}, err
	}
	rawData.InventoryItems = inventoryItems

	backstory, err := queries.GetBackstory(ctx, character.ID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return characterDTO.CharacterModel{}, err
		}
	} else {
		rawData.Backstory = &backstory
		rawData.BackstoryItems, err = queries.GetBackstoryItemsByBackstoryID(ctx, backstory.ID)
		if err != nil {
			return characterDTO.CharacterModel{}, err
		}
	}

	finances, err := queries.GetFinances(ctx, db.GetFinancesParams{
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
