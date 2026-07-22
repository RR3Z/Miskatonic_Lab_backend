package tests

import (
	"testing"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func requireCleanupDeletedRoomIDs(
	t *testing.T,
	result model.CleanupRoomsResult,
	invalidRoomIDs ...pgtype.UUID,
) {
	t.Helper()

	for _, roomID := range invalidRoomIDs {
		require.Contains(t, result.InvalidDeletedRoomIDs, roomID)
		require.Contains(t, result.DeletedRoomIDs, roomID)
	}
}
