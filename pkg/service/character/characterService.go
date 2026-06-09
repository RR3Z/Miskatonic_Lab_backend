package character

import (
	"context"
	"errors"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/calculators"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

	derivedStats, err := s.repos.Queries.GetDerivedStats(ctx, db.GetDerivedStatsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.DerivedStats = derivedStats

	hp, err := s.repos.Queries.GetHealthState(ctx, db.GetHealthStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.HP = hp

	sanity, err := s.repos.Queries.GetSanityState(ctx, db.GetSanityStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.Sanity = sanity

	mp, err := s.repos.Queries.GetMagicState(ctx, db.GetMagicStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.CharacterModel{}, err
	}
	rawData.MP = mp

	luck, err := s.repos.Queries.GetLuckState(ctx, db.GetLuckStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
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

	finances, err := s.repos.Queries.GetFinances(ctx, db.GetFinancesParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
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

	characteristics, shouldRecalculate := s.getCharacteristicsForDerivedStatsRecalculation(ctx, input.UserID, input.ID)
	if shouldRecalculate {
		s.recalculateDerivedStats(ctx, character.UserID, character.ID, character.Age, characteristics, "character_update")
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

// Health
func (s *CharacterService) GetHealth(ctx context.Context, input db.GetHealthStateParams) (db.HealthState, error) {
	health, err := s.repos.Queries.GetHealthState(ctx, input)
	if err != nil {
		return db.HealthState{}, err
	}

	return health, nil
}

func (s *CharacterService) UpsertHealth(ctx context.Context, input db.UpsertHealthStateParams) (db.HealthState, error) {
	if err := s.validateHealthState(ctx, input); err != nil {
		return db.HealthState{}, err
	}

	health, err := s.repos.Queries.UpsertHealthState(ctx, input)
	if err != nil {
		return db.HealthState{}, err
	}

	return health, nil
}

func (s *CharacterService) DeleteHealth(ctx context.Context, input db.DeleteHealthStateParams) error {
	if _, err := s.repos.Queries.DeleteHealthState(ctx, input); err != nil {
		return err
	}

	return nil
}

// Sanity
func (s *CharacterService) GetSanity(ctx context.Context, input db.GetSanityStateParams) (db.SanityState, error) {
	sanity, err := s.repos.Queries.GetSanityState(ctx, input)
	if err != nil {
		return db.SanityState{}, err
	}

	return sanity, nil
}

func (s *CharacterService) UpsertSanity(ctx context.Context, input db.UpsertSanityStateParams) (db.SanityState, error) {
	if err := s.validateSanityState(ctx, input); err != nil {
		return db.SanityState{}, err
	}

	sanity, err := s.repos.Queries.UpsertSanityState(ctx, input)
	if err != nil {
		return db.SanityState{}, err
	}

	return sanity, nil
}

func (s *CharacterService) DeleteSanity(ctx context.Context, input db.DeleteSanityStateParams) error {
	if _, err := s.repos.Queries.DeleteSanityState(ctx, input); err != nil {
		return err
	}

	return nil
}

// Magic
func (s *CharacterService) GetMagic(ctx context.Context, input db.GetMagicStateParams) (db.MagicState, error) {
	magic, err := s.repos.Queries.GetMagicState(ctx, input)
	if err != nil {
		return db.MagicState{}, err
	}

	return magic, nil
}

func (s *CharacterService) UpsertMagic(ctx context.Context, input db.UpsertMagicStateParams) (db.MagicState, error) {
	if err := s.validateMagicState(ctx, input); err != nil {
		return db.MagicState{}, err
	}

	magic, err := s.repos.Queries.UpsertMagicState(ctx, input)
	if err != nil {
		return db.MagicState{}, err
	}

	return magic, nil
}

func (s *CharacterService) DeleteMagic(ctx context.Context, input db.DeleteMagicStateParams) error {
	if _, err := s.repos.Queries.DeleteMagicState(ctx, input); err != nil {
		return err
	}

	return nil
}

// Luck
func (s *CharacterService) GetLuck(ctx context.Context, input db.GetLuckStateParams) (db.LuckState, error) {
	luck, err := s.repos.Queries.GetLuckState(ctx, input)
	if err != nil {
		return db.LuckState{}, err
	}

	return luck, nil
}

func (s *CharacterService) UpsertLuck(ctx context.Context, input db.UpsertLuckStateParams) (db.LuckState, error) {
	if err := s.validateLuckState(ctx, input); err != nil {
		return db.LuckState{}, err
	}

	luck, err := s.repos.Queries.UpsertLuckState(ctx, input)
	if err != nil {
		return db.LuckState{}, err
	}

	return luck, nil
}

func (s *CharacterService) DeleteLuck(ctx context.Context, input db.DeleteLuckStateParams) error {
	if _, err := s.repos.Queries.DeleteLuckState(ctx, input); err != nil {
		return err
	}

	return nil
}

// Finances
func (s *CharacterService) GetFinances(ctx context.Context, input db.GetFinancesParams) (db.Finance, error) {
	finances, err := s.repos.Queries.GetFinances(ctx, input)
	if err != nil {
		return db.Finance{}, err
	}

	return finances, nil
}

func (s *CharacterService) UpsertFinances(ctx context.Context, input db.UpsertFinancesParams) (db.Finance, error) {
	finances, err := s.repos.Queries.UpsertFinances(ctx, input)
	if err != nil {
		return db.Finance{}, err
	}

	return finances, nil
}

func (s *CharacterService) DeleteFinances(ctx context.Context, input db.DeleteFinancesParams) error {
	if _, err := s.repos.Queries.DeleteFinances(ctx, input); err != nil {
		return err
	}

	return nil
}

// Backstory
func (s *CharacterService) GetBackstory(ctx context.Context, input db.GetBackstoryByCharacterParams) (model.BackstoryModel, error) {
	backstory, err := s.repos.Queries.GetBackstoryByCharacter(ctx, input)
	if err != nil {
		return model.BackstoryModel{}, err
	}

	items, err := s.repos.Queries.GetBackstoryItems(ctx, db.GetBackstoryItemsParams(input))
	if err != nil {
		return model.BackstoryModel{}, err
	}

	return model.ToBackstoryModel(backstory, items), nil
}

func (s *CharacterService) UpsertBackstory(ctx context.Context, input db.UpsertBackstoryParams) (model.BackstoryModel, error) {
	backstory, err := s.repos.Queries.UpsertBackstory(ctx, input)
	if err != nil {
		return model.BackstoryModel{}, err
	}

	items, err := s.repos.Queries.GetBackstoryItems(ctx, db.GetBackstoryItemsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return model.BackstoryModel{}, err
	}

	return model.ToBackstoryModel(backstory, items), nil
}

func (s *CharacterService) DeleteBackstory(ctx context.Context, input db.DeleteBackstoryParams) error {
	if _, err := s.repos.Queries.DeleteBackstory(ctx, input); err != nil {
		return err
	}

	return nil
}

func (s *CharacterService) GetBackstoryItems(ctx context.Context, input db.GetBackstoryItemsParams) ([]model.BackstoryItemModel, error) {
	items, err := s.repos.Queries.GetBackstoryItems(ctx, input)
	if err != nil {
		return nil, err
	}

	return model.ToBackstoryItemModels(items), nil
}

func (s *CharacterService) GetBackstoryItem(ctx context.Context, input db.GetBackstoryItemParams) (model.BackstoryItemModel, error) {
	item, err := s.repos.Queries.GetBackstoryItem(ctx, input)
	if err != nil {
		return model.BackstoryItemModel{}, err
	}

	return model.ToBackstoryItemModel(item), nil
}

func (s *CharacterService) CreateBackstoryItem(ctx context.Context, input db.CreateBackstoryItemParams) (model.BackstoryItemModel, error) {
	item, err := s.repos.Queries.CreateBackstoryItem(ctx, input)
	if err != nil {
		return model.BackstoryItemModel{}, err
	}

	return model.ToBackstoryItemModel(item), nil
}

func (s *CharacterService) UpdateBackstoryItem(ctx context.Context, input db.UpdateBackstoryItemParams) (model.BackstoryItemModel, error) {
	item, err := s.repos.Queries.UpdateBackstoryItem(ctx, input)
	if err != nil {
		return model.BackstoryItemModel{}, err
	}

	return model.ToBackstoryItemModel(item), nil
}

func (s *CharacterService) DeleteBackstoryItem(ctx context.Context, input db.DeleteBackstoryItemParams) error {
	_, err := s.repos.Queries.DeleteBackstoryItem(ctx, input)
	return err
}

// Skills
func (s *CharacterService) GetSkills(ctx context.Context, input db.GetCharacterSkillsParams) ([]model.SkillModel, error) {
	skills, err := s.repos.Queries.GetCharacterSkills(ctx, input)
	if err != nil {
		return nil, err
	}

	return model.ToCharacterSkillModels(skills), nil
}

func (s *CharacterService) GetSkill(ctx context.Context, input db.GetCharacterSkillParams) (model.SkillModel, error) {
	skill, err := s.repos.Queries.GetCharacterSkill(ctx, input)
	if err != nil {
		return model.SkillModel{}, err
	}

	return model.ToSingleCharacterSkillModel(skill), nil
}

func (s *CharacterService) CreateSkill(ctx context.Context, input db.CreateCharacterSkillParams) (model.SkillModel, error) {
	skill, err := s.repos.Queries.CreateCharacterSkill(ctx, input)
	if err != nil {
		return model.SkillModel{}, err
	}

	return model.ToCreatedCharacterSkillModel(skill), nil
}

func (s *CharacterService) UpdateSkill(ctx context.Context, input db.UpdateCharacterSkillParams) (model.SkillModel, error) {
	skill, err := s.repos.Queries.UpdateCharacterSkill(ctx, input)
	if err != nil {
		return model.SkillModel{}, err
	}

	return model.ToUpdatedCharacterSkillModel(skill), nil
}

func (s *CharacterService) DeleteSkill(ctx context.Context, input db.DeleteCharacterSkillParams) error {
	_, err := s.repos.Queries.DeleteCharacterSkill(ctx, input)
	return err
}

// DerivedStats
func (s *CharacterService) GetDerivedStats(ctx context.Context, input db.GetDerivedStatsParams) (db.DerivedStat, error) {
	derivedStats, err := s.repos.Queries.GetDerivedStats(ctx, input)
	if err != nil {
		return db.DerivedStat{}, err
	}

	return derivedStats, nil
}

func (s *CharacterService) UpsertDerivedStats(ctx context.Context, input db.UpsertDerivedStatsParams) (db.DerivedStat, error) {
	derivedStats, err := s.repos.Queries.UpsertDerivedStats(ctx, input)
	if err != nil {
		return db.DerivedStat{}, err
	}

	return derivedStats, nil
}

func (s *CharacterService) DeleteDerivedStats(ctx context.Context, input db.DeleteDerivedStatsParams) error {
	if _, err := s.repos.Queries.DeleteDerivedStats(ctx, input); err != nil {
		return err
	}

	return nil
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

	character, err := s.repos.Queries.GetCharacter(ctx, db.GetCharacterParams{
		UserID: input.UserID,
		ID:     input.CharacterID,
	})
	if err == nil {
		s.recalculateDerivedStats(ctx, input.UserID, input.CharacterID, character.Age, characteristics, "characteristics_upsert")
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

// Helpers
func (s *CharacterService) validateMagicState(ctx context.Context, input db.UpsertMagicStateParams) error {
	if input.MaxMp != nil && input.CurrentMp != nil {
		if *input.CurrentMp > *input.MaxMp {
			return myErrors.ErrCurrentMagicExceedsMax
		}
		return nil
	}

	if input.MaxMp == nil && input.CurrentMp == nil {
		return nil
	}

	existing, err := s.repos.Queries.GetMagicState(ctx, db.GetMagicStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	maxMp := existing.MaxMp
	if input.MaxMp != nil {
		maxMp = *input.MaxMp
	}

	currentMp := existing.CurrentMp
	if input.CurrentMp != nil {
		currentMp = *input.CurrentMp
	}

	if currentMp > maxMp {
		return myErrors.ErrCurrentMagicExceedsMax
	}

	return nil
}

func (s *CharacterService) validateLuckState(ctx context.Context, input db.UpsertLuckStateParams) error {
	if input.StartingLuck != nil && input.CurrentLuck != nil {
		if *input.CurrentLuck > *input.StartingLuck {
			return myErrors.ErrCurrentLuckExceedsStarting
		}
		return nil
	}

	if input.StartingLuck == nil && input.CurrentLuck == nil {
		return nil
	}

	existing, err := s.repos.Queries.GetLuckState(ctx, db.GetLuckStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	startingLuck := existing.StartingLuck
	if input.StartingLuck != nil {
		startingLuck = *input.StartingLuck
	}

	currentLuck := existing.CurrentLuck
	if input.CurrentLuck != nil {
		currentLuck = *input.CurrentLuck
	}

	if currentLuck > startingLuck {
		return myErrors.ErrCurrentLuckExceedsStarting
	}

	return nil
}

func (s *CharacterService) validateHealthState(ctx context.Context, input db.UpsertHealthStateParams) error {
	if input.MaxHp != nil && input.CurrentHp != nil {
		if *input.CurrentHp > *input.MaxHp {
			return myErrors.ErrCurrentHealthExceedsMax
		}
		return nil
	}

	if input.MaxHp == nil && input.CurrentHp == nil {
		return nil
	}

	existing, err := s.repos.Queries.GetHealthState(ctx, db.GetHealthStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	maxHp := existing.MaxHp
	if input.MaxHp != nil {
		maxHp = *input.MaxHp
	}

	currentHp := existing.CurrentHp
	if input.CurrentHp != nil {
		currentHp = *input.CurrentHp
	}

	if currentHp > maxHp {
		return myErrors.ErrCurrentHealthExceedsMax
	}

	return nil
}

func (s *CharacterService) validateSanityState(ctx context.Context, input db.UpsertSanityStateParams) error {
	if input.MaxSanity != nil && input.CurrentSanity != nil {
		if *input.CurrentSanity > *input.MaxSanity {
			return myErrors.ErrCurrentSanityExceedsMax
		}
		return nil
	}

	if input.MaxSanity == nil && input.CurrentSanity == nil {
		return nil
	}

	existing, err := s.repos.Queries.GetSanityState(ctx, db.GetSanityStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	maxSanity := existing.MaxSanity
	if input.MaxSanity != nil {
		maxSanity = *input.MaxSanity
	}

	currentSanity := existing.CurrentSanity
	if input.CurrentSanity != nil {
		currentSanity = *input.CurrentSanity
	}

	if currentSanity > maxSanity {
		return myErrors.ErrCurrentSanityExceedsMax
	}

	return nil
}

func (s *CharacterService) getCharacteristicsForDerivedStatsRecalculation(
	ctx context.Context,
	userID string,
	characterID pgtype.UUID,
) (db.Characteristic, bool) {
	characteristics, err := s.repos.Queries.GetCharacteristics(ctx, db.GetCharacteristicsParams{
		UserID:      userID,
		CharacterID: characterID,
	})
	if err != nil {
		return db.Characteristic{}, false
	}

	return characteristics, true
}

func (s *CharacterService) recalculateDerivedStats(
	ctx context.Context,
	userID string,
	characterID pgtype.UUID,
	age *int16,
	characteristics db.Characteristic,
	source string,
) {
	if reason, canCalculate := derivedStatsRecalculationReadiness(age, characteristics); !canCalculate {
		s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsAutoRecalculateSkipped{
			UserID:      userID,
			CharacterID: characterID.String(),
			Source:      source,
			Reason:      reason,
		})
		return
	}

	derivedStatsInput := calculators.CalculateDerivedStats(
		userID,
		characterID,
		*age,
		characteristics,
	)

	_, err := s.UpsertDerivedStats(ctx, derivedStatsInput)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsAutoRecalculateFailed{
			UserID:      userID,
			CharacterID: characterID.String(),
			Source:      source,
			Err:         err,
		})
		return
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsAutoRecalculateSucceeded{
		UserID:      userID,
		CharacterID: characterID.String(),
		Source:      source,
	})
}

func derivedStatsRecalculationReadiness(age *int16, characteristics db.Characteristic) (string, bool) {
	if age == nil {
		return "age_missing", false
	}

	if characteristics.Strength == nil || characteristics.Size == nil || characteristics.Dexterity == nil {
		return "required_characteristics_missing", false
	}

	return "", true
}
