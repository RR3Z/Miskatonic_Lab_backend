package skills

import "github.com/jackc/pgx/v5/pgtype"

type GetSkillsInput struct {
	UserID      string
	CharacterID pgtype.UUID
}

type GetSkillInput struct {
	UserID      string
	CharacterID pgtype.UUID
	SkillID     pgtype.UUID
}

type CreateSkillInput struct {
	Name        string
	CategoryID  pgtype.UUID
	BaseValue   int16
	Value       int16
	Checked     bool
	Specialized bool
	SpecialtyID pgtype.UUID
	UserID      string
	CharacterID pgtype.UUID
}

type UpdateSkillInput struct {
	Name        string
	CategoryID  pgtype.UUID
	BaseValue   int16
	Value       int16
	Checked     bool
	Specialized bool
	SpecialtyID pgtype.UUID
	UserID      string
	CharacterID pgtype.UUID
	SkillID     pgtype.UUID
}

type DeleteSkillInput struct {
	UserID      string
	CharacterID pgtype.UUID
	SkillID     pgtype.UUID
}
