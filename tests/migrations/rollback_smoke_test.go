package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrationRollbackSmoke(t *testing.T) {
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
