package character

import (
	"context"

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
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

// Implements ICharacter
type EventPublishingCharacterService struct {
	next      ICharacter
	publisher events.EventPublisher
}

func NewEventPublishingCharacterService(next ICharacter, publisher events.EventPublisher) *EventPublishingCharacterService {
	return &EventPublishingCharacterService{
		next:      next,
		publisher: publisher,
	}
}

// Characters
func (s *EventPublishingCharacterService) GetAllCharacters(ctx context.Context, userID string) ([]characterDTO.CharacterShortModel, error) {
	characters, err := s.next.GetAllCharacters(ctx, userID)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharactersListFailed{
			UserID: userID,
			Err:    err,
		})

		return nil, err
	}

	s.publisher.Publish(ctx, characterEvents.CharactersListSucceeded{
		UserID: userID,
		Count:  len(characters),
	})

	return characters, nil
}

func (s *EventPublishingCharacterService) GetCharacter(ctx context.Context, input characterDTO.GetCharacterInput) (characterDTO.CharacterModel, error) {
	character, err := s.next.GetCharacter(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return characterDTO.CharacterModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		Name:        character.Name,
	})

	return character, nil
}

func (s *EventPublishingCharacterService) CreateCharacter(ctx context.Context, input characterDTO.CreateCharacterInput) (characterDTO.CharacterShortModel, error) {
	character, err := s.next.CreateCharacter(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCreateFailed{
			UserID: input.UserID,
			Err:    err,
		})
		return characterDTO.CharacterShortModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCreateSucceeded{
		UserID:      input.UserID,
		CharacterID: character.ID.String(),
		Name:        character.Name,
	})

	return character, nil
}

func (s *EventPublishingCharacterService) UpdateCharacter(ctx context.Context, input characterDTO.UpdateCharacterInput) (characterDTO.CharacterShortModel, error) {
	character, err := s.next.UpdateCharacter(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterUpdateFailed{
			UserID:      input.UserID,
			CharacterID: input.ID.String(),
			Err:         err,
		})
		return characterDTO.CharacterShortModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterUpdateSucceeded{
		UserID:      input.UserID,
		CharacterID: character.ID.String(),
		Name:        character.Name,
	})

	return character, nil
}

func (s *EventPublishingCharacterService) DeleteCharacter(ctx context.Context, input characterDTO.DeleteCharacterInput) error {
	if err := s.next.DeleteCharacter(ctx, input); err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.ID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.ID.String(),
	})

	return nil
}

// Health
func (s *EventPublishingCharacterService) GetHealth(ctx context.Context, input healthDTO.GetHealthInput) (db.HealthState, error) {
	health, err := s.next.GetHealth(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterHealthGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.HealthState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterHealthGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return health, nil
}

func (s *EventPublishingCharacterService) UpsertHealth(ctx context.Context, input healthDTO.UpsertHealthInput) (db.HealthState, error) {
	health, err := s.next.UpsertHealth(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterHealthUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.HealthState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterHealthUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return health, nil
}

func (s *EventPublishingCharacterService) DeleteHealth(ctx context.Context, input healthDTO.DeleteHealthInput) error {
	err := s.next.DeleteHealth(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterHealthDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterHealthDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Sanity
func (s *EventPublishingCharacterService) GetSanity(ctx context.Context, input sanityDTO.GetSanityInput) (db.SanityState, error) {
	sanity, err := s.next.GetSanity(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSanityGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.SanityState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSanityGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return sanity, nil
}

func (s *EventPublishingCharacterService) UpsertSanity(ctx context.Context, input sanityDTO.UpsertSanityInput) (db.SanityState, error) {
	sanity, err := s.next.UpsertSanity(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSanityUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.SanityState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSanityUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return sanity, nil
}

func (s *EventPublishingCharacterService) DeleteSanity(ctx context.Context, input sanityDTO.DeleteSanityInput) error {
	err := s.next.DeleteSanity(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSanityDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSanityDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Magic
func (s *EventPublishingCharacterService) GetMagic(ctx context.Context, input magicDTO.GetMagicInput) (db.MagicState, error) {
	magic, err := s.next.GetMagic(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterMagicGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.MagicState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterMagicGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return magic, nil
}

func (s *EventPublishingCharacterService) UpsertMagic(ctx context.Context, input magicDTO.UpsertMagicInput) (db.MagicState, error) {
	magic, err := s.next.UpsertMagic(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterMagicUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.MagicState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterMagicUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return magic, nil
}

func (s *EventPublishingCharacterService) DeleteMagic(ctx context.Context, input magicDTO.DeleteMagicInput) error {
	err := s.next.DeleteMagic(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterMagicDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterMagicDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Luck
func (s *EventPublishingCharacterService) GetLuck(ctx context.Context, input luckDTO.GetLuckInput) (db.LuckState, error) {
	luck, err := s.next.GetLuck(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterLuckGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.LuckState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterLuckGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return luck, nil
}

func (s *EventPublishingCharacterService) UpsertLuck(ctx context.Context, input luckDTO.UpsertLuckInput) (db.LuckState, error) {
	luck, err := s.next.UpsertLuck(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterLuckUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.LuckState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterLuckUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return luck, nil
}

func (s *EventPublishingCharacterService) DeleteLuck(ctx context.Context, input luckDTO.DeleteLuckInput) error {
	err := s.next.DeleteLuck(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterLuckDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterLuckDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Finances
func (s *EventPublishingCharacterService) GetFinances(ctx context.Context, input financesDTO.GetFinancesInput) (db.Finance, error) {
	finances, err := s.next.GetFinances(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterFinancesGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Finance{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterFinancesGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return finances, nil
}

func (s *EventPublishingCharacterService) UpsertFinances(ctx context.Context, input financesDTO.UpsertFinancesInput) (db.Finance, error) {
	finances, err := s.next.UpsertFinances(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterFinancesUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Finance{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterFinancesUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return finances, nil
}

func (s *EventPublishingCharacterService) DeleteFinances(ctx context.Context, input financesDTO.DeleteFinancesInput) error {
	err := s.next.DeleteFinances(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterFinancesDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterFinancesDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Backstory
func (s *EventPublishingCharacterService) GetBackstory(ctx context.Context, input backstoriesDTO.GetBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	backstory, err := s.next.GetBackstory(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return backstoriesDTO.BackstoryModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return backstory, nil
}

func (s *EventPublishingCharacterService) UpsertBackstory(ctx context.Context, input backstoriesDTO.UpsertBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	backstory, err := s.next.UpsertBackstory(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return backstoriesDTO.BackstoryModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return backstory, nil
}

func (s *EventPublishingCharacterService) DeleteBackstory(ctx context.Context, input backstoriesDTO.DeleteBackstoryInput) error {
	err := s.next.DeleteBackstory(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

func (s *EventPublishingCharacterService) GetBackstoryItems(ctx context.Context, input backstoriesDTO.GetBackstoryItemsInput) ([]backstoriesDTO.BackstoryItemModel, error) {
	items, err := s.next.GetBackstoryItems(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemsListFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return nil, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemsListSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		Count:       len(items),
	})

	return items, nil
}

func (s *EventPublishingCharacterService) GetBackstoryItem(ctx context.Context, input backstoriesDTO.GetBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	item, err := s.next.GetBackstoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemGetFailed{
			UserID:          input.UserID,
			CharacterID:     input.CharacterID.String(),
			BackstoryItemID: input.BackstoryItemID.String(),
			Err:             err,
		})
		return backstoriesDTO.BackstoryItemModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemGetSucceeded{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID.String(),
		BackstoryItemID: input.BackstoryItemID.String(),
		Section:         item.Section,
		Title:           item.Title,
	})

	return item, nil
}

func (s *EventPublishingCharacterService) CreateBackstoryItem(ctx context.Context, input backstoriesDTO.CreateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	item, err := s.next.CreateBackstoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemCreateFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return backstoriesDTO.BackstoryItemModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemCreateSucceeded{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID.String(),
		BackstoryItemID: item.ID.String(),
		Section:         item.Section,
		Title:           item.Title,
	})

	return item, nil
}

func (s *EventPublishingCharacterService) UpdateBackstoryItem(ctx context.Context, input backstoriesDTO.UpdateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	item, err := s.next.UpdateBackstoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemUpdateFailed{
			UserID:          input.UserID,
			CharacterID:     input.CharacterID.String(),
			BackstoryItemID: input.BackstoryItemID.String(),
			Err:             err,
		})
		return backstoriesDTO.BackstoryItemModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemUpdateSucceeded{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID.String(),
		BackstoryItemID: input.BackstoryItemID.String(),
		Section:         item.Section,
		Title:           item.Title,
	})

	return item, nil
}

func (s *EventPublishingCharacterService) DeleteBackstoryItem(ctx context.Context, input backstoriesDTO.DeleteBackstoryItemInput) error {
	err := s.next.DeleteBackstoryItem(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemDeleteFailed{
			UserID:          input.UserID,
			CharacterID:     input.CharacterID.String(),
			BackstoryItemID: input.BackstoryItemID.String(),
			Err:             err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterBackstoryItemDeleteSucceeded{
		UserID:          input.UserID,
		CharacterID:     input.CharacterID.String(),
		BackstoryItemID: input.BackstoryItemID.String(),
	})

	return nil
}

// Skills
func (s *EventPublishingCharacterService) GetSkills(ctx context.Context, input skillsDTO.GetSkillsInput) ([]skillsDTO.SkillModel, error) {
	skills, err := s.next.GetSkills(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSkillsListFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return nil, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSkillsListSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		Count:       len(skills),
	})

	return skills, nil
}

func (s *EventPublishingCharacterService) GetSkill(ctx context.Context, input skillsDTO.GetSkillInput) (skillsDTO.SkillModel, error) {
	skill, err := s.next.GetSkill(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSkillGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			SkillID:     input.SkillID.String(),
			Err:         err,
		})
		return skillsDTO.SkillModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSkillGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		SkillID:     input.SkillID.String(),
		Name:        skill.Name,
	})

	return skill, nil
}

func (s *EventPublishingCharacterService) CreateSkill(ctx context.Context, input skillsDTO.CreateSkillInput) (skillsDTO.SkillModel, error) {
	skill, err := s.next.CreateSkill(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSkillCreateFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return skillsDTO.SkillModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSkillCreateSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		SkillID:     skill.ID.String(),
		Name:        skill.Name,
	})

	return skill, nil
}

func (s *EventPublishingCharacterService) UpdateSkill(ctx context.Context, input skillsDTO.UpdateSkillInput) (skillsDTO.SkillModel, error) {
	skill, err := s.next.UpdateSkill(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSkillUpdateFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			SkillID:     input.SkillID.String(),
			Err:         err,
		})
		return skillsDTO.SkillModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSkillUpdateSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		SkillID:     input.SkillID.String(),
		Name:        skill.Name,
	})

	return skill, nil
}

func (s *EventPublishingCharacterService) DeleteSkill(ctx context.Context, input skillsDTO.DeleteSkillInput) error {
	err := s.next.DeleteSkill(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSkillDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			SkillID:     input.SkillID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSkillDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		SkillID:     input.SkillID.String(),
	})

	return nil
}

// DerivedStats
func (s *EventPublishingCharacterService) GetDerivedStats(ctx context.Context, input derivedStatsDTO.GetDerivedStatsInput) (db.DerivedStat, error) {
	derivedStats, err := s.next.GetDerivedStats(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.DerivedStat{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return derivedStats, nil
}

func (s *EventPublishingCharacterService) UpsertDerivedStats(ctx context.Context, input derivedStatsDTO.UpsertDerivedStatsInput) (db.DerivedStat, error) {
	derivedStats, err := s.next.UpsertDerivedStats(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.DerivedStat{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return derivedStats, nil
}

func (s *EventPublishingCharacterService) DeleteDerivedStats(ctx context.Context, input derivedStatsDTO.DeleteDerivedStatsInput) error {
	err := s.next.DeleteDerivedStats(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDerivedStatsDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Characteristics
func (s *EventPublishingCharacterService) GetCharacteristics(ctx context.Context, input characteristicsDTO.GetCharacteristicsInput) (db.Characteristic, error) {
	characteristics, err := s.next.GetCharacteristics(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Characteristic{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return characteristics, nil
}

func (s *EventPublishingCharacterService) UpsertCharacteristics(ctx context.Context, input characteristicsDTO.UpsertCharacteristicsInput) (db.Characteristic, error) {
	characteristics, err := s.next.UpsertCharacteristics(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Characteristic{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return characteristics, nil
}

func (s *EventPublishingCharacterService) DeleteCharacteristics(ctx context.Context, input characteristicsDTO.DeleteCharacteristicsInput) error {
	err := s.next.DeleteCharacteristics(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Notes
func (s *EventPublishingCharacterService) GetNotes(ctx context.Context, input notesDTO.GetNotesInput) ([]db.Note, error) {
	notes, err := s.next.GetNotes(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterNotesListFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return nil, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterNotesListSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		Count:       len(notes),
	})

	return notes, nil
}

func (s *EventPublishingCharacterService) GetNote(ctx context.Context, input notesDTO.GetNoteInput) (db.Note, error) {
	note, err := s.next.GetNote(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterNoteGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			NoteID:      input.NoteID.String(),
			Err:         err,
		})
		return db.Note{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterNoteGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		NoteID:      input.NoteID.String(),
		Title:       note.Title,
	})

	return note, nil
}

func (s *EventPublishingCharacterService) CreateNote(ctx context.Context, input notesDTO.CreateNoteInput) (db.Note, error) {
	note, err := s.next.CreateNote(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterNoteCreateFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Note{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterNoteCreateSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		NoteID:      note.ID.String(),
		Title:       note.Title,
	})

	return note, nil
}

func (s *EventPublishingCharacterService) UpdateNote(ctx context.Context, input notesDTO.UpdateNoteInput) (db.Note, error) {
	note, err := s.next.UpdateNote(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterNoteUpdateFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			NoteID:      input.NoteID.String(),
			Err:         err,
		})
		return db.Note{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterNoteUpdateSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		NoteID:      input.NoteID.String(),
		Title:       note.Title,
	})

	return note, nil
}

func (s *EventPublishingCharacterService) DeleteNote(ctx context.Context, input notesDTO.DeleteNoteInput) error {
	err := s.next.DeleteNote(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterNoteDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			NoteID:      input.NoteID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterNoteDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		NoteID:      input.NoteID.String(),
	})

	return nil
}
