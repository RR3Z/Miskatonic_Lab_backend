package character

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

type ICharacter interface {
	GetAllCharacters(ctx context.Context, userID string) ([]model.CharacterModel, error)
	GetCharacter(ctx context.Context, input model.GetCharacterInput) (model.CharacterModel, error)
	CreateCharacter(ctx context.Context, input db.CreateCharacterParams) (model.CharacterModel, error)
	UpdateCharacter(ctx context.Context, input db.UpdateCharacterParams) (model.CharacterModel, error)
	DeleteCharacter(ctx context.Context, input db.DeleteCharacterParams) error

	GetHealth(ctx context.Context, input db.GetHealthStateParams) (db.HealthState, error)
	UpsertHealth(ctx context.Context, input db.UpsertHealthStateParams) (db.HealthState, error)
	DeleteHealth(ctx context.Context, input db.DeleteHealthStateParams) error

	GetSanity(ctx context.Context, input db.GetSanityStateParams) (db.SanityState, error)
	UpsertSanity(ctx context.Context, input db.UpsertSanityStateParams) (db.SanityState, error)
	DeleteSanity(ctx context.Context, input db.DeleteSanityStateParams) error

	GetCharacteristics(ctx context.Context, input db.GetCharacteristicsParams) (db.Characteristic, error)
	UpsertCharacteristics(ctx context.Context, input db.UpsertCharacteristicsParams) (db.Characteristic, error)
	DeleteCharacteristics(ctx context.Context, input db.DeleteCharacteristicsParams) error

	GetNotes(ctx context.Context, input db.GetNotesParams) ([]db.Note, error)
	GetNote(ctx context.Context, input db.GetNoteParams) (db.Note, error)
	CreateNote(ctx context.Context, input db.CreateNoteParams) (db.Note, error)
	UpdateNote(ctx context.Context, input db.UpdateNoteParams) (db.Note, error)
	DeleteNote(ctx context.Context, input db.DeleteNoteParams) error
}
