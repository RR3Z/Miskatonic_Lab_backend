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
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type CharacterDBData struct {
	Character       db.Character
	Characteristics db.Characteristic
	DerivedStats    db.DerivedStat
	HP              db.HealthState
	MP              db.MagicState
	Sanity          db.SanityState
	Luck            db.LuckState
	Skills          []db.GetSkillsRow
	Backstory       *db.Backstory
	BackstoryItems  []db.BackstoryItem
	Finances        *db.Finance
	Notes           []db.Note
}

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

	Skills          []skills.SkillModel                `json:"skills"`
	Characteristics characteristics.CharacteristicsModel `json:"characteristics"`
	DerivedStats    derivedstats.DerivedStatsModel     `json:"derived_stats"`
	HP              health.HealthModel                 `json:"hp"`
	MP              magic.MagicModel                   `json:"mp"`
	Sanity          sanity.SanityModel                 `json:"sanity"`
	Luck            luck.LuckModel                     `json:"luck"`
	Backstory       backstories.BackstoryModel         `json:"backstory"`
	Finances        finances.FinancesModel             `json:"finances"`
	Notes           []notes.NoteModel                  `json:"notes"`
}

func ToCharacterShortModel(c db.Character) CharacterShortModel {
	return CharacterShortModel{
		ID:         c.ID,
		UserID:     c.UserID,
		Name:       c.Name,
		PlayerName: c.PlayerName,
		Occupation: c.Occupation,
		Age:        c.Age,
		Sex:        c.Sex,
		Residence:  c.Residence,
		Birthplace: c.Birthplace,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

func ToCharacterModel(d CharacterDBData) CharacterModel {
	m := CharacterModel{
		CharacterShortModel: ToCharacterShortModel(d.Character),
	}

	if d.Characteristics.ID.Valid {
		m.Characteristics = characteristics.ToCharacteristicsModel(d.Characteristics)
	}
	if d.DerivedStats.ID.Valid {
		m.DerivedStats = derivedstats.ToDerivedStatsModel(d.DerivedStats)
	}
	if d.HP.ID.Valid {
		m.HP = health.ToHealthModel(d.HP)
	}
	if d.MP.ID.Valid {
		m.MP = magic.ToMagicModel(d.MP)
	}
	if d.Sanity.ID.Valid {
		m.Sanity = sanity.ToSanityModel(d.Sanity)
	}
	if d.Luck.ID.Valid {
		m.Luck = luck.ToLuckModel(d.Luck)
	}

	m.Notes = notes.ToNoteModels(d.Notes)

	if len(d.Skills) > 0 {
		skillModels := make([]skills.SkillModel, len(d.Skills))
		for i, s := range d.Skills {
			skillModels[i] = skills.ToSkillModel(s)
		}
		m.Skills = skillModels
	}

	if d.Backstory != nil {
		m.Backstory = backstories.ToBackstoryModel(*d.Backstory, d.BackstoryItems)
	}

	if d.Finances != nil {
		var creditRating *skills.SkillModel
		if d.Finances.CreditRatingSkillID.Valid {
			for _, skill := range m.Skills {
				if skill.ID == d.Finances.CreditRatingSkillID {
					creditRating = &skill
					break
				}
			}
		}
		m.Finances = finances.ToFinancesModel(*d.Finances, creditRating)
	}

	return m
}
