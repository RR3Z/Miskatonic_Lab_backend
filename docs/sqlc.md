# sqlc

`sqlc` generates type-safe Go code from plain SQL. In this project the schema comes from `migrations/`, queries live in `pkg/repository/queries/`, and generated Go files go to `pkg/repository/db/`.

## Generate

```powershell
npm run sqlc:generate
```

The script uses a globally installed `sqlc` when it exists. Otherwise it runs the pinned CLI version through Go:

```powershell
go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1 generate
```

## Workflow

1. Add or change a migration in `migrations/`.
2. Add named SQL queries in `pkg/repository/queries/*.sql`.
3. Run `npm run sqlc:generate`.
4. Use generated methods through `repository.Repository.Queries`.

Example query:

```sql
-- name: GetUserByClerkID :one
SELECT *
FROM users
WHERE clerk_user_id = $1;
```

This becomes a Go method:

```go
user, err := repos.Queries.GetUserByClerkID(ctx, clerkUserID)
```

## Query Modes

`-- name: Something :one` returns one row and an error.

`-- name: Something :many` returns a slice.

`-- name: Something :exec` runs a statement without returning rows.

`INSERT/UPDATE/DELETE ... RETURNING *` can use `:one` when you want the changed row back.

## Types

The project uses `pgx/v5`, so generated UUID and timestamp fields use `pgtype.UUID` and `pgtype.Timestamptz`. Nullable text and smallint columns are generated as pointers because `emit_pointers_for_null_types` is enabled in `sqlc.yaml`.
