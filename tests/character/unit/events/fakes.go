package tests

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	characterModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
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

	Characters      []characterModel.CharacterShortModel
	Character       characterModel.CharacterModel
	Health          db.HealthState
	Sanity          db.SanityState
	Magic           db.MagicState
	Luck            db.LuckState
	Finances        db.Finance
	Backstory       model.BackstoryModel
	BackstoryItems  []model.BackstoryItemModel
	BackstoryItem   model.BackstoryItemModel
	Skills          []model.SkillModel
	Skill           model.SkillModel
	DerivedStats    db.DerivedStat
	Characteristics db.Characteristic
	Notes           []db.Note
	Note            db.Note
}

func (f *FakeCharacterService) GetAllCharacters(context.Context, string) ([]characterModel.CharacterShortModel, error) {
	return f.Characters, f.Err
}

func (f *FakeCharacterService) GetCharacter(context.Context, characterModel.GetCharacterInput) (characterModel.CharacterModel, error) {
	return f.Character, f.Err
}

func (f *FakeCharacterService) CreateCharacter(context.Context, characterModel.CreateCharacterInput) (characterModel.CharacterShortModel, error) {
	return f.Character.CharacterShortModel, f.Err
}

func (f *FakeCharacterService) UpdateCharacter(context.Context, characterModel.UpdateCharacterInput) (characterModel.CharacterShortModel, error) {
	return f.Character.CharacterShortModel, f.Err
}

func (f *FakeCharacterService) DeleteCharacter(context.Context, characterModel.DeleteCharacterInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetHealth(context.Context, characterModel.GetHealthInput) (db.HealthState, error) {
	return f.Health, f.Err
}

func (f *FakeCharacterService) UpsertHealth(context.Context, characterModel.UpsertHealthInput) (db.HealthState, error) {
	return f.Health, f.Err
}

func (f *FakeCharacterService) DeleteHealth(context.Context, characterModel.DeleteHealthInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetSanity(context.Context, characterModel.GetSanityInput) (db.SanityState, error) {
	return f.Sanity, f.Err
}

func (f *FakeCharacterService) UpsertSanity(context.Context, characterModel.UpsertSanityInput) (db.SanityState, error) {
	return f.Sanity, f.Err
}

func (f *FakeCharacterService) DeleteSanity(context.Context, characterModel.DeleteSanityInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetMagic(context.Context, characterModel.GetMagicInput) (db.MagicState, error) {
	return f.Magic, f.Err
}

func (f *FakeCharacterService) UpsertMagic(context.Context, characterModel.UpsertMagicInput) (db.MagicState, error) {
	return f.Magic, f.Err
}

func (f *FakeCharacterService) DeleteMagic(context.Context, characterModel.DeleteMagicInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetLuck(context.Context, characterModel.GetLuckInput) (db.LuckState, error) {
	return f.Luck, f.Err
}

func (f *FakeCharacterService) UpsertLuck(context.Context, characterModel.UpsertLuckInput) (db.LuckState, error) {
	return f.Luck, f.Err
}

func (f *FakeCharacterService) DeleteLuck(context.Context, characterModel.DeleteLuckInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetFinances(context.Context, characterModel.GetFinancesInput) (db.Finance, error) {
	return f.Finances, f.Err
}

func (f *FakeCharacterService) UpsertFinances(context.Context, characterModel.UpsertFinancesInput) (db.Finance, error) {
	return f.Finances, f.Err
}

func (f *FakeCharacterService) DeleteFinances(context.Context, characterModel.DeleteFinancesInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetBackstory(context.Context, characterModel.GetBackstoryInput) (model.BackstoryModel, error) {
	return f.Backstory, f.Err
}

func (f *FakeCharacterService) UpsertBackstory(context.Context, characterModel.UpsertBackstoryInput) (model.BackstoryModel, error) {
	return f.Backstory, f.Err
}

func (f *FakeCharacterService) DeleteBackstory(context.Context, characterModel.DeleteBackstoryInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetBackstoryItems(context.Context, characterModel.GetBackstoryItemsInput) ([]model.BackstoryItemModel, error) {
	return f.BackstoryItems, f.Err
}

func (f *FakeCharacterService) GetBackstoryItem(context.Context, characterModel.GetBackstoryItemInput) (model.BackstoryItemModel, error) {
	return f.BackstoryItem, f.Err
}

func (f *FakeCharacterService) CreateBackstoryItem(context.Context, characterModel.CreateBackstoryItemInput) (model.BackstoryItemModel, error) {
	return f.BackstoryItem, f.Err
}

func (f *FakeCharacterService) UpdateBackstoryItem(context.Context, characterModel.UpdateBackstoryItemInput) (model.BackstoryItemModel, error) {
	return f.BackstoryItem, f.Err
}

func (f *FakeCharacterService) DeleteBackstoryItem(context.Context, characterModel.DeleteBackstoryItemInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetSkills(context.Context, characterModel.GetSkillsInput) ([]model.SkillModel, error) {
	return f.Skills, f.Err
}

func (f *FakeCharacterService) GetSkill(context.Context, characterModel.GetSkillInput) (model.SkillModel, error) {
	return f.Skill, f.Err
}

func (f *FakeCharacterService) CreateSkill(context.Context, characterModel.CreateSkillInput) (model.SkillModel, error) {
	return f.Skill, f.Err
}

func (f *FakeCharacterService) UpdateSkill(context.Context, characterModel.UpdateSkillInput) (model.SkillModel, error) {
	return f.Skill, f.Err
}

func (f *FakeCharacterService) DeleteSkill(context.Context, characterModel.DeleteSkillInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetDerivedStats(context.Context, characterModel.GetDerivedStatsInput) (db.DerivedStat, error) {
	return f.DerivedStats, f.Err
}

func (f *FakeCharacterService) UpsertDerivedStats(context.Context, characterModel.UpsertDerivedStatsInput) (db.DerivedStat, error) {
	return f.DerivedStats, f.Err
}

func (f *FakeCharacterService) DeleteDerivedStats(context.Context, characterModel.DeleteDerivedStatsInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetCharacteristics(context.Context, characterModel.GetCharacteristicsInput) (db.Characteristic, error) {
	return f.Characteristics, f.Err
}

func (f *FakeCharacterService) UpsertCharacteristics(context.Context, characterModel.UpsertCharacteristicsInput) (db.Characteristic, error) {
	return f.Characteristics, f.Err
}

func (f *FakeCharacterService) DeleteCharacteristics(context.Context, characterModel.DeleteCharacteristicsInput) error {
	return f.Err
}

func (f *FakeCharacterService) GetNotes(context.Context, characterModel.GetNotesInput) ([]db.Note, error) {
	return f.Notes, f.Err
}

func (f *FakeCharacterService) GetNote(context.Context, characterModel.GetNoteInput) (db.Note, error) {
	return f.Note, f.Err
}

func (f *FakeCharacterService) CreateNote(context.Context, characterModel.CreateNoteInput) (db.Note, error) {
	return f.Note, f.Err
}

func (f *FakeCharacterService) UpdateNote(context.Context, characterModel.UpdateNoteInput) (db.Note, error) {
	return f.Note, f.Err
}

func (f *FakeCharacterService) DeleteNote(context.Context, characterModel.DeleteNoteInput) error {
	return f.Err
}
