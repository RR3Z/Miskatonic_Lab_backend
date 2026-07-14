package testdb

import "testing"

func TestValidateLocal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		databaseURL string
		wantError   bool
	}{
		{name: "localhost test database", databaseURL: defaultDatabaseURL},
		{name: "loopback test database", databaseURL: "postgres://user:pass@127.0.0.1:5432/app_test?sslmode=disable"},
		{name: "reject production name", databaseURL: "postgres://user:pass@localhost:5432/app", wantError: true},
		{name: "reject Supabase", databaseURL: "postgres://user:pass@aws-0-eu.pooler.supabase.com:5432/app_test", wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateLocal(test.databaseURL)
			if test.wantError && err == nil {
				t.Fatal("expected validation error")
			}
			if !test.wantError && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestValidateMigrationSmokeURL(t *testing.T) {
	t.Parallel()

	testURL := "postgres://user:pass@localhost:5433/miskatonic_lab_test?sslmode=disable"
	tests := []struct {
		name      string
		testURL   string
		smokeURL  string
		wantError bool
	}{
		{name: "dedicated local smoke database", testURL: testURL, smokeURL: "postgres://user:pass@localhost:5433/miskatonic_lab_migration_smoke_test?sslmode=disable"},
		{name: "reject missing smoke URL", testURL: testURL, smokeURL: "", wantError: true},
		{name: "reject Supabase smoke URL", testURL: testURL, smokeURL: "postgres://user:pass@aws-0-eu.pooler.supabase.com:5432/miskatonic_lab_migration_smoke_test?sslmode=require", wantError: true},
		{name: "reject production test URL", testURL: "postgres://user:pass@aws-0-eu.pooler.supabase.com:5432/miskatonic_lab_test?sslmode=require", smokeURL: "postgres://user:pass@localhost:5433/miskatonic_lab_migration_smoke_test?sslmode=disable", wantError: true},
		{name: "reject normal test database", testURL: testURL, smokeURL: testURL, wantError: true},
		{name: "reject generic test database", testURL: testURL, smokeURL: "postgres://user:pass@localhost:5433/miskatonic_lab_other_test?sslmode=disable", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateMigrationSmokeURL(tt.smokeURL, tt.testURL)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected validation error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}
