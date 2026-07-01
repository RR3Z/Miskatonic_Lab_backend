package character

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characters"
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
)

type CharacterShortModel = characters.CharacterShortModel
type CharacterModel = characters.CharacterModel
type CharacterDBData = characters.CharacterDBData

type HealthState = db.HealthState
type SanityState = db.SanityState
type MagicState = db.MagicState
type LuckState = db.LuckState
type Characteristic = db.Characteristic
type DerivedStat = db.DerivedStat
type Finance = db.Finance
type Note = db.Note

type BackstoryModel = backstories.BackstoryModel
type BackstoryItemModel = backstories.BackstoryItemModel
type FinancesModel = finances.FinancesModel
type SkillModel = skills.SkillModel

type CharacteristicsModel = characteristics.CharacteristicsModel
type DerivedStatsModel = derivedstats.DerivedStatsModel
type HealthModel = health.HealthModel
type SanityModel = sanity.SanityModel
type MagicModel = magic.MagicModel
type LuckModel = luck.LuckModel
type NoteModel = notes.NoteModel

func ToCharacterShortModel(c db.Character) CharacterShortModel {
	return characters.ToCharacterShortModel(c)
}

func ToCharacterModel(d CharacterDBData) CharacterModel {
	return characters.ToCharacterModel(d)
}

func ToBackstoryModel(b db.Backstory, items []db.BackstoryItem) BackstoryModel {
	return backstories.ToBackstoryModel(b, items)
}

func ToBackstoryItemModel(item db.BackstoryItem) BackstoryItemModel {
	return backstories.ToBackstoryItemModel(item)
}

func ToBackstoryItemModels(items []db.BackstoryItem) []BackstoryItemModel {
	return backstories.ToBackstoryItemModels(items)
}

func ToFinancesModel(f db.Finance, creditRating *skills.SkillModel) FinancesModel {
	return finances.ToFinancesModel(f, creditRating)
}

func ToSkillModel(row db.GetSkillsRow) SkillModel {
	return skills.ToSkillModel(row)
}

func ToCharacterSkillModels(rows []db.GetCharacterSkillsRow) []SkillModel {
	return skills.ToCharacterSkillModels(rows)
}

func ToCharacterSkillModel(row db.GetCharacterSkillsRow) SkillModel {
	return skills.ToCharacterSkillModel(row)
}

func ToSingleCharacterSkillModel(row db.GetCharacterSkillRow) SkillModel {
	return skills.ToSingleCharacterSkillModel(row)
}

func ToCreatedCharacterSkillModel(row db.CreateCharacterSkillRow) SkillModel {
	return skills.ToCreatedCharacterSkillModel(row)
}

func ToUpdatedCharacterSkillModel(row db.UpdateCharacterSkillRow) SkillModel {
	return skills.ToUpdatedCharacterSkillModel(row)
}
