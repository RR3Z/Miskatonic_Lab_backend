package character

import (
	"context"

	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

// Notes
func (s *CharacterService) GetNotes(ctx context.Context, input notesDTO.GetNotesInput) ([]db.Note, error) {
	notes, err := s.repos.Queries.GetNotes(ctx, db.GetNotesParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *CharacterService) GetNote(ctx context.Context, input notesDTO.GetNoteInput) (db.Note, error) {
	note, err := s.repos.Queries.GetNote(ctx, db.GetNoteParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		NoteID:      input.NoteID,
	})
	if err != nil {
		return db.Note{}, err
	}

	return note, nil
}

func (s *CharacterService) CreateNote(ctx context.Context, input notesDTO.CreateNoteInput) (db.Note, error) {
	note, err := s.repos.Queries.CreateNote(ctx, db.CreateNoteParams{
		Title:       input.Title,
		Body:        input.Body,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		return db.Note{}, err
	}

	return note, nil
}

func (s *CharacterService) UpdateNote(ctx context.Context, input notesDTO.UpdateNoteInput) (db.Note, error) {
	note, err := s.repos.Queries.UpdateNote(ctx, db.UpdateNoteParams{
		Title:       input.Title,
		Body:        input.Body,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		NoteID:      input.NoteID,
	})
	if err != nil {
		return db.Note{}, err
	}

	return note, nil
}

func (s *CharacterService) DeleteNote(ctx context.Context, input notesDTO.DeleteNoteInput) error {
	_, err := s.repos.Queries.DeleteNote(ctx, db.DeleteNoteParams{
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
		NoteID:      input.NoteID,
	})
	return err
}
