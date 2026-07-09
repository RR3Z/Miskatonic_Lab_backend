package tests

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	derivedStatsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	luckDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	sanityDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

type FakeEventPublisher struct {
	Events []events.Event
}

func (f *FakeEventPublisher) Publish(_ context.Context, event events.Event) {
	f.Events = append(f.Events, event)
}

type FakeCharacterService struct {
	Err error

	Characters      []characterDTO.CharacterSummaryModel
	Character       characterDTO.CharacterModel
	Health          db.HealthState
	Sanity          db.SanityState
	Magic           db.MagicState
	Luck            db.LuckState
	Finances        db.Finance
	Backstory       backstoriesDTO.BackstoryModel
	BackstoryItems  []backstoriesDTO.BackstoryItemModel
	BackstoryItem   backstoriesDTO.BackstoryItemModel
	Skills          []skillsDTO.SkillModel
	Skill           skillsDTO.SkillModel
	DerivedStats    db.DerivedStat
	Characteristics db.Characteristic
	Notes           []db.Note
	Note            db.Note
}

func (f *FakeCharacterService) GetAllCharacters(context.Context, string) ([]characterDTO.CharacterSummaryModel, error) {
	return f.Characters, f.Err
}

func (f *FakeCharacterService) GetCharacter(context.Context, characterDTO.GetCharacterInput) (characterDTO.CharacterModel, error) {
	return f.Character, f.Err
}

func (f *FakeCharacterService) CreateCharacter(context.Context, characterDTO.CreateCharacterInput) (characterDTO.CharacterShortModel, error) {
	return f.Character.CharacterShortModel, f.Err
}

func (f *FakeCharacterService) UpdateCharacter(context.Context, characterDTO.UpdateCharacterInput) (characterDTO.CharacterShortModel, error) {
	return f.Character.CharacterShortModel, f.Err
}

func (f *FakeCharacterService) DeleteCharacter(context.Context, characterDTO.DeleteCharacterInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetHealth(context.Context, healthDTO.GetHealthInput) (db.HealthState, error) {
	return f.Health, f.Err
}

func (f *FakeCharacterService) UpsertHealth(context.Context, healthDTO.UpsertHealthInput) (db.HealthState, error) {
	return f.Health, f.Err
}

func (f *FakeCharacterService) DeleteHealth(context.Context, healthDTO.DeleteHealthInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetSanity(context.Context, sanityDTO.GetSanityInput) (db.SanityState, error) {
	return f.Sanity, f.Err
}

func (f *FakeCharacterService) UpsertSanity(context.Context, sanityDTO.UpsertSanityInput) (db.SanityState, error) {
	return f.Sanity, f.Err
}

func (f *FakeCharacterService) DeleteSanity(context.Context, sanityDTO.DeleteSanityInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetMagic(context.Context, magicDTO.GetMagicInput) (db.MagicState, error) {
	return f.Magic, f.Err
}

func (f *FakeCharacterService) UpsertMagic(context.Context, magicDTO.UpsertMagicInput) (db.MagicState, error) {
	return f.Magic, f.Err
}

func (f *FakeCharacterService) DeleteMagic(context.Context, magicDTO.DeleteMagicInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetLuck(context.Context, luckDTO.GetLuckInput) (db.LuckState, error) {
	return f.Luck, f.Err
}

func (f *FakeCharacterService) UpsertLuck(context.Context, luckDTO.UpsertLuckInput) (db.LuckState, error) {
	return f.Luck, f.Err
}

func (f *FakeCharacterService) DeleteLuck(context.Context, luckDTO.DeleteLuckInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetFinances(context.Context, financesDTO.GetFinancesInput) (db.Finance, error) {
	return f.Finances, f.Err
}

func (f *FakeCharacterService) UpsertFinances(context.Context, financesDTO.UpsertFinancesInput) (db.Finance, error) {
	return f.Finances, f.Err
}

func (f *FakeCharacterService) DeleteFinances(context.Context, financesDTO.DeleteFinancesInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetBackstory(context.Context, backstoriesDTO.GetBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	return f.Backstory, f.Err
}

func (f *FakeCharacterService) UpsertBackstory(context.Context, backstoriesDTO.UpsertBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	return f.Backstory, f.Err
}

func (f *FakeCharacterService) DeleteBackstory(context.Context, backstoriesDTO.DeleteBackstoryInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetBackstoryItems(context.Context, backstoriesDTO.GetBackstoryItemsInput) ([]backstoriesDTO.BackstoryItemModel, error) {
	return f.BackstoryItems, f.Err
}

func (f *FakeCharacterService) GetBackstoryItem(context.Context, backstoriesDTO.GetBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	return f.BackstoryItem, f.Err
}

func (f *FakeCharacterService) CreateBackstoryItem(context.Context, backstoriesDTO.CreateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	return f.BackstoryItem, f.Err
}

func (f *FakeCharacterService) UpdateBackstoryItem(context.Context, backstoriesDTO.UpdateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	return f.BackstoryItem, f.Err
}

func (f *FakeCharacterService) DeleteBackstoryItem(context.Context, backstoriesDTO.DeleteBackstoryItemInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetSkills(context.Context, skillsDTO.GetSkillsInput) ([]skillsDTO.SkillModel, error) {
	return f.Skills, f.Err
}

func (f *FakeCharacterService) GetSkill(context.Context, skillsDTO.GetSkillInput) (skillsDTO.SkillModel, error) {
	return f.Skill, f.Err
}

func (f *FakeCharacterService) CreateSkill(context.Context, skillsDTO.CreateSkillInput) (skillsDTO.SkillModel, error) {
	return f.Skill, f.Err
}

func (f *FakeCharacterService) UpdateSkill(context.Context, skillsDTO.UpdateSkillInput) (skillsDTO.SkillModel, error) {
	return f.Skill, f.Err
}

func (f *FakeCharacterService) DeleteSkill(context.Context, skillsDTO.DeleteSkillInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetDerivedStats(context.Context, derivedStatsDTO.GetDerivedStatsInput) (db.DerivedStat, error) {
	return f.DerivedStats, f.Err
}

func (f *FakeCharacterService) UpsertDerivedStats(context.Context, derivedStatsDTO.UpsertDerivedStatsInput) (db.DerivedStat, error) {
	return f.DerivedStats, f.Err
}

func (f *FakeCharacterService) DeleteDerivedStats(context.Context, derivedStatsDTO.DeleteDerivedStatsInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetCharacteristics(context.Context, characteristicsDTO.GetCharacteristicsInput) (db.Characteristic, error) {
	return f.Characteristics, f.Err
}

func (f *FakeCharacterService) UpsertCharacteristics(context.Context, characteristicsDTO.UpsertCharacteristicsInput) (db.Characteristic, error) {
	return f.Characteristics, f.Err
}

func (f *FakeCharacterService) DeleteCharacteristics(context.Context, characteristicsDTO.DeleteCharacteristicsInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetNotes(context.Context, notesDTO.GetNotesInput) ([]db.Note, error) {
	return f.Notes, f.Err
}

func (f *FakeCharacterService) GetNote(context.Context, notesDTO.GetNoteInput) (db.Note, error) {
	return f.Note, f.Err
}

func (f *FakeCharacterService) CreateNote(context.Context, notesDTO.CreateNoteInput) (db.Note, error) {
	return f.Note, f.Err
}

func (f *FakeCharacterService) UpdateNote(context.Context, notesDTO.UpdateNoteInput) (db.Note, error) {
	return f.Note, f.Err
}

func (f *FakeCharacterService) DeleteNote(context.Context, notesDTO.DeleteNoteInput) error {
	return f.Err
}
