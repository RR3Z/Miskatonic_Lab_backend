package tests

import (
	"context"
	"log/slog"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	eventLogging "github.com/RR3Z/Miskatonic_Lab_backend/pkg/observability/logging"
	"github.com/stretchr/testify/require"
)

func TestEventLoggerLogsCharacterSuccessWithDescriptorFields(t *testing.T) {
	handler := &recordingSlogHandler{}
	logger := eventLogging.NewDefaultEventLogger(slog.New(handler))

	logger.Handle(context.Background(), characterEvents.CharacterHealthUpsertSucceeded{
		UserID:      "user_1",
		CharacterID: "character_1",
	})

	record := requireSingleRecord(t, handler)
	require.Equal(t, slog.LevelInfo, record.level)
	require.Equal(t, "domain event succeeded", record.message)
	require.Equal(t, "character.health.upsert_succeeded", record.attrs["event"])
	require.Equal(t, "character", record.attrs["domain"])
	require.Equal(t, "health", record.attrs["resource"])
	require.Equal(t, "upsert", record.attrs["action"])
	require.Equal(t, "succeeded", record.attrs["outcome"])
	require.Equal(t, "user_1", record.attrs["user_id"])
	require.Equal(t, "character_1", record.attrs["character_id"])
}

func TestEventLoggerLogsFailedEventAtErrorLevel(t *testing.T) {
	handler := &recordingSlogHandler{}
	logger := eventLogging.NewDefaultEventLogger(slog.New(handler))
	err := errEventLoggerTest

	logger.Handle(context.Background(), characterEvents.CharacterUpdateFailed{
		UserID:      "user_1",
		CharacterID: "character_1",
		Err:         err,
	})

	record := requireSingleRecord(t, handler)
	require.Equal(t, slog.LevelError, record.level)
	require.Equal(t, "domain event failed", record.message)
	require.Equal(t, "failed", record.attrs["outcome"])
	require.Equal(t, err, record.attrs["error"])
}

func TestEventLoggerLogsSkippedEventAtWarnLevel(t *testing.T) {
	handler := &recordingSlogHandler{}
	logger := eventLogging.NewDefaultEventLogger(slog.New(handler))

	logger.Handle(context.Background(), characterEvents.CharacterDerivedStatsAutoRecalculateSkipped{
		UserID:      "user_1",
		CharacterID: "character_1",
		Source:      "characteristics",
		Reason:      "missing_values",
	})

	record := requireSingleRecord(t, handler)
	require.Equal(t, slog.LevelWarn, record.level)
	require.Equal(t, "domain event skipped", record.message)
	require.Equal(t, "skipped", record.attrs["outcome"])
	require.Equal(t, "characteristics", record.attrs["source"])
	require.Equal(t, "missing_values", record.attrs["reason"])
}

func TestEventLoggerLogsDiceRoomIDAndOmitsDetails(t *testing.T) {
	handler := &recordingSlogHandler{}
	logger := eventLogging.NewDefaultEventLogger(slog.New(handler))
	roomID := "room_1"

	logger.Handle(context.Background(), diceEvents.DiceRollMakeSucceeded{
		UserID:      "user_1",
		CharacterID: "character_1",
		RollID:      "roll_1",
		Expression:  "1d20",
		Result:      13,
		Details:     []byte(`[{"type":"dice"}]`),
		RoomID:      &roomID,
	})

	record := requireSingleRecord(t, handler)
	require.Equal(t, slog.LevelInfo, record.level)
	require.Equal(t, "dice", record.attrs["domain"])
	require.Equal(t, "dice_roll", record.attrs["resource"])
	require.Equal(t, "room_1", record.attrs["room_id"])
	require.Equal(t, "roll_1", record.attrs["roll_id"])
	require.Equal(t, int64(13), record.attrs["result"])
	require.NotContains(t, record.attrs, "details")
}

func TestEventLoggerOmitsNilPointerFields(t *testing.T) {
	handler := &recordingSlogHandler{}
	logger := eventLogging.NewDefaultEventLogger(slog.New(handler))

	logger.Handle(context.Background(), diceEvents.DiceRollMakeSucceeded{
		UserID:      "user_1",
		CharacterID: "character_1",
		RollID:      "roll_1",
		Expression:  "1d20",
		Result:      13,
	})

	record := requireSingleRecord(t, handler)
	require.NotContains(t, record.attrs, "room_id")
}

func TestEventLoggerLogsRoomEvents(t *testing.T) {
	handler := &recordingSlogHandler{}
	logger := eventLogging.NewDefaultEventLogger(slog.New(handler))

	logger.Handle(context.Background(), roomEvents.RoomCleanupSucceeded{
		InactiveDeleted: 1,
		InvalidDeleted:  2,
		DeletedCount:    3,
	})

	record := requireSingleRecord(t, handler)
	require.Equal(t, slog.LevelInfo, record.level)
	require.Equal(t, "room", record.attrs["domain"])
	require.Equal(t, "room_cleanup", record.attrs["resource"])
	require.Equal(t, int64(1), record.attrs["inactive_deleted"])
	require.Equal(t, int64(2), record.attrs["invalid_deleted"])
	require.Equal(t, int64(3), record.attrs["deleted_count"])
}

func TestEventLoggerIgnoresUnknownEvents(t *testing.T) {
	handler := &recordingSlogHandler{}
	logger := eventLogging.NewEventLogger(slog.New(handler), events.NewDescriptorRegistry(characterEvents.Descriptors()))

	logger.Handle(context.Background(), diceEvents.DiceRollMakeSucceeded{})

	require.Empty(t, handler.records)
}
