package characters

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/jackc/pgx/v5/pgtype"
)

type CharacterShortModel struct {
	ID     pgtype.UUID `json:"id"`
	UserID string      `json:"user_id"`

	Name       string  `json:"name"`
	PlayerName *string `json:"player_name"`
	Occupation *string `json:"occupation"`
	Age        *int16  `json:"age"`
	Sex        *string `json:"sex"`
	Residence  *string `json:"residence"`
	Birthplace *string `json:"birthplace"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type CharacterModel struct {
	CharacterShortModel

	Skills          []skills.SkillModel             `json:"skills"`
	Characteristics characteristics.CharacteristicsModel `json:"characteristics"`
	DerivedStats    derivedstats.DerivedStatsModel  `json:"derived_stats"`
	HP              health.HealthModel              `json:"hp"`
	MP              magic.MagicModel                `json:"mp"`
	Sanity          sanity.SanityModel              `json:"sanity"`
	Luck            luck.LuckModel                  `json:"luck"`
	Backstory       backstories.BackstoryModel      `json:"backstory"`
	Finances        finances.FinancesModel          `json:"finances"`
	Notes           []notes.NoteModel               `json:"notes"`
}
