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
	}
}

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
