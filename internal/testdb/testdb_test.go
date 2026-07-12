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
