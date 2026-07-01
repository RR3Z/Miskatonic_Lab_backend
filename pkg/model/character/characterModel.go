package character

import (
	rootModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
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

type CharacterModel struct {
	CharacterShortModel

	Skills          []rootModel.SkillModel   `json:"skills"`
	Characteristics db.Characteristic        `json:"characteristics"`
	DerivedStats    db.DerivedStat           `json:"derived_stats"`
	HP              db.HealthState           `json:"hp"`
	MP              db.MagicState            `json:"mp"`
	Sanity          db.SanityState           `json:"sanity"`
	Luck            db.LuckState             `json:"luck"`
	Backstory       rootModel.BackstoryModel `json:"backstory"`
	Finances        rootModel.FinancesModel  `json:"finances"`
	Notes           []db.Note                `json:"notes"`
}

func ToCharacterModel(d CharacterDBData) CharacterModel {
	m := CharacterModel{
		CharacterShortModel: ToCharacterShortModel(d.Character),
	}

	if d.Characteristics.ID.Valid {
		m.Characteristics = d.Characteristics
	}
	if d.DerivedStats.ID.Valid {
		m.DerivedStats = d.DerivedStats
	}
	if d.HP.ID.Valid {
		m.HP = d.HP
	}
	if d.MP.ID.Valid {
		m.MP = d.MP
	}
	if d.Sanity.ID.Valid {
		m.Sanity = d.Sanity
	}
	if d.Luck.ID.Valid {
		m.Luck = d.Luck
	}

	m.Notes = d.Notes

	if len(d.Skills) > 0 {
		skills := make([]rootModel.SkillModel, len(d.Skills))
		for i, s := range d.Skills {
			skills[i] = rootModel.ToSkillModel(s)
		}
		m.Skills = skills
	}

	if d.Backstory != nil {
		m.Backstory = rootModel.ToBackstoryModel(*d.Backstory, d.BackstoryItems)
	}

	if d.Finances != nil {
		var creditRating *rootModel.SkillModel
		if d.Finances.CreditRatingSkillID.Valid {
			for _, skill := range m.Skills {
				if skill.ID == d.Finances.CreditRatingSkillID {
					creditRating = &skill
					break
				}
			}
		}
		m.Finances = rootModel.ToFinancesModel(*d.Finances, creditRating)
	}

	return m
}
