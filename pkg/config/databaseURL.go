package config

import (
	"fmt"
	"strings"
)

func DatabaseURL(cfg PostgresDBConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
}

func DatabaseURLWithFallback(databaseURL string, cfg PostgresDBConfig) string {
	trimmedDatabaseURL := strings.TrimSpace(databaseURL)
	if trimmedDatabaseURL != "" {
		return trimmedDatabaseURL
	}

	return DatabaseURL(cfg)
}
