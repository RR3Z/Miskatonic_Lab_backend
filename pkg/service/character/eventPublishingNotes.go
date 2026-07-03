package character

import (
	"context"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	notesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/notes"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (s *EventPublishingCharacterService) GetNotes(ctx context.Context, input notesDTO.GetNotesInput) ([]db.Note, error) {
	notes, err := s.next.GetNotes(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterNotesListFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return nil, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterNotesListSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		Count:       len(notes),
	})

	return notes, nil
}

func (s *EventPublishingCharacterService) GetNote(ctx context.Context, input notesDTO.GetNoteInput) (db.Note, error) {
	note, err := s.next.GetNote(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterNoteGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			NoteID:      input.NoteID.String(),
			Err:         err,
		})
		return db.Note{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterNoteGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		NoteID:      input.NoteID.String(),
		Title:       note.Title,
	})

	return note, nil
}

func (s *EventPublishingCharacterService) CreateNote(ctx context.Context, input notesDTO.CreateNoteInput) (db.Note, error) {
	note, err := s.next.CreateNote(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterNoteCreateFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Note{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterNoteCreateSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		NoteID:      note.ID.String(),
		Title:       note.Title,
	})

	return note, nil
}

func (s *EventPublishingCharacterService) UpdateNote(ctx context.Context, input notesDTO.UpdateNoteInput) (db.Note, error) {
	note, err := s.next.UpdateNote(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterNoteUpdateFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			NoteID:      input.NoteID.String(),
			Err:         err,
		})
		return db.Note{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterNoteUpdateSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		NoteID:      input.NoteID.String(),
		Title:       note.Title,
	})

	return note, nil
}

func (s *EventPublishingCharacterService) DeleteNote(ctx context.Context, input notesDTO.DeleteNoteInput) error {
	err := s.next.DeleteNote(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterNoteDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			NoteID:      input.NoteID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterNoteDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		NoteID:      input.NoteID.String(),
	})

	return nil
}
