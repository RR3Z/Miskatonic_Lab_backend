package skillsDTO

import "github.com/jackc/pgx/v5/pgtype"

type SkillRequest struct {
	Name        string      `json:"name"`
	CategoryID  pgtype.UUID `json:"category_id"`
	BaseValue   int16       `json:"base_value"`
	Value       int16       `json:"value"`
	Checked     bool        `json:"checked"`
	Specialized bool        `json:"specialized"`
	SpecialtyID pgtype.UUID `json:"specialty_id"`
}
