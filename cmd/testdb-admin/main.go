package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/internal/testdb"
)

func main() {
	validateOnly := flag.Bool("validate", false, "validate migration smoke database configuration")
	reset := flag.Bool("reset", false, "drop and recreate the migration smoke database")
	flag.Parse()

	if *validateOnly == *reset {
		fmt.Fprintln(os.Stderr, "set exactly one of -validate or -reset")
		os.Exit(2)
	}
	if err := testdb.LoadEnv(); err != nil {
		fmt.Fprintln(os.Stderr, "load environment:", err)
		os.Exit(1)
	}

	smokeURL := strings.TrimSpace(os.Getenv("MIGRATION_SMOKE_DATABASE_URL"))
	testURL := testdb.DatabaseURL()
	if *validateOnly {
		if _, err := testdb.ValidateMigrationSmokeURL(smokeURL, testURL); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("migration smoke database configuration is valid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := testdb.ResetMigrationSmokeDatabase(ctx, smokeURL, testURL); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("migration smoke database reset")
}
