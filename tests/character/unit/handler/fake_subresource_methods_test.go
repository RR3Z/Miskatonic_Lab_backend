package tests

import (
	"context"

	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

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
