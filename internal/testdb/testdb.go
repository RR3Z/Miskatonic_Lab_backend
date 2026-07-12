package testdb

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const defaultDatabaseURL = "postgres://miskatonic_user:miskatonic_password@localhost:5433/miskatonic_lab_test?sslmode=disable"

func LoadEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return nil
		}
		dir = parent
	}
}

func DatabaseURL() string {
	value := strings.TrimSpace(os.Getenv("TEST_DATABASE_URL"))
	if value == "" {
		return defaultDatabaseURL
	}
	return value
}

func ValidateLocal(databaseURL string) error {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return fmt.Errorf("parse TEST_DATABASE_URL: %w", err)
	}

	host := strings.ToLower(parsed.Hostname())
	if host != "localhost" && !net.ParseIP(host).IsLoopback() {
		return fmt.Errorf("TEST_DATABASE_URL must use localhost or loopback IP, got %q", host)
	}
	if !strings.HasSuffix(strings.TrimPrefix(parsed.Path, "/"), "_test") {
		return fmt.Errorf("TEST_DATABASE_URL database name must end with _test")
	}
	return nil
}

func Open(t testing.TB) *pgxpool.Pool {
	t.Helper()
	if err := LoadEnv(); err != nil {
		t.Fatalf("load test env: %v", err)
	}

	databaseURL := DatabaseURL()
	if err := ValidateLocal(databaseURL); err != nil {
		t.Fatal(err)
	}

	pool, err := repository.NewPostgresDBFromURL(context.Background(), databaseURL)
	if err != nil {
		t.Fatalf("connect to local test database: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}
