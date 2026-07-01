package character

import (
	"context"
	"errors"
	"strings"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	derivedStatsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	luckDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	sanityDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/calculators"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
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
	if strings.TrimSpace(input.Name) == "" {
		return characterDTO.CharacterShortModel{}, characterErrors.ErrNameRequired
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
	if strings.TrimSpace(input.Name) == "" {
		return characterDTO.CharacterShortModel{}, characterErrors.ErrNameRequired
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

// Health
func (s *CharacterService) GetHealth(ctx context.Context, input healthDTO.GetHealthInput) (db.HealthState, error) {
	health, err := s.repos.Queries.GetHealthState(ctx, db.GetHealthStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.HealthState{}, err
	}

	return health, nil
}

func (s *CharacterService) UpsertHealth(ctx context.Context, input healthDTO.UpsertHealthInput) (db.HealthState, error) {
	if err := s.validateHealthState(ctx, input); err != nil {
		return db.HealthState{}, err
	}

	health, err := s.repos.Queries.UpsertHealthState(ctx, db.UpsertHealthStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		MaxHp:       input.MaxHp,
		CurrentHp:   input.CurrentHp,
		MajorWound:  input.MajorWound,
		Unconscious: input.Unconscious,
		Dying:       input.Dying,
		Dead:        input.Dead,
	})
	if err != nil {
		return db.HealthState{}, characterErrors.MapCharacterConstraintError(err)
	}

	return health, nil
}

func (s *CharacterService) DeleteHealth(ctx context.Context, input healthDTO.DeleteHealthInput) error {
	if _, err := s.repos.Queries.DeleteHealthState(ctx, db.DeleteHealthStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}

// Sanity
func (s *CharacterService) GetSanity(ctx context.Context, input sanityDTO.GetSanityInput) (db.SanityState, error) {
	sanity, err := s.repos.Queries.GetSanityState(ctx, db.GetSanityStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.SanityState{}, err
	}

	return sanity, nil
}

func (s *CharacterService) UpsertSanity(ctx context.Context, input sanityDTO.UpsertSanityInput) (db.SanityState, error) {
	if err := s.validateSanityState(ctx, input); err != nil {
		return db.SanityState{}, err
	}

	sanity, err := s.repos.Queries.UpsertSanityState(ctx, db.UpsertSanityStateParams{
		UserID:        input.UserID,
		CharacterID:   input.CharacterID,
		MaxSanity:     input.MaxSanity,
		CurrentSanity: input.CurrentSanity,
		TempInsanity:  input.TempInsanity,
		IndefInsanity: input.IndefInsanity,
	})
	if err != nil {
		return db.SanityState{}, characterErrors.MapCharacterConstraintError(err)
	}

	return sanity, nil
}

func (s *CharacterService) DeleteSanity(ctx context.Context, input sanityDTO.DeleteSanityInput) error {
	if _, err := s.repos.Queries.DeleteSanityState(ctx, db.DeleteSanityStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}

// Magic
func (s *CharacterService) GetMagic(ctx context.Context, input magicDTO.GetMagicInput) (db.MagicState, error) {
	magic, err := s.repos.Queries.GetMagicState(ctx, db.GetMagicStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.MagicState{}, err
	}

	return magic, nil
}

func (s *CharacterService) UpsertMagic(ctx context.Context, input magicDTO.UpsertMagicInput) (db.MagicState, error) {
	if err := s.validateMagicState(ctx, input); err != nil {
		return db.MagicState{}, err
	}

	magic, err := s.repos.Queries.UpsertMagicState(ctx, db.UpsertMagicStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		MaxMp:       input.MaxMp,
		CurrentMp:   input.CurrentMp,
	})
	if err != nil {
		return db.MagicState{}, characterErrors.MapCharacterConstraintError(err)
	}

	return magic, nil
}

func (s *CharacterService) DeleteMagic(ctx context.Context, input magicDTO.DeleteMagicInput) error {
	if _, err := s.repos.Queries.DeleteMagicState(ctx, db.DeleteMagicStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}

// Luck
func (s *CharacterService) GetLuck(ctx context.Context, input luckDTO.GetLuckInput) (db.LuckState, error) {
	luck, err := s.repos.Queries.GetLuckState(ctx, db.GetLuckStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.LuckState{}, err
	}

	return luck, nil
}

func (s *CharacterService) UpsertLuck(ctx context.Context, input luckDTO.UpsertLuckInput) (db.LuckState, error) {
	if err := s.validateLuckState(ctx, input); err != nil {
		return db.LuckState{}, err
	}

	luck, err := s.repos.Queries.UpsertLuckState(ctx, db.UpsertLuckStateParams{
		UserID:       input.UserID,
		CharacterID:  input.CharacterID,
		StartingLuck: input.StartingLuck,
		CurrentLuck:  input.CurrentLuck,
	})
	if err != nil {
		return db.LuckState{}, characterErrors.MapCharacterConstraintError(err)
	}

	return luck, nil
}

func (s *CharacterService) DeleteLuck(ctx context.Context, input luckDTO.DeleteLuckInput) error {
	if _, err := s.repos.Queries.DeleteLuckState(ctx, db.DeleteLuckStateParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}

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
	finances, err := s.repos.Queries.UpsertFinances(ctx, db.UpsertFinancesParams{
		UserID:              input.UserID,
		CharacterID:         input.CharacterID,
		SpendingLimit:       input.SpendingLimit,
		Cash:                input.Cash,
		Assets:              input.Assets,
		CreditRatingSkillID: input.CreditRatingSkillID,
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

// Backstory
func (s *CharacterService) GetBackstory(ctx context.Context, input backstoriesDTO.GetBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	backstory, err := s.repos.Queries.GetBackstoryByCharacter(ctx, db.GetBackstoryByCharacterParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryModel{}, err
	}

	items, err := s.repos.Queries.GetBackstoryItems(ctx, db.GetBackstoryItemsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryModel{}, err
	}

	return backstoriesDTO.ToBackstoryModel(backstory, items), nil
}

func (s *CharacterService) UpsertBackstory(ctx context.Context, input backstoriesDTO.UpsertBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	backstory, err := s.repos.Queries.UpsertBackstory(ctx, db.UpsertBackstoryParams{
		UserID:              input.UserID,
		CharacterID:         input.CharacterID,
		PersonalDescription: input.PersonalDescription,
	})
	if err != nil {
		return backstoriesDTO.BackstoryModel{}, err
	}

	items, err := s.repos.Queries.GetBackstoryItems(ctx, db.GetBackstoryItemsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryModel{}, err
	}

	return backstoriesDTO.ToBackstoryModel(backstory, items), nil
}

func (s *CharacterService) DeleteBackstory(ctx context.Context, input backstoriesDTO.DeleteBackstoryInput) error {
	if _, err := s.repos.Queries.DeleteBackstory(ctx, db.DeleteBackstoryParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *CharacterService) GetBackstoryItems(ctx context.Context, input backstoriesDTO.GetBackstoryItemsInput) ([]backstoriesDTO.BackstoryItemModel, error) {
	items, err := s.repos.Queries.GetBackstoryItems(ctx, db.GetBackstoryItemsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return nil, err
	}

	return backstoriesDTO.ToBackstoryItemModels(items), nil
}

func (s *CharacterService) GetBackstoryItem(ctx context.Context, input backstoriesDTO.GetBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	item, err := s.repos.Queries.GetBackstoryItem(ctx, db.GetBackstoryItemParams{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID,
		BackstoryItemID: input.BackstoryItemID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryItemModel{}, err
	}

	return backstoriesDTO.ToBackstoryItemModel(item), nil
}

func (s *CharacterService) CreateBackstoryItem(ctx context.Context, input backstoriesDTO.CreateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	item, err := s.repos.Queries.CreateBackstoryItem(ctx, db.CreateBackstoryItemParams{
		Section:     input.Section,
		Title:       input.Title,
		Text:        input.Text,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryItemModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	return backstoriesDTO.ToBackstoryItemModel(item), nil
}

func (s *CharacterService) UpdateBackstoryItem(ctx context.Context, input backstoriesDTO.UpdateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	item, err := s.repos.Queries.UpdateBackstoryItem(ctx, db.UpdateBackstoryItemParams{
		Section:         input.Section,
		Title:           input.Title,
		Text:            input.Text,
		UserID:          input.UserID,
		CharacterID:     input.CharacterID,
		BackstoryItemID: input.BackstoryItemID,
	})
	if err != nil {
		return backstoriesDTO.BackstoryItemModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	return backstoriesDTO.ToBackstoryItemModel(item), nil
}

func (s *CharacterService) DeleteBackstoryItem(ctx context.Context, input backstoriesDTO.DeleteBackstoryItemInput) error {
	_, err := s.repos.Queries.DeleteBackstoryItem(ctx, db.DeleteBackstoryItemParams{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID,
		BackstoryItemID: input.BackstoryItemID,
	})
	return err
}

// Skills
func (s *CharacterService) GetSkills(ctx context.Context, input skillsDTO.GetSkillsInput) ([]skillsDTO.SkillModel, error) {
	dbSkills, err := s.repos.Queries.GetCharacterSkills(ctx, db.GetCharacterSkillsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return nil, err
	}

	return skillsDTO.ToCharacterSkillModels(dbSkills), nil
}

func (s *CharacterService) GetSkill(ctx context.Context, input skillsDTO.GetSkillInput) (skillsDTO.SkillModel, error) {
	skill, err := s.repos.Queries.GetCharacterSkill(ctx, db.GetCharacterSkillParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		SkillID:     input.SkillID,
	})
	if err != nil {
		return skillsDTO.SkillModel{}, err
	}

	return skillsDTO.ToSingleCharacterSkillModel(skill), nil
}

func (s *CharacterService) CreateSkill(ctx context.Context, input skillsDTO.CreateSkillInput) (skillsDTO.SkillModel, error) {
	skill, err := s.repos.Queries.CreateCharacterSkill(ctx, db.CreateCharacterSkillParams{
		Name:        input.Name,
		CategoryID:  input.CategoryID,
		BaseValue:   input.BaseValue,
		Value:       input.Value,
		Checked:     input.Checked,
		Specialized: input.Specialized,
		SpecialtyID: input.SpecialtyID,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return skillsDTO.SkillModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	return skillsDTO.ToCreatedCharacterSkillModel(skill), nil
}

func (s *CharacterService) UpdateSkill(ctx context.Context, input skillsDTO.UpdateSkillInput) (skillsDTO.SkillModel, error) {
	skill, err := s.repos.Queries.UpdateCharacterSkill(ctx, db.UpdateCharacterSkillParams{
		Name:        input.Name,
		CategoryID:  input.CategoryID,
		BaseValue:   input.BaseValue,
		Value:       input.Value,
		Checked:     input.Checked,
		Specialized: input.Specialized,
		SpecialtyID: input.SpecialtyID,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		SkillID:     input.SkillID,
	})
	if err != nil {
		return skillsDTO.SkillModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	return skillsDTO.ToUpdatedCharacterSkillModel(skill), nil
}

func (s *CharacterService) DeleteSkill(ctx context.Context, input skillsDTO.DeleteSkillInput) error {
	_, err := s.repos.Queries.DeleteCharacterSkill(ctx, db.DeleteCharacterSkillParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		SkillID:     input.SkillID,
	})
	return characterErrors.MapDeleteSkillError(err)
}

// DerivedStats
func (s *CharacterService) GetDerivedStats(ctx context.Context, input derivedStatsDTO.GetDerivedStatsInput) (db.DerivedStat, error) {
	derivedStats, err := s.repos.Queries.GetDerivedStats(ctx, db.GetDerivedStatsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.DerivedStat{}, err
	}

	return derivedStats, nil
}

func (s *CharacterService) UpsertDerivedStats(ctx context.Context, input derivedStatsDTO.UpsertDerivedStatsInput) (db.DerivedStat, error) {
	derivedStats, err := s.repos.Queries.UpsertDerivedStats(ctx, db.UpsertDerivedStatsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		Speed:       input.Speed,
		Physique:    input.Physique,
		DamageBonus: input.DamageBonus,
		DodgeValue:  input.DodgeValue,
	})
	if err != nil {
		return db.DerivedStat{}, characterErrors.MapCharacterConstraintError(err)
	}

	return derivedStats, nil
}

func (s *CharacterService) DeleteDerivedStats(ctx context.Context, input derivedStatsDTO.DeleteDerivedStatsInput) error {
	if _, err := s.repos.Queries.DeleteDerivedStats(ctx, db.DeleteDerivedStatsParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}); err != nil {
		return err
	}

	return nil
}

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

// Notes
func (s *CharacterService) GetNotes(ctx context.Context, input notesDTO.GetNotesInput) ([]db.Note, error) {
	notes, err := s.repos.Queries.GetNotes(ctx, db.GetNotesParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *CharacterService) GetNote(ctx context.Context, input notesDTO.GetNoteInput) (db.Note, error) {
	note, err := s.repos.Queries.GetNote(ctx, db.GetNoteParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		NoteID:      input.NoteID,
	})
	if err != nil {
		return db.Note{}, err
	}

	return note, nil
}

func (s *CharacterService) CreateNote(ctx context.Context, input notesDTO.CreateNoteInput) (db.Note, error) {
	note, err := s.repos.Queries.CreateNote(ctx, db.CreateNoteParams{
		Title:       input.Title,
		Body:        input.Body,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.Note{}, err
	}

	return note, nil
}

func (s *CharacterService) UpdateNote(ctx context.Context, input notesDTO.UpdateNoteInput) (db.Note, error) {
	note, err := s.repos.Queries.UpdateNote(ctx, db.UpdateNoteParams{
		Title:       input.Title,
		Body:        input.Body,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		NoteID:      input.NoteID,
	})
	if err != nil {
		return db.Note{}, err
	}

	return note, nil
}

func (s *CharacterService) DeleteNote(ctx context.Context, input notesDTO.DeleteNoteInput) error {
	_, err := s.repos.Queries.DeleteNote(ctx, db.DeleteNoteParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		NoteID:      input.NoteID,
	})
	return err
}

// Helpers
func (s *CharacterService) validateMagicState(ctx context.Context, input magicDTO.UpsertMagicInput) error {
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

func (s *CharacterService) validateLuckState(ctx context.Context, input luckDTO.UpsertLuckInput) error {
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

func (s *CharacterService) validateHealthState(ctx context.Context, input healthDTO.UpsertHealthInput) error {
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

func (s *CharacterService) validateSanityState(ctx context.Context, input sanityDTO.UpsertSanityInput) error {
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

	_, err := s.UpsertDerivedStats(ctx, derivedStatsDTO.UpsertDerivedStatsInput{
		UserID:      derivedStatsInput.UserID,
		CharacterID: derivedStatsInput.CharacterID,
		Speed:       derivedStatsInput.Speed,
		Physique:    derivedStatsInput.Physique,
		DamageBonus: derivedStatsInput.DamageBonus,
		DodgeValue:  derivedStatsInput.DodgeValue,
	})
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
