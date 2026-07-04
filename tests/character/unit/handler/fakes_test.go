package tests

import (
	"context"

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

type fakeCharacterHandlerService struct {
	err error

	characters []characterDTO.CharacterShortModel
	character  characterDTO.CharacterModel

	getAllCalls  int
	getAllUserID string

	getCalls int
	getInput characterDTO.GetCharacterInput

	createCalls int
	createInput characterDTO.CreateCharacterInput

	updateCalls int
	updateInput characterDTO.UpdateCharacterInput

	deleteCalls int
	deleteInput characterDTO.DeleteCharacterInput

	getHealthCalls int
	getHealthInput healthDTO.GetHealthInput

	upsertHealthCalls int
	upsertHealthInput healthDTO.UpsertHealthInput

	deleteHealthCalls int
	deleteHealthInput healthDTO.DeleteHealthInput

	getCharacteristicsCalls int
	getCharacteristicsInput characteristicsDTO.GetCharacteristicsInput

	upsertCharacteristicsCalls int
	upsertCharacteristicsInput characteristicsDTO.UpsertCharacteristicsInput

	deleteCharacteristicsCalls int
	deleteCharacteristicsInput characteristicsDTO.DeleteCharacteristicsInput

	getBackstoryItemCalls int
	getBackstoryItemInput backstoriesDTO.GetBackstoryItemInput

	createBackstoryItemCalls int
	createBackstoryItemInput backstoriesDTO.CreateBackstoryItemInput

	deleteSkillCalls int
	deleteSkillInput skillsDTO.DeleteSkillInput

	createNoteCalls int
	createNoteInput notesDTO.CreateNoteInput

	deleteNoteCalls int
	deleteNoteInput notesDTO.DeleteNoteInput
}

func (f *fakeCharacterHandlerService) totalCalls() int {
	return f.getAllCalls + f.getCalls + f.createCalls + f.updateCalls + f.deleteCalls +
		f.getHealthCalls + f.upsertHealthCalls + f.deleteHealthCalls +
		f.getCharacteristicsCalls + f.upsertCharacteristicsCalls + f.deleteCharacteristicsCalls +
		f.getBackstoryItemCalls + f.createBackstoryItemCalls + f.deleteSkillCalls +
		f.createNoteCalls + f.deleteNoteCalls
}

func (f *fakeCharacterHandlerService) GetAllCharacters(_ context.Context, userID string) ([]characterDTO.CharacterShortModel, error) {
	f.getAllCalls++
	f.getAllUserID = userID
	return f.characters, f.err
}

func (f *fakeCharacterHandlerService) GetCharacter(_ context.Context, input characterDTO.GetCharacterInput) (characterDTO.CharacterModel, error) {
	f.getCalls++
	f.getInput = input
	return f.character, f.err
}

func (f *fakeCharacterHandlerService) CreateCharacter(_ context.Context, input characterDTO.CreateCharacterInput) (characterDTO.CharacterShortModel, error) {
	f.createCalls++
	f.createInput = input
	return f.character.CharacterShortModel, f.err
}

func (f *fakeCharacterHandlerService) UpdateCharacter(_ context.Context, input characterDTO.UpdateCharacterInput) (characterDTO.CharacterShortModel, error) {
	f.updateCalls++
	f.updateInput = input
	return f.character.CharacterShortModel, f.err
}

func (f *fakeCharacterHandlerService) DeleteCharacter(_ context.Context, input characterDTO.DeleteCharacterInput) error {
	f.deleteCalls++
	f.deleteInput = input
	return f.err
}

func (f *fakeCharacterHandlerService) GetHealth(_ context.Context, input healthDTO.GetHealthInput) (db.HealthState, error) {
	f.getHealthCalls++
	f.getHealthInput = input
	return db.HealthState{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertHealth(_ context.Context, input healthDTO.UpsertHealthInput) (db.HealthState, error) {
	f.upsertHealthCalls++
	f.upsertHealthInput = input
	return db.HealthState{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteHealth(_ context.Context, input healthDTO.DeleteHealthInput) error {
	f.deleteHealthCalls++
	f.deleteHealthInput = input
	return f.err
}

func (f *fakeCharacterHandlerService) GetSanity(_ context.Context, _ sanityDTO.GetSanityInput) (db.SanityState, error) {
	return db.SanityState{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertSanity(_ context.Context, _ sanityDTO.UpsertSanityInput) (db.SanityState, error) {
	return db.SanityState{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteSanity(_ context.Context, _ sanityDTO.DeleteSanityInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetMagic(_ context.Context, _ magicDTO.GetMagicInput) (db.MagicState, error) {
	return db.MagicState{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertMagic(_ context.Context, _ magicDTO.UpsertMagicInput) (db.MagicState, error) {
	return db.MagicState{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteMagic(_ context.Context, _ magicDTO.DeleteMagicInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetLuck(_ context.Context, _ luckDTO.GetLuckInput) (db.LuckState, error) {
	return db.LuckState{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertLuck(_ context.Context, _ luckDTO.UpsertLuckInput) (db.LuckState, error) {
	return db.LuckState{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteLuck(_ context.Context, _ luckDTO.DeleteLuckInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetFinances(_ context.Context, _ financesDTO.GetFinancesInput) (db.Finance, error) {
	return db.Finance{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertFinances(_ context.Context, _ financesDTO.UpsertFinancesInput) (db.Finance, error) {
	return db.Finance{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteFinances(_ context.Context, _ financesDTO.DeleteFinancesInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetBackstory(_ context.Context, _ backstoriesDTO.GetBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	return backstoriesDTO.BackstoryModel{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertBackstory(_ context.Context, _ backstoriesDTO.UpsertBackstoryInput) (backstoriesDTO.BackstoryModel, error) {
	return backstoriesDTO.BackstoryModel{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteBackstory(_ context.Context, _ backstoriesDTO.DeleteBackstoryInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetBackstoryItems(_ context.Context, _ backstoriesDTO.GetBackstoryItemsInput) ([]backstoriesDTO.BackstoryItemModel, error) {
	return nil, f.err
}

func (f *fakeCharacterHandlerService) GetBackstoryItem(_ context.Context, input backstoriesDTO.GetBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	f.getBackstoryItemCalls++
	f.getBackstoryItemInput = input
	return backstoriesDTO.BackstoryItemModel{}, f.err
}

func (f *fakeCharacterHandlerService) CreateBackstoryItem(_ context.Context, input backstoriesDTO.CreateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	f.createBackstoryItemCalls++
	f.createBackstoryItemInput = input
	return backstoriesDTO.BackstoryItemModel{}, f.err
}

func (f *fakeCharacterHandlerService) UpdateBackstoryItem(_ context.Context, _ backstoriesDTO.UpdateBackstoryItemInput) (backstoriesDTO.BackstoryItemModel, error) {
	return backstoriesDTO.BackstoryItemModel{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteBackstoryItem(_ context.Context, _ backstoriesDTO.DeleteBackstoryItemInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetSkills(_ context.Context, _ skillsDTO.GetSkillsInput) ([]skillsDTO.SkillModel, error) {
	return nil, f.err
}

func (f *fakeCharacterHandlerService) GetSkill(_ context.Context, _ skillsDTO.GetSkillInput) (skillsDTO.SkillModel, error) {
	return skillsDTO.SkillModel{}, f.err
}

func (f *fakeCharacterHandlerService) CreateSkill(_ context.Context, _ skillsDTO.CreateSkillInput) (skillsDTO.SkillModel, error) {
	return skillsDTO.SkillModel{}, f.err
}

func (f *fakeCharacterHandlerService) UpdateSkill(_ context.Context, _ skillsDTO.UpdateSkillInput) (skillsDTO.SkillModel, error) {
	return skillsDTO.SkillModel{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteSkill(_ context.Context, input skillsDTO.DeleteSkillInput) error {
	f.deleteSkillCalls++
	f.deleteSkillInput = input
	return f.err
}

func (f *fakeCharacterHandlerService) GetDerivedStats(_ context.Context, _ derivedStatsDTO.GetDerivedStatsInput) (db.DerivedStat, error) {
	return db.DerivedStat{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertDerivedStats(_ context.Context, _ derivedStatsDTO.UpsertDerivedStatsInput) (db.DerivedStat, error) {
	return db.DerivedStat{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteDerivedStats(_ context.Context, _ derivedStatsDTO.DeleteDerivedStatsInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetCharacteristics(_ context.Context, input characteristicsDTO.GetCharacteristicsInput) (db.Characteristic, error) {
	f.getCharacteristicsCalls++
	f.getCharacteristicsInput = input
	return db.Characteristic{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertCharacteristics(_ context.Context, input characteristicsDTO.UpsertCharacteristicsInput) (db.Characteristic, error) {
	f.upsertCharacteristicsCalls++
	f.upsertCharacteristicsInput = input
	return db.Characteristic{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteCharacteristics(_ context.Context, input characteristicsDTO.DeleteCharacteristicsInput) error {
	f.deleteCharacteristicsCalls++
	f.deleteCharacteristicsInput = input
	return f.err
}

func (f *fakeCharacterHandlerService) GetNotes(_ context.Context, _ notesDTO.GetNotesInput) ([]db.Note, error) {
	return nil, f.err
}

func (f *fakeCharacterHandlerService) GetNote(_ context.Context, _ notesDTO.GetNoteInput) (db.Note, error) {
	return db.Note{}, f.err
}

func (f *fakeCharacterHandlerService) CreateNote(_ context.Context, input notesDTO.CreateNoteInput) (db.Note, error) {
	f.createNoteCalls++
	f.createNoteInput = input
	return db.Note{}, f.err
}

func (f *fakeCharacterHandlerService) UpdateNote(_ context.Context, _ notesDTO.UpdateNoteInput) (db.Note, error) {
	return db.Note{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteNote(_ context.Context, input notesDTO.DeleteNoteInput) error {
	f.deleteNoteCalls++
	f.deleteNoteInput = input
	return f.err
}
