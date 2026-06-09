package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func testCharacteristics(strength int16, size int16, dexterity int16) db.Characteristic {
	return db.Characteristic{
		Strength:  int16Pointer(strength),
		Size:      int16Pointer(size),
		Dexterity: int16Pointer(dexterity),
	}
}

func splitTotal(total int16) (int16, int16) {
	strength := total / 2
	size := total - strength
	return strength, size
}

func int16Pointer(value int16) *int16 {
	return &value
}

func requireInt16PointerValue(t *testing.T, actual *int16, expected int16) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected, *actual)
}

func requireStringPointerValue(t *testing.T, actual *string, expected string) {
	t.Helper()

	require.NotNil(t, actual)
	require.Equal(t, expected, *actual)
}

func testUUID(value string) pgtype.UUID {
	var uuid pgtype.UUID
	if err := uuid.Scan(value); err != nil {
		panic(err)
	}

	return uuid
}
