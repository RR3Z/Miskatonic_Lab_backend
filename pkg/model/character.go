package model

import (
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

type CharacterModel struct {
	ID     pgtype.UUID `json:"id"`
	UserID string      `json:"user_id"`

	Name            string            `json:"name"`
	PlayerName      *string           `json:"player_name"`
	Occupation      *string           `json:"occupation"`
	Age             *int16            `json:"age"`
	Sex             *string           `json:"sex"`
	Residence       *string           `json:"residence"`
	Birthplace      *string           `json:"birthplace"`
	Skills          []SkillModel      `json:"skills"`
	Characteristics db.Characteristic `json:"characteristics"`
	DerivedStats    db.DerivedStat    `json:"derived_stats"`
	HP              db.HealthState    `json:"hp"`
	MP              db.MagicState     `json:"mp"`
	Sanity          db.SanityState    `json:"sanity"`
	Luck            db.LuckState      `json:"luck"`
	Backstory       BackstoryModel    `json:"backstory"`
	Finances        FinancesModel     `json:"finances"`
	Notes           []db.Note         `json:"notes"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type GetCharacterInput struct {
	UserID      string      `json:"user_id"`
	CharacterID pgtype.UUID `json:"character_id"`
}

type CreateCharacterInput struct {
}

type UpdateCharacterInput struct {
}

type DeleteCharacterInput struct {
	UserID      string      `json:"user_id"`
	CharacterID pgtype.UUID `json:"character_id"`
}

func ToShortCharacterModel(c db.Character) CharacterModel {
	return CharacterModel{
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

func ToFullCharacterModel(d CharacterDBData) CharacterModel {
	m := ToShortCharacterModel(d.Character)

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
		skills := make([]SkillModel, len(d.Skills))
		for i, s := range d.Skills {
			skills[i] = ToSkillModel(s)
		}
		m.Skills = skills
	}

	if d.Backstory != nil {
		m.Backstory = ToBackstoryModel(*d.Backstory, d.BackstoryItems)
	}

	if d.Finances != nil {
		var creditRating *SkillModel
		if d.Finances.CreditRatingSkillID.Valid {
			for _, skill := range m.Skills {
				if skill.ID == d.Finances.CreditRatingSkillID {
					creditRating = &skill
					break
				}
			}
		}
		m.Finances = ToFinancesModel(*d.Finances, creditRating)
	}

	return m
}
