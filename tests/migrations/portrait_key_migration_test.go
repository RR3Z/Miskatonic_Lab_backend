package tests

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestPortraitKeyMigrationRenamesColumnClearsLegacyURLsAndRollsBack(t *testing.T) {
	loadMigrationTestEnv(t)

	if !migrationSmokeEnabled() {
		t.Skip("set MIGRATION_SMOKE_TESTS=1 to run migration rollback smoke")
	}

	databaseURL := strings.TrimSpace(os.Getenv("MIGRATION_SMOKE_DATABASE_URL"))
	require.NotEmpty(t, databaseURL, "MIGRATION_SMOKE_DATABASE_URL must point to a dedicated disposable database")
	migratePath, err := exec.LookPath("migrate")
	require.NoError(t, err, "migrate CLI must be available in PATH")
	root := migrationRepoRoot(t)
	ensureLatestMigrationOnCleanup(t, root, migratePath, databaseURL)

	runMigrate(t, root, migratePath, databaseURL, "up")
	runMigrate(t, root, migratePath, databaseURL, "down", "1")

	connection, err := pgx.Connect(context.Background(), databaseURL)
	require.NoError(t, err)
	t.Cleanup(func() { _ = connection.Close(context.Background()) })

	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)
	userID := "portrait_migration_user_" + suffix
	_, err = connection.Exec(context.Background(), `
		INSERT INTO users (id, username, email)
		VALUES ($1, $2, $3)
	`, userID, "portrait_migration_"+suffix, "portrait.migration+"+suffix+"@example.com")
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = connection.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
	})

	characterID := "11111111-1111-1111-1111-" + suffix[len(suffix)-12:]
	legacyURL := "https://legacy.example.test/portrait.png"
	_, err = connection.Exec(context.Background(), `
		INSERT INTO characters (id, user_id, name, portrait_url)
		VALUES ($1, $2, 'Legacy Portrait Investigator', $3)
	`, characterID, userID, legacyURL)
	require.NoError(t, err)
	requireColumnExists(t, connection, "portrait_url", true)
	requireColumnExists(t, connection, "portrait_key", false)

	runMigrate(t, root, migratePath, databaseURL, "up", "1")
	requireColumnExists(t, connection, "portrait_url", false)
	requireColumnExists(t, connection, "portrait_key", true)
	var portraitKey *string
	require.NoError(t, connection.QueryRow(context.Background(), "SELECT portrait_key FROM characters WHERE id = $1", characterID).Scan(&portraitKey))
	require.Nil(t, portraitKey)

	runMigrate(t, root, migratePath, databaseURL, "down", "1")
	requireColumnExists(t, connection, "portrait_key", false)
	requireColumnExists(t, connection, "portrait_url", true)
	var rolledBackURL *string
	require.NoError(t, connection.QueryRow(context.Background(), "SELECT portrait_url FROM characters WHERE id = $1", characterID).Scan(&rolledBackURL))
	require.Nil(t, rolledBackURL)

	runMigrate(t, root, migratePath, databaseURL, "up", "1")
}

func requireColumnExists(t *testing.T, connection *pgx.Conn, column string, expected bool) {
	t.Helper()
	var exists bool
	err := connection.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_schema = 'public'
			  AND table_name = 'characters'
			  AND column_name = $1
		)
	`, column).Scan(&exists)
	require.NoError(t, err)
	require.Equal(t, expected, exists)
}

func ensureLatestMigrationOnCleanup(t *testing.T, root string, migratePath string, databaseURL string) {
	t.Helper()
	t.Cleanup(func() {
		cmd := exec.Command(migratePath, "-path", "migrations", "-database", databaseURL, "up")
		cmd.Dir = root
		output, err := cmd.CombinedOutput()
		if err != nil && !strings.Contains(strings.ToLower(string(output)), "no change") {
			t.Errorf("restore latest migration failed: %v: %s", err, output)
		}
	})
}
