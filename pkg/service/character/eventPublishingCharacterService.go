package character

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

// Implements ICharacter
type EventPublishingCharacterService struct {
	next      ICharacter
	publisher events.EventPublisher
}

func NewEventPublishingCharacterService(next ICharacter, publisher events.EventPublisher) *EventPublishingCharacterService {
	return &EventPublishingCharacterService{
		next:      next,
		publisher: publisher,
	}
}

// Characters
func (s *EventPublishingCharacterService) GetAllCharacters(ctx context.Context, userID string) ([]model.CharacterModel, error) {
	characters, err := s.next.GetAllCharacters(ctx, userID)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharactersListFailed{
			UserID: userID,
			Err:    err,
		})

		return nil, err
	}

	s.publisher.Publish(ctx, characterEvents.CharactersListSucceeded{
		UserID: userID,
		Count:  len(characters),
	})

	return characters, nil
}

func (s *EventPublishingCharacterService) GetCharacter(ctx context.Context, input model.GetCharacterInput) (model.CharacterModel, error) {
	character, err := s.next.GetCharacter(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return model.CharacterModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
		Name:        character.Name,
	})

	return character, nil
}

func (s *EventPublishingCharacterService) CreateCharacter(ctx context.Context, input db.CreateCharacterParams) (model.CharacterModel, error) {
	character, err := s.next.CreateCharacter(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCreateFailed{
			UserID: input.UserID,
			Err:    err,
		})
		return model.CharacterModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCreateSucceeded{
		UserID:      input.UserID,
		CharacterID: character.ID.String(),
		Name:        character.Name,
	})

	return character, nil
}

func (s *EventPublishingCharacterService) UpdateCharacter(ctx context.Context, input db.UpdateCharacterParams) (model.CharacterModel, error) {
	character, err := s.next.UpdateCharacter(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterUpdateFailed{
			UserID:      input.UserID,
			CharacterID: input.ID.String(),
			Err:         err,
		})
		return model.CharacterModel{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterUpdateSucceeded{
		UserID:      input.UserID,
		CharacterID: character.ID.String(),
		Name:        character.Name,
	})

	return character, nil
}

func (s *EventPublishingCharacterService) DeleteCharacter(ctx context.Context, input db.DeleteCharacterParams) error {
	if err := s.next.DeleteCharacter(ctx, input); err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.ID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.ID.String(),
	})

	return nil
}

// Health
func (s *EventPublishingCharacterService) GetHealth(ctx context.Context, input db.GetHealthStateParams) (db.HealthState, error) {
	health, err := s.next.GetHealth(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterHealthGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.HealthState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterHealthGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return health, nil
}

func (s *EventPublishingCharacterService) UpsertHealth(ctx context.Context, input db.UpsertHealthStateParams) (db.HealthState, error) {
	health, err := s.next.UpsertHealth(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterHealthUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.HealthState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterHealthUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return health, nil
}

func (s *EventPublishingCharacterService) DeleteHealth(ctx context.Context, input db.DeleteHealthStateParams) error {
	err := s.next.DeleteHealth(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterHealthDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterHealthDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Sanity
func (s *EventPublishingCharacterService) GetSanity(ctx context.Context, input db.GetSanityStateParams) (db.SanityState, error) {
	sanity, err := s.next.GetSanity(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSanityGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.SanityState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSanityGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return sanity, nil
}

func (s *EventPublishingCharacterService) UpsertSanity(ctx context.Context, input db.UpsertSanityStateParams) (db.SanityState, error) {
	sanity, err := s.next.UpsertSanity(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSanityUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.SanityState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSanityUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return sanity, nil
}

func (s *EventPublishingCharacterService) DeleteSanity(ctx context.Context, input db.DeleteSanityStateParams) error {
	err := s.next.DeleteSanity(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterSanityDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterSanityDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Magic
func (s *EventPublishingCharacterService) GetMagic(ctx context.Context, input db.GetMagicStateParams) (db.MagicState, error) {
	magic, err := s.next.GetMagic(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterMagicGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.MagicState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterMagicGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return magic, nil
}

func (s *EventPublishingCharacterService) UpsertMagic(ctx context.Context, input db.UpsertMagicStateParams) (db.MagicState, error) {
	magic, err := s.next.UpsertMagic(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterMagicUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.MagicState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterMagicUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return magic, nil
}

func (s *EventPublishingCharacterService) DeleteMagic(ctx context.Context, input db.DeleteMagicStateParams) error {
	err := s.next.DeleteMagic(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterMagicDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterMagicDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Luck
func (s *EventPublishingCharacterService) GetLuck(ctx context.Context, input db.GetLuckStateParams) (db.LuckState, error) {
	luck, err := s.next.GetLuck(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterLuckGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.LuckState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterLuckGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return luck, nil
}

func (s *EventPublishingCharacterService) UpsertLuck(ctx context.Context, input db.UpsertLuckStateParams) (db.LuckState, error) {
	luck, err := s.next.UpsertLuck(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterLuckUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.LuckState{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterLuckUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return luck, nil
}

func (s *EventPublishingCharacterService) DeleteLuck(ctx context.Context, input db.DeleteLuckStateParams) error {
	err := s.next.DeleteLuck(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterLuckDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterLuckDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Characteristics
func (s *EventPublishingCharacterService) GetCharacteristics(ctx context.Context, input db.GetCharacteristicsParams) (db.Characteristic, error) {
	characteristics, err := s.next.GetCharacteristics(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsGetFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Characteristic{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsGetSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return characteristics, nil
}

func (s *EventPublishingCharacterService) UpsertCharacteristics(ctx context.Context, input db.UpsertCharacteristicsParams) (db.Characteristic, error) {
	characteristics, err := s.next.UpsertCharacteristics(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsUpsertFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return db.Characteristic{}, err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsUpsertSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return characteristics, nil
}

func (s *EventPublishingCharacterService) DeleteCharacteristics(ctx context.Context, input db.DeleteCharacteristicsParams) error {
	err := s.next.DeleteCharacteristics(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsDeleteFailed{
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return err
	}

	s.publisher.Publish(ctx, characterEvents.CharacterCharacteristicsDeleteSucceeded{
		UserID:      input.UserID,
		CharacterID: input.CharacterID.String(),
	})

	return nil
}

// Notes
func (s *EventPublishingCharacterService) GetNotes(ctx context.Context, input db.GetNotesParams) ([]db.Note, error) {
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

func (s *EventPublishingCharacterService) GetNote(ctx context.Context, input db.GetNoteParams) (db.Note, error) {
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

func (s *EventPublishingCharacterService) CreateNote(ctx context.Context, input db.CreateNoteParams) (db.Note, error) {
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

func (s *EventPublishingCharacterService) UpdateNote(ctx context.Context, input db.UpdateNoteParams) (db.Note, error) {
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

func (s *EventPublishingCharacterService) DeleteNote(ctx context.Context, input db.DeleteNoteParams) error {
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
