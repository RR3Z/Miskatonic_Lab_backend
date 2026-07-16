package tests

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/internal/testdb"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestMigrationRollbackSmoke(t *testing.T) {
	loadMigrationTestEnv(t)

	if !migrationSmokeEnabled() {
		t.Skip("set MIGRATION_SMOKE_TESTS=1 to run migration rollback smoke")
	}

	databaseURL := strings.TrimSpace(os.Getenv("MIGRATION_SMOKE_DATABASE_URL"))
	require.NotEmpty(t, databaseURL, "MIGRATION_SMOKE_DATABASE_URL must point to a dedicated disposable database")

	migratePath, err := exec.LookPath("migrate")
	require.NoError(t, err, "migrate CLI must be available in PATH")

	root := migrationRepoRoot(t)
	runMigrate(t, root, migratePath, databaseURL, "up")
	runMigrate(t, root, migratePath, databaseURL, "down", "1")
	runMigrate(t, root, migratePath, databaseURL, "up", "1")
	versionOutput := runMigrate(t, root, migratePath, databaseURL, "version")
	require.NotEmpty(t, strings.TrimSpace(versionOutput))
}

func TestRemoveCharacterPlayerNameMigrationDropsDataAndRollsBackSchema(t *testing.T) {
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

	requireColumnExists(t, connection, "player_name", true)
	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)
	userID := "player_name_migration_user_" + suffix
	_, err = connection.Exec(context.Background(), `
		INSERT INTO users (id, username, email)
		VALUES ($1, $2, $3)
	`, userID, "player_name_migration_"+suffix, "player.name.migration+"+suffix+"@example.com")
	require.NoError(t, err)

	characterID := "22222222-2222-2222-2222-" + suffix[len(suffix)-12:]
	_, err = connection.Exec(context.Background(), `
		INSERT INTO characters (id, user_id, name, player_name)
		VALUES ($1, $2, 'Legacy Player Name', 'Legacy Player')
	`, characterID, userID)
	require.NoError(t, err)

	runMigrate(t, root, migratePath, databaseURL, "up", "1")
	requireColumnExists(t, connection, "player_name", false)

	runMigrate(t, root, migratePath, databaseURL, "down", "1")
	requireColumnExists(t, connection, "player_name", true)
	var playerName *string
	require.NoError(t, connection.QueryRow(context.Background(), "SELECT player_name FROM characters WHERE id = $1", characterID).Scan(&playerName))
	require.Nil(t, playerName)

	runMigrate(t, root, migratePath, databaseURL, "up", "1")
}

func loadMigrationTestEnv(t *testing.T) {
	t.Helper()
	require.NoError(t, testdb.LoadEnv())
}

func migrationSmokeEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("MIGRATION_SMOKE_TESTS")))
	return value == "1" || value == "true" || value == "yes"
}

func migrationRepoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	require.NoError(t, err)
	for {
		if fileExists(filepath.Join(dir, "go.mod")) && fileExists(filepath.Join(dir, "migrations")) {
			return dir
		}
		parent := filepath.Dir(dir)
		require.NotEqual(t, dir, parent, "could not find repository root")
		dir = parent
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func runMigrate(t *testing.T, root string, migratePath string, databaseURL string, args ...string) string {
	t.Helper()

	fullArgs := append([]string{"-path", "migrations", "-database", databaseURL}, args...)
	cmd := exec.Command(migratePath, fullArgs...)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "migrate %s failed:\n%s", strings.Join(args, " "), string(output))
	return string(output)
}
