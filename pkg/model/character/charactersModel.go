package characterDTO

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/inventory"
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
	Skills          []db.Skill
	Backstory       *db.Backstory
	BackstoryItems  []db.BackstoryItem
	Finances        *db.Finance
	InventoryItems  []db.CharacterInventoryItem
	Notes           []db.Note
}

type CharacterShortModel struct {
	ID     pgtype.UUID `json:"id"`
	UserID string      `json:"user_id"`

	Name        string  `json:"name"`
	Occupation  *string `json:"occupation"`
	Age         *int16  `json:"age"`
	Sex         *string `json:"sex"`
	Residence   *string `json:"residence"`
	Birthplace  *string `json:"birthplace"`
	PortraitUrl *string `json:"portrait_url"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type CharacterSummaryModel struct {
	ID pgtype.UUID `json:"id"`

	Name        string  `json:"name"`
	Occupation  *string `json:"occupation"`
	Age         *int16  `json:"age"`
	Sex         *string `json:"sex"`
	Residence   *string `json:"residence"`
	PortraitUrl *string `json:"portrait_url"`

	HP struct {
		Current int16 `json:"current_hp"`
		Max     int16 `json:"max_hp"`
	} `json:"hp"`

	MP struct {
		Current int16 `json:"current_mp"`
		Max     int16 `json:"max_mp"`
	} `json:"mp"`

	Sanity struct {
		Current int16 `json:"current_sanity"`
		Max     int16 `json:"max_sanity"`
	} `json:"sanity"`

	Luck struct {
		Current  int16 `json:"current_luck"`
		Starting int16 `json:"starting_luck"`
	} `json:"luck"`
}

type CharacterModel struct {
	CharacterShortModel

	Skills          []skillsDTO.SkillModel                  `json:"skills"`
	Characteristics characteristicsDTO.CharacteristicsModel `json:"characteristics"`
	DerivedStats    derivedStatsDTO.DerivedStatsModel       `json:"derived_stats"`
	HP              healthDTO.HealthModel                   `json:"hp"`
	MP              magicDTO.MagicModel                     `json:"mp"`
	Sanity          sanityDTO.SanityModel                   `json:"sanity"`
	Luck            luckDTO.LuckModel                       `json:"luck"`
	Backstory       backstoriesDTO.BackstoryModel           `json:"backstory"`
	Finances        financesDTO.FinancesModel               `json:"finances"`
	Inventory       []inventoryDTO.InventoryItemModel       `json:"inventory"`
	Notes           []notesDTO.NoteModel                    `json:"notes"`
}

func ToCharacterShortModel(c db.Character) CharacterShortModel {
	return CharacterShortModel{
		ID:         c.ID,
		UserID:     c.UserID,
		Name:       c.Name,
		Occupation: c.Occupation,
		Age:        c.Age,
		Sex:        c.Sex,
		Residence:  c.Residence,
		Birthplace: c.Birthplace,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

func ToCharacterSummaryModel(row db.GetAllUserCharacterCardsRow) CharacterSummaryModel {
	m := CharacterSummaryModel{
		ID:         row.ID,
		Name:       row.Name,
		Occupation: row.Occupation,
		Age:        row.Age,
		Sex:        row.Sex,
		Residence:  row.Residence,
	}
	m.HP.Current = row.CurrentHp
	m.HP.Max = row.MaxHp
	m.MP.Current = row.CurrentMp
	m.MP.Max = row.MaxMp
	m.Sanity.Current = row.CurrentSanity
	m.Sanity.Max = row.MaxSanity
	m.Luck.Current = row.CurrentLuck
	m.Luck.Starting = row.StartingLuck
	return m
}

func ToCharacterModel(d CharacterDBData) CharacterModel {
	m := CharacterModel{
		CharacterShortModel: ToCharacterShortModel(d.Character),
	}

	if d.Characteristics.ID.Valid {
		m.Characteristics = characteristicsDTO.ToCharacteristicsModel(d.Characteristics)
	}
	if d.DerivedStats.ID.Valid {
		m.DerivedStats = derivedStatsDTO.ToDerivedStatsModel(d.DerivedStats)
	}
	if d.HP.ID.Valid {
		m.HP = healthDTO.ToHealthModel(d.HP)
	}
	if d.MP.ID.Valid {
		m.MP = magicDTO.ToMagicModel(d.MP)
	}
	if d.Sanity.ID.Valid {
		m.Sanity = sanityDTO.ToSanityModel(d.Sanity)
	}
	if d.Luck.ID.Valid {
		m.Luck = luckDTO.ToLuckModel(d.Luck)
	}

	m.Notes = notesDTO.ToNoteModels(d.Notes)
	m.Inventory = inventoryDTO.ToInventoryItemModels(d.InventoryItems)

	if len(d.Skills) > 0 {
		skillModels := make([]skillsDTO.SkillModel, len(d.Skills))
		for i, s := range d.Skills {
			skillModels[i] = skillsDTO.ToSkillModel(s)
		}
		m.Skills = skillModels
	}

	if d.Backstory != nil {
		m.Backstory = backstoriesDTO.ToBackstoryModel(*d.Backstory, d.BackstoryItems)
	}

	if d.Finances != nil {
		m.Finances = financesDTO.ToFinancesModel(*d.Finances)
	}

	return m
}
