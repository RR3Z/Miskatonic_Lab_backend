package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
)

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
