package main

import (
	"context"
	"fmt"
	"log"

	"github.com/RR3Z/Miskatonic_Lab_backend/internal/testdb"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
)

const seedSQL = `
DO $$
DECLARE
    tables_to_truncate text;
BEGIN
    SELECT string_agg(format('%I.%I', schemaname, tablename), ', ')
      INTO tables_to_truncate
      FROM pg_tables
     WHERE schemaname = 'public'
       AND tablename <> 'schema_migrations';

    IF tables_to_truncate IS NOT NULL THEN
        EXECUTE 'TRUNCATE TABLE ' || tables_to_truncate || ' RESTART IDENTITY CASCADE';
    END IF;
END $$;

INSERT INTO users (id, username, email)
VALUES ('test_seed_user', 'test_seed_user', 'test.seed@example.com');

INSERT INTO characters (id, user_id, name, occupation, age, residence, birthplace)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'test_seed_user',
    'Dr. Henry Armitage',
    'Librarian',
    63,
    'Arkham',
    'New Hampshire'
);

INSERT INTO skills_categories (id, name)
VALUES ('00000000-0000-0000-0000-000000000101', 'Investigation');

INSERT INTO skills_specialties (id, name, description, base_value)
VALUES (
    '00000000-0000-0000-0000-000000000201',
    'History',
    'Seed specialty for database tests',
    5
);
`

func main() {
	if err := testdb.LoadEnv(); err != nil {
		log.Fatal(err)
	}
	databaseURL := testdb.DatabaseURL()
	if err := testdb.ValidateLocal(databaseURL); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := repository.NewPostgresDBFromURL(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, seedSQL); err != nil {
		log.Fatal(err)
	}
	fmt.Println("local test database seeded")
}
