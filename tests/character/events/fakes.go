package tests

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
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

	Characters      []model.CharacterModel
	Character       model.CharacterModel
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

func (f *FakeCharacterService) GetAllCharacters(context.Context, string) ([]model.CharacterModel, error) {
	return f.Characters, f.Err
}

func (f *FakeCharacterService) GetCharacter(context.Context, model.GetCharacterInput) (model.CharacterModel, error) {
	return f.Character, f.Err
}

func (f *FakeCharacterService) CreateCharacter(context.Context, db.CreateCharacterParams) (model.CharacterModel, error) {
	return f.Character, f.Err
}

func (f *FakeCharacterService) UpdateCharacter(context.Context, db.UpdateCharacterParams) (model.CharacterModel, error) {
	return f.Character, f.Err
}

func (f *FakeCharacterService) DeleteCharacter(context.Context, db.DeleteCharacterParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetHealth(context.Context, db.GetHealthStateParams) (db.HealthState, error) {
	return f.Health, f.Err
}

func (f *FakeCharacterService) UpsertHealth(context.Context, db.UpsertHealthStateParams) (db.HealthState, error) {
	return f.Health, f.Err
}

func (f *FakeCharacterService) DeleteHealth(context.Context, db.DeleteHealthStateParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetSanity(context.Context, db.GetSanityStateParams) (db.SanityState, error) {
	return f.Sanity, f.Err
}

func (f *FakeCharacterService) UpsertSanity(context.Context, db.UpsertSanityStateParams) (db.SanityState, error) {
	return f.Sanity, f.Err
}

func (f *FakeCharacterService) DeleteSanity(context.Context, db.DeleteSanityStateParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetMagic(context.Context, db.GetMagicStateParams) (db.MagicState, error) {
	return f.Magic, f.Err
}

func (f *FakeCharacterService) UpsertMagic(context.Context, db.UpsertMagicStateParams) (db.MagicState, error) {
	return f.Magic, f.Err
}

func (f *FakeCharacterService) DeleteMagic(context.Context, db.DeleteMagicStateParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetLuck(context.Context, db.GetLuckStateParams) (db.LuckState, error) {
	return f.Luck, f.Err
}

func (f *FakeCharacterService) UpsertLuck(context.Context, db.UpsertLuckStateParams) (db.LuckState, error) {
	return f.Luck, f.Err
}

func (f *FakeCharacterService) DeleteLuck(context.Context, db.DeleteLuckStateParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetFinances(context.Context, db.GetFinancesParams) (db.Finance, error) {
	return f.Finances, f.Err
}

func (f *FakeCharacterService) UpsertFinances(context.Context, db.UpsertFinancesParams) (db.Finance, error) {
	return f.Finances, f.Err
}

func (f *FakeCharacterService) DeleteFinances(context.Context, db.DeleteFinancesParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetBackstory(context.Context, db.GetBackstoryByCharacterParams) (model.BackstoryModel, error) {
	return f.Backstory, f.Err
}

func (f *FakeCharacterService) UpsertBackstory(context.Context, db.UpsertBackstoryParams) (model.BackstoryModel, error) {
	return f.Backstory, f.Err
}

func (f *FakeCharacterService) DeleteBackstory(context.Context, db.DeleteBackstoryParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetBackstoryItems(context.Context, db.GetBackstoryItemsParams) ([]model.BackstoryItemModel, error) {
	return f.BackstoryItems, f.Err
}

func (f *FakeCharacterService) GetBackstoryItem(context.Context, db.GetBackstoryItemParams) (model.BackstoryItemModel, error) {
	return f.BackstoryItem, f.Err
}

func (f *FakeCharacterService) CreateBackstoryItem(context.Context, db.CreateBackstoryItemParams) (model.BackstoryItemModel, error) {
	return f.BackstoryItem, f.Err
}

func (f *FakeCharacterService) UpdateBackstoryItem(context.Context, db.UpdateBackstoryItemParams) (model.BackstoryItemModel, error) {
	return f.BackstoryItem, f.Err
}

func (f *FakeCharacterService) DeleteBackstoryItem(context.Context, db.DeleteBackstoryItemParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetSkills(context.Context, db.GetCharacterSkillsParams) ([]model.SkillModel, error) {
	return f.Skills, f.Err
}

func (f *FakeCharacterService) GetSkill(context.Context, db.GetCharacterSkillParams) (model.SkillModel, error) {
	return f.Skill, f.Err
}

func (f *FakeCharacterService) CreateSkill(context.Context, db.CreateCharacterSkillParams) (model.SkillModel, error) {
	return f.Skill, f.Err
}

func (f *FakeCharacterService) UpdateSkill(context.Context, db.UpdateCharacterSkillParams) (model.SkillModel, error) {
	return f.Skill, f.Err
}

func (f *FakeCharacterService) DeleteSkill(context.Context, db.DeleteCharacterSkillParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetDerivedStats(context.Context, db.GetDerivedStatsParams) (db.DerivedStat, error) {
	return f.DerivedStats, f.Err
}

func (f *FakeCharacterService) UpsertDerivedStats(context.Context, db.UpsertDerivedStatsParams) (db.DerivedStat, error) {
	return f.DerivedStats, f.Err
}

func (f *FakeCharacterService) DeleteDerivedStats(context.Context, db.DeleteDerivedStatsParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetCharacteristics(context.Context, db.GetCharacteristicsParams) (db.Characteristic, error) {
	return f.Characteristics, f.Err
}

func (f *FakeCharacterService) UpsertCharacteristics(context.Context, db.UpsertCharacteristicsParams) (db.Characteristic, error) {
	return f.Characteristics, f.Err
}

func (f *FakeCharacterService) DeleteCharacteristics(context.Context, db.DeleteCharacteristicsParams) error {
	return f.Err
}

func (f *FakeCharacterService) GetNotes(context.Context, db.GetNotesParams) ([]db.Note, error) {
	return f.Notes, f.Err
}

func (f *FakeCharacterService) GetNote(context.Context, db.GetNoteParams) (db.Note, error) {
	return f.Note, f.Err
}

func (f *FakeCharacterService) CreateNote(context.Context, db.CreateNoteParams) (db.Note, error) {
	return f.Note, f.Err
}

func (f *FakeCharacterService) UpdateNote(context.Context, db.UpdateNoteParams) (db.Note, error) {
	return f.Note, f.Err
}

func (f *FakeCharacterService) DeleteNote(context.Context, db.DeleteNoteParams) error {
	return f.Err
}
