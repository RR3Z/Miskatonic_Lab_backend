package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestParseAllowedOriginsTrimsAndDropsBlanks(t *testing.T) {
	origins := config.ParseAllowedOrigins(" http://localhost:3000, ,https://app.example.com,, * ")

	require.Equal(t, []string{"http://localhost:3000", "https://app.example.com", "*"}, origins)
}

func TestParseAllowedOriginsReturnsEmptySliceForEmptyInput(t *testing.T) {
	origins := config.ParseAllowedOrigins("  , ")

	require.Empty(t, origins)
	require.NotNil(t, origins)
}

func TestDatabaseURLFormatsPostgresConnectionString(t *testing.T) {
	url := config.DatabaseURL(config.PostgresDBConfig{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "secret",
		DBName:   "miskatonic",
		SSLMode:  "disable",
	})

	require.Equal(t, "postgres://postgres:secret@localhost:5432/miskatonic?sslmode=disable", url)
}
