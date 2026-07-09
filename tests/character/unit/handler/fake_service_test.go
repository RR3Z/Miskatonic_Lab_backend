package tests

import (
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	backstoriesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/backstories"
	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	skillsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/skills"
)

type fakeCharacterHandlerService struct {
	err error

	characters []characterDTO.CharacterSummaryModel
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
