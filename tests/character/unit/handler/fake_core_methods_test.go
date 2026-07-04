package tests

import (
	"context"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
)

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
