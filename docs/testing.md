# Testing Strategy

This project uses the standard Go test runner as the base test framework.

Run the test suite with standard Go output:

```powershell
go test ./...
```

Run domain tests with per-test pretty output:

```powershell
npm run test:pretty
```

Run every configured suite, including real Clerk integration, live E2E, and migration rollback smoke:

```powershell
$env:E2E_AUTH_TOKEN="<real Clerk session token>"
$env:E2E_SECOND_AUTH_TOKEN="<second real Clerk session token>"
$env:MIGRATION_SMOKE_DATABASE_URL="<dedicated disposable database url>"
npm run test:all
```

`test:all` runs `go test ./...`, then enables real Clerk integration tests, live E2E tests, and migration rollback smoke. Configure Clerk secrets, webhook tunnel, PostgreSQL env, E2E tokens, and a dedicated migration smoke database before using it.

## Unit Tests

- Use `testing` for test structure.
- Use `github.com/stretchr/testify/require` for assertions that should stop the test immediately.
- Prefer table tests and small named subtests with `t.Run`.
- Keep unit tests close to the package they test when they need package-private constructors or helpers.
- Use narrow hand-written fakes for package boundaries instead of mocking `pgx` or generated `sqlc` internals.
- Do not change service, handler, or other production logic just to make tests easier to write; test-only helpers and fakes should live under `tests/<domain>/`.

Service unit tests should focus on business rules and service flow. Avoid changing service architecture just to make database-heavy methods easy to fake; if a method mostly delegates to `sqlc`, prefer repository or API tests.

## Handler/API Tests

- Use `testing`, `httptest`, and `require`.
- Use fake services behind handler interfaces.
- Verify status codes, JSON bodies, route params, malformed request bodies, auth context, and `AppError` mapping.

Character CRUD handler tests live in `tests/character`. They build the public handler router with `NewHandler`, inject Clerk claims into the request context, and keep fake services/helpers split from happy-path, validation, error, binding, and method tests.

Character portrait HTTP tests live under `tests/character/unit/handler`, lifecycle/concurrency/reconciliation coverage lives under `tests/character/integration`, local filesystem storage tests are split into local-store, file-server, image-validation, and reconciliation files under `tests/storage/portrait`, and DB/storage coordination plus worker scheduling tests live under `tests/maintenance/portrait`.

Focused Character HTTP contract tests live in `tests/character/unit/handler`. They cover route mounting, malformed JSON, invalid UUIDs, representative DTO handoff, delete status, and service-error JSON mapping without touching production service logic.

Config and middleware unit tests live in `tests/config` and `tests/middleware`. They cover CORS origin parsing/headers, database URL formatting, request logging levels, and request error-code logging.

## Repository Tests

- Use a real PostgreSQL test database.
- Verify generated `sqlc` queries, migrations, constraints, foreign keys, upserts, and delete behavior.
- Do not mock `pgx` or hand-edit generated files.

Character table integration tests live in `tests/character/integration`. They use the real PostgreSQL database and generated `sqlc` queries to verify character create/get/list/update/delete behavior, user ownership scoping, foreign-key constraints, check constraints, nullable fields, and cascade deletion from `users` to `characters`.

Prepare the local PostgreSQL test database and run database-backed suites:

```powershell
npm run testdb:prepare
go test -count=1 ./tests/character/integration ./tests/diceRoller/integration ./tests/room/integration
```

Character-limit integration coverage verifies the 30-character boundary, slot reuse after deletion, per-user isolation, and concurrent creation at the boundary. Portrait integration coverage verifies DB/file lifecycle, ownership, concurrent replacement, compensation cleanup after DB or context failure, and reconciliation against referenced DB keys.

Health, sanity, magic, and luck table integration tests also live in `tests/character/integration`. They verify state upsert/get/delete behavior, database defaults, partial updates, owner scoping, negative-value CHECK constraints, and cascade deletion from `characters` to the related state row.

Room realtime integration and workability tests live in `tests/room/integration`. They cover password room creation, invite/password joins, selected-character visibility by role, room event persistence and old-to-new history order, room-wide dice/chat events, `character.changed` privacy filtering, owner leave transfer, last-member deletion, cleanup result IDs, and room cleanup deletion behavior.

Room WebSocket unit tests live in `tests/ws/unit`. They cover persisted chat broadcast, command errors, room-wide delivery, targeted delivery support through listener tests, slow-client isolation, and closing active clients when rooms are deleted.

Event infrastructure tests are split by responsibility:

- `tests/events` covers publisher behavior plus event descriptor lookup and exported prototype groups.
- `tests/observability/logging` covers descriptor-driven app logging for Character, DiceRoller, and Room events.
- `tests/room/unit/events` covers `EventPublishingRoomService`.
- `tests/listeners` covers async room side-effect listeners under `pkg/listeners/room`.

## End-To-End Tests

- Use real HTTP calls against a configured test server and real PostgreSQL test database.
- Keep this suite small and scenario-based: user provisioning, `/api/me`, character creation, subresource updates, full character read, deletion, and access denial for another user's character.

Live backend E2E tests live in `tests/e2e`. They are opt-in so the normal suite stays local and deterministic:

```powershell
$env:E2E_AUTH_TOKEN="<real Clerk session token>"
npm run test:e2e
```

Optional settings:

- `E2E_TESTS=1` enables the package; `npm run test:e2e` sets it for the current PowerShell command.
- `E2E_BASE_URL` points at the running backend, defaulting to `http://localhost:${PORT}` or `http://localhost:8000`.
- `E2E_AUTH_TOKEN` may be a raw JWT or `Bearer <jwt>`. The tests decode only the JWT `sub` to prepare local database rows; the backend still verifies the token through Clerk middleware.

Current E2E scenarios cover missing-token rejection, `/api/me`, Character create/list/get/health/delete over HTTP, ownership denial for another user's character seeded in PostgreSQL, character-limit conflict, backend-managed portrait upload/replacement/deletion with public GET/HEAD delivery, portrait validation errors, and room-context dice rolling that produces a `dice.roll` room event.

Additional live E2E gap coverage:

- `E2E_SECOND_AUTH_TOKEN` is optional. When it is present and belongs to a different Clerk subject, the multi-user E2E creates a room with the primary token, joins it with the second token, selects characters for both users, verifies GM/owner visibility, and verifies player-only selected-character visibility. When it is absent, that specific test skips with an explicit reason.
- Live WebSocket E2E creates a room over HTTP, opens `/api/rooms/{roomID}/ws` with the real auth token, sends `chat.message`, verifies the received `chat.message` socket event, then polls room history and verifies the persisted event payload.

## Clerk Integration Tests

User Clerk integration tests live in `tests/user/integration`. They call the real Clerk API through `clerk-sdk-go`, then wait until Clerk sends a real webhook to the running backend and the local PostgreSQL row changes.

These tests are opt-in because they need external services and a reachable webhook URL:

```powershell
npm run test:clerk
```

Before running them:

- Start local PostgreSQL and apply migrations.
- Start the backend with `CLERK_SECRET_KEY`, `CLERK_WEBHOOK_SIGNING_SECRET`, and PostgreSQL env vars configured.
- Expose the backend webhook URL with a public tunnel, for example `<public-url>/webhooks/clerk/user`; Clerk cannot send webhooks directly to your machine's `localhost`.
- VS Code port forwarding works for this too: forward the backend port, copy the public forwarded URL, and use `<forwarded-url>/webhooks/clerk/user` as the Clerk webhook endpoint.
- Configure the Clerk test instance webhook endpoint to send `user.created`, `user.updated`, and `user.deleted` events to that URL.
- If the backend is not running on `http://localhost:8000`, set `CLERK_LOCAL_BACKEND_URL` for the test preflight check, for example `http://localhost:3000`.
- The tests wait up to 2 minutes for each webhook by default. If Clerk delivery or port forwarding is slow, set `CLERK_WEBHOOK_WAIT_TIMEOUT`, for example `180s`.

The tests create Clerk users with `+clerk_test` email subaddresses, poll the local DB for webhook results, delete the Clerk user, and clean up their local DB row. Because they depend on real Clerk webhook delivery and a public forwarding service, occasional timing instability usually points to delayed webhook delivery, a refreshed forwarded URL, or the backend/test process using different PostgreSQL settings.

`npm run test:all` runs this suite by setting `CLERK_INTEGRATION_TESTS=1` and executing `go test ./tests/user/integration -v`. The suite still skips outside that opt-in path.

## Migration Smoke Tests

Migration rollback smoke tests live in `tests/migrations`. They are opt-in because they intentionally run `migrate up`, `migrate down 1`, `migrate up 1`, and `migrate version`. The portrait migration scenario also verifies that migration `000014` renames `portrait_url` to `portrait_key`, clears legacy client-owned URLs, and restores the expected column name across down/up.

```powershell
$env:MIGRATION_SMOKE_DATABASE_URL="<dedicated disposable database url>"
npm run test:migrations
```

- `MIGRATION_SMOKE_TESTS=1` enables the package; `npm run test:migrations` sets it for the current PowerShell command.
- `MIGRATION_SMOKE_DATABASE_URL` is required and must point to a dedicated disposable database, not the normal development database.
- The `migrate` CLI must be available in `PATH`.

## Test Trace

The trace generator runs every package under `tests/` with `go test -json -count=1` and writes one row per concrete test or subtest to `docs/test-trace.md`:

```powershell
python .agent/scripts/generate_test_trace.py
```

Rows preserve the actual result. Local tests should show `pass`; opt-in Clerk, E2E, and migration scenarios show `skip` until their guards and required environment are enabled. To record a real live E2E run in the trace, start the test backend, configure `E2E_AUTH_TOKEN`, set `E2E_TESTS=1`, and rerun the generator.
