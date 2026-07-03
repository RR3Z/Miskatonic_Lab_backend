package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func requireSingleRecord(t *testing.T, handler *recordingSlogHandler) recordedLog {
	t.Helper()

	require.Len(t, handler.records, 1)
	return handler.records[0]
}
