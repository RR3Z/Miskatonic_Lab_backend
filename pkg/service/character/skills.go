package character

import (
	"context"

	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
)

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
	if err := validateSkillInput(input.Name, input.BaseValue, input.Value); err != nil {
		return skillsDTO.SkillModel{}, err
	}

	skill, err := s.repos.Queries.CreateCharacterSkill(ctx, db.CreateCharacterSkillParams{
		Name:        input.Name,
		BaseValue:   input.BaseValue,
		Value:       input.Value,
		Checked:     input.Checked,
		IsProtected: false,
		BaseRule:    nil,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return skillsDTO.SkillModel{}, characterErrors.MapCharacterConstraintError(err)
	}

	return skillsDTO.ToCreatedCharacterSkillModel(skill), nil
}

func (s *CharacterService) UpdateSkill(ctx context.Context, input skillsDTO.UpdateSkillInput) (skillsDTO.SkillModel, error) {
	if err := validateSkillInput(input.Name, input.BaseValue, input.Value); err != nil {
		return skillsDTO.SkillModel{}, err
	}

	existing, err := s.repos.Queries.GetCharacterSkill(ctx, db.GetCharacterSkillParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		SkillID:     input.SkillID,
	})
	if err != nil {
		return skillsDTO.SkillModel{}, err
	}
	if existing.IsProtected && (existing.Name != input.Name || existing.BaseValue != input.BaseValue) {
		return skillsDTO.SkillModel{}, characterErrors.ErrProtectedSkill
	}

	skill, err := s.repos.Queries.UpdateCharacterSkill(ctx, db.UpdateCharacterSkillParams{
		Name:        input.Name,
		BaseValue:   input.BaseValue,
		Value:       input.Value,
		Checked:     input.Checked,
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
	existing, err := s.repos.Queries.GetCharacterSkill(ctx, db.GetCharacterSkillParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		SkillID:     input.SkillID,
	})
	if err != nil {
		return err
	}
	if existing.IsProtected {
		return characterErrors.ErrProtectedSkill
	}

	_, err = s.repos.Queries.DeleteCharacterSkill(ctx, db.DeleteCharacterSkillParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		SkillID:     input.SkillID,
	})
	return err
}
