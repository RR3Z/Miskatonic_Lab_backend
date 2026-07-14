package testdb

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
)

const migrationSmokeDatabaseSuffix = "_migration_smoke_test"

// ValidateMigrationSmokeURL rejects any target except a dedicated local smoke database.
func ValidateMigrationSmokeURL(smokeURL string, testURL string) (*url.URL, error) {
	if err := ValidateLocal(testURL); err != nil {
		return nil, fmt.Errorf("TEST_DATABASE_URL: %w", err)
	}
	if err := ValidateLocal(smokeURL); err != nil {
		return nil, fmt.Errorf("MIGRATION_SMOKE_DATABASE_URL: %w", err)
	}

	parsedSmokeURL, err := url.Parse(smokeURL)
	if err != nil {
		return nil, fmt.Errorf("parse MIGRATION_SMOKE_DATABASE_URL: %w", err)
	}
	smokeDatabaseName := strings.TrimPrefix(parsedSmokeURL.Path, "/")
	if !strings.HasSuffix(smokeDatabaseName, migrationSmokeDatabaseSuffix) {
		return nil, fmt.Errorf("MIGRATION_SMOKE_DATABASE_URL database name must end with %q", migrationSmokeDatabaseSuffix)
	}
	if strings.ContainsAny(smokeDatabaseName, `"'\\`) {
		return nil, fmt.Errorf("MIGRATION_SMOKE_DATABASE_URL database name contains unsupported characters")
	}

	parsedTestURL, err := url.Parse(testURL)
	if err != nil {
		return nil, fmt.Errorf("parse TEST_DATABASE_URL: %w", err)
	}
	if sameDatabaseTarget(parsedSmokeURL, parsedTestURL) {
		return nil, fmt.Errorf("MIGRATION_SMOKE_DATABASE_URL must differ from TEST_DATABASE_URL")
	}

	return parsedSmokeURL, nil
}

// ResetMigrationSmokeDatabase drops and recreates the validated disposable smoke database.
func ResetMigrationSmokeDatabase(ctx context.Context, smokeURL string, testURL string) error {
	parsedSmokeURL, err := ValidateMigrationSmokeURL(smokeURL, testURL)
	if err != nil {
		return err
	}

	databaseName := strings.TrimPrefix(parsedSmokeURL.Path, "/")
	adminURL := *parsedSmokeURL
	adminURL.Path = "/postgres"

	connection, err := pgx.Connect(ctx, adminURL.String())
	if err != nil {
		return fmt.Errorf("connect to local postgres admin database: %w", err)
	}
	defer connection.Close(ctx)

	if _, err := connection.Exec(ctx, `
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = $1 AND pid <> pg_backend_pid()
	`, databaseName); err != nil {
		return fmt.Errorf("terminate migration smoke database connections: %w", err)
	}

	quotedDatabaseName := `"` + strings.ReplaceAll(databaseName, `"`, `""`) + `"`
	if _, err := connection.Exec(ctx, "DROP DATABASE IF EXISTS "+quotedDatabaseName); err != nil {
		return fmt.Errorf("drop migration smoke database: %w", err)
	}
	if _, err := connection.Exec(ctx, "CREATE DATABASE "+quotedDatabaseName); err != nil {
		return fmt.Errorf("create migration smoke database: %w", err)
	}

	return nil
}

func sameDatabaseTarget(left *url.URL, right *url.URL) bool {
	return strings.EqualFold(left.Hostname(), right.Hostname()) &&
		left.Port() == right.Port() &&
		strings.TrimPrefix(left.Path, "/") == strings.TrimPrefix(right.Path, "/")
}
