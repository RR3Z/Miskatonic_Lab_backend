package logging

import (
	"context"
	"log/slog"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
)

// Implements EventHandler
type CharacterEventLogger struct {
	logger *slog.Logger
}

func NewCharacterEventLogger(logger *slog.Logger) *CharacterEventLogger {
	return &CharacterEventLogger{logger: logger}
}

func (l *CharacterEventLogger) Handle(ctx context.Context, event events.Event) {
	switch e := event.(type) {
	case characterEvents.CharactersListSucceeded:
		l.logCharactersListSucceeded(ctx, e)
	case characterEvents.CharactersListFailed:
		l.logCharactersListFailed(ctx, e)

	case characterEvents.CharacterGetSucceeded:
		l.logCharacterGetSucceeded(ctx, e)
	case characterEvents.CharacterGetFailed:
		l.logCharacterGetFailed(ctx, e)

	case characterEvents.CharacterCreateSucceeded:
		l.logCharacterCreateSucceeded(ctx, e)
	case characterEvents.CharacterCreateFailed:
		l.logCharacterCreateFailed(ctx, e)

	case characterEvents.CharacterUpdateSucceeded:
		l.logCharacterUpdateSucceeded(ctx, e)
	case characterEvents.CharacterUpdateFailed:
		l.logCharacterUpdateFailed(ctx, e)

	case characterEvents.CharacterDeleteSucceeded:
		l.logCharacterDeleteSucceeded(ctx, e)
	case characterEvents.CharacterDeleteFailed:
		l.logCharacterDeleteFailed(ctx, e)

	case characterEvents.CharacterHealthGetSucceeded:
		l.logCharacterHealthGetSucceeded(ctx, e)
	case characterEvents.CharacterHealthGetFailed:
		l.logCharacterHealthGetFailed(ctx, e)

	case characterEvents.CharacterHealthUpsertSucceeded:
		l.logCharacterHealthUpsertSucceeded(ctx, e)
	case characterEvents.CharacterHealthUpsertFailed:
		l.logCharacterHealthUpsertFailed(ctx, e)

	case characterEvents.CharacterHealthDeleteSucceeded:
		l.logCharacterHealthDeleteSucceeded(ctx, e)
	case characterEvents.CharacterHealthDeleteFailed:
		l.logCharacterHealthDeleteFailed(ctx, e)

	case characterEvents.CharacterCharacteristicsGetSucceeded:
		l.logCharacterCharacteristicsGetSucceeded(ctx, e)
	case characterEvents.CharacterCharacteristicsGetFailed:
		l.logCharacterCharacteristicsGetFailed(ctx, e)

	case characterEvents.CharacterCharacteristicsUpsertSucceeded:
		l.logCharacterCharacteristicsUpsertSucceeded(ctx, e)
	case characterEvents.CharacterCharacteristicsUpsertFailed:
		l.logCharacterCharacteristicsUpsertFailed(ctx, e)

	case characterEvents.CharacterCharacteristicsDeleteSucceeded:
		l.logCharacterCharacteristicsDeleteSucceeded(ctx, e)
	case characterEvents.CharacterCharacteristicsDeleteFailed:
		l.logCharacterCharacteristicsDeleteFailed(ctx, e)

	case characterEvents.CharacterNotesListSucceeded:
		l.logCharacterNotesListSucceeded(ctx, e)
	case characterEvents.CharacterNotesListFailed:
		l.logCharacterNotesListFailed(ctx, e)

	case characterEvents.CharacterNoteGetSucceeded:
		l.logCharacterNoteGetSucceeded(ctx, e)
	case characterEvents.CharacterNoteGetFailed:
		l.logCharacterNoteGetFailed(ctx, e)

	case characterEvents.CharacterNoteCreateSucceeded:
		l.logCharacterNoteCreateSucceeded(ctx, e)
	case characterEvents.CharacterNoteCreateFailed:
		l.logCharacterNoteCreateFailed(ctx, e)

	case characterEvents.CharacterNoteUpdateSucceeded:
		l.logCharacterNoteUpdateSucceeded(ctx, e)
	case characterEvents.CharacterNoteUpdateFailed:
		l.logCharacterNoteUpdateFailed(ctx, e)

	case characterEvents.CharacterNoteDeleteSucceeded:
		l.logCharacterNoteDeleteSucceeded(ctx, e)
	case characterEvents.CharacterNoteDeleteFailed:
		l.logCharacterNoteDeleteFailed(ctx, e)
	}
}

// Character
func (l *CharacterEventLogger) logCharactersListSucceeded(ctx context.Context, event characterEvents.CharactersListSucceeded) {
	l.logger.InfoContext(ctx, "characters listed",
		"event", event.EventName(),
		"user_id", event.UserID,
		"count", event.Count,
	)
}
func (l *CharacterEventLogger) logCharactersListFailed(ctx context.Context, event characterEvents.CharactersListFailed) {
	l.logger.ErrorContext(ctx, "failed to list characters",
		"event", event.EventName(),
		"user_id", event.UserID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterGetSucceeded(ctx context.Context, event characterEvents.CharacterGetSucceeded) {
	l.logger.InfoContext(ctx, "character fetched",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"name", event.Name,
	)
}
func (l *CharacterEventLogger) logCharacterGetFailed(ctx context.Context, event characterEvents.CharacterGetFailed) {
	l.logger.ErrorContext(ctx, "failed to fetch character",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterCreateSucceeded(ctx context.Context, event characterEvents.CharacterCreateSucceeded) {
	l.logger.InfoContext(ctx, "character created",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"name", event.Name,
	)
}
func (l *CharacterEventLogger) logCharacterCreateFailed(ctx context.Context, event characterEvents.CharacterCreateFailed) {
	l.logger.ErrorContext(ctx, "failed to create character",
		"event", event.EventName(),
		"user_id", event.UserID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterUpdateSucceeded(ctx context.Context, event characterEvents.CharacterUpdateSucceeded) {
	l.logger.InfoContext(ctx, "character updated",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"name", event.Name,
	)
}
func (l *CharacterEventLogger) logCharacterUpdateFailed(ctx context.Context, event characterEvents.CharacterUpdateFailed) {
	l.logger.ErrorContext(ctx, "failed to update character",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterDeleteSucceeded(ctx context.Context, event characterEvents.CharacterDeleteSucceeded) {
	l.logger.InfoContext(ctx, "character deleted",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
	)
}
func (l *CharacterEventLogger) logCharacterDeleteFailed(ctx context.Context, event characterEvents.CharacterDeleteFailed) {
	l.logger.ErrorContext(ctx, "failed to delete character",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

// HealthState
func (l *CharacterEventLogger) logCharacterHealthGetSucceeded(ctx context.Context, event characterEvents.CharacterHealthGetSucceeded) {
	l.logger.InfoContext(ctx, "character health fetched",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
	)
}
func (l *CharacterEventLogger) logCharacterHealthGetFailed(ctx context.Context, event characterEvents.CharacterHealthGetFailed) {
	l.logger.ErrorContext(ctx, "failed to fetch character health",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterHealthUpsertSucceeded(ctx context.Context, event characterEvents.CharacterHealthUpsertSucceeded) {
	l.logger.InfoContext(ctx, "character health upserted",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
	)
}
func (l *CharacterEventLogger) logCharacterHealthUpsertFailed(ctx context.Context, event characterEvents.CharacterHealthUpsertFailed) {
	l.logger.ErrorContext(ctx, "failed to upsert character health",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterHealthDeleteSucceeded(ctx context.Context, event characterEvents.CharacterHealthDeleteSucceeded) {
	l.logger.InfoContext(ctx, "character health deleted",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
	)
}
func (l *CharacterEventLogger) logCharacterHealthDeleteFailed(ctx context.Context, event characterEvents.CharacterHealthDeleteFailed) {
	l.logger.ErrorContext(ctx, "failed to delete character health",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

// Characteristics
func (l *CharacterEventLogger) logCharacterCharacteristicsGetSucceeded(ctx context.Context, event characterEvents.CharacterCharacteristicsGetSucceeded) {
	l.logger.InfoContext(ctx, "character characteristics fetched",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
	)
}
func (l *CharacterEventLogger) logCharacterCharacteristicsGetFailed(ctx context.Context, event characterEvents.CharacterCharacteristicsGetFailed) {
	l.logger.ErrorContext(ctx, "failed to fetch character characteristics",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterCharacteristicsUpsertSucceeded(ctx context.Context, event characterEvents.CharacterCharacteristicsUpsertSucceeded) {
	l.logger.InfoContext(ctx, "character characteristics upserted",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
	)
}
func (l *CharacterEventLogger) logCharacterCharacteristicsUpsertFailed(ctx context.Context, event characterEvents.CharacterCharacteristicsUpsertFailed) {
	l.logger.ErrorContext(ctx, "failed to upsert character characteristics",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterCharacteristicsDeleteSucceeded(ctx context.Context, event characterEvents.CharacterCharacteristicsDeleteSucceeded) {
	l.logger.InfoContext(ctx, "character characteristics deleted",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
	)
}
func (l *CharacterEventLogger) logCharacterCharacteristicsDeleteFailed(ctx context.Context, event characterEvents.CharacterCharacteristicsDeleteFailed) {
	l.logger.ErrorContext(ctx, "failed to delete character characteristics",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

// Note
func (l *CharacterEventLogger) logCharacterNotesListSucceeded(ctx context.Context, event characterEvents.CharacterNotesListSucceeded) {
	l.logger.InfoContext(ctx, "character notes listed",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"count", event.Count,
	)
}
func (l *CharacterEventLogger) logCharacterNotesListFailed(ctx context.Context, event characterEvents.CharacterNotesListFailed) {
	l.logger.ErrorContext(ctx, "failed to list character notes",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterNoteGetSucceeded(ctx context.Context, event characterEvents.CharacterNoteGetSucceeded) {
	l.logger.InfoContext(ctx, "character note fetched",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"note_id", event.NoteID,
		"title", event.Title,
	)
}
func (l *CharacterEventLogger) logCharacterNoteGetFailed(ctx context.Context, event characterEvents.CharacterNoteGetFailed) {
	l.logger.ErrorContext(ctx, "failed to fetch character note",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"note_id", event.NoteID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterNoteCreateSucceeded(ctx context.Context, event characterEvents.CharacterNoteCreateSucceeded) {
	l.logger.InfoContext(ctx, "character note created",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"note_id", event.NoteID,
		"title", event.Title,
	)
}
func (l *CharacterEventLogger) logCharacterNoteCreateFailed(ctx context.Context, event characterEvents.CharacterNoteCreateFailed) {
	l.logger.ErrorContext(ctx, "failed to create character note",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterNoteUpdateSucceeded(ctx context.Context, event characterEvents.CharacterNoteUpdateSucceeded) {
	l.logger.InfoContext(ctx, "character note updated",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"note_id", event.NoteID,
		"title", event.Title,
	)
}
func (l *CharacterEventLogger) logCharacterNoteUpdateFailed(ctx context.Context, event characterEvents.CharacterNoteUpdateFailed) {
	l.logger.ErrorContext(ctx, "failed to update character note",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"note_id", event.NoteID,
		"error", event.Err,
	)
}

func (l *CharacterEventLogger) logCharacterNoteDeleteSucceeded(ctx context.Context, event characterEvents.CharacterNoteDeleteSucceeded) {
	l.logger.InfoContext(ctx, "character note deleted",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"note_id", event.NoteID,
	)
}
func (l *CharacterEventLogger) logCharacterNoteDeleteFailed(ctx context.Context, event characterEvents.CharacterNoteDeleteFailed) {
	l.logger.ErrorContext(ctx, "failed to delete character note",
		"event", event.EventName(),
		"user_id", event.UserID,
		"character_id", event.CharacterID,
		"note_id", event.NoteID,
		"error", event.Err,
	)
}
