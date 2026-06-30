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

## Repository Tests

- Use a real PostgreSQL test database.
- Verify generated `sqlc` queries, migrations, constraints, foreign keys, upserts, and delete behavior.
- Do not mock `pgx` or hand-edit generated files.

Character table integration tests live in `tests/character/integration`. They use the real PostgreSQL database and generated `sqlc` queries to verify character create/get/list/update/delete behavior, user ownership scoping, foreign-key constraints, check constraints, nullable fields, and cascade deletion from `users` to `characters`.

Health, sanity, magic, and luck table integration tests also live in `tests/character/integration`. They verify state upsert/get/delete behavior, database defaults, partial updates, owner scoping, negative-value CHECK constraints, and cascade deletion from `characters` to the related state row.

## End-To-End Tests

- Use real HTTP calls against a configured test server and real PostgreSQL test database.
- Keep this suite small and scenario-based: user provisioning, `/api/me`, character creation, subresource updates, full character read, deletion, and access denial for another user's character.

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
