# Testing

`npm run test:all` is the complete test command. After the one-time setup and a filled `.env`, it prepares local databases, runs all local tests, starts an isolated API plus Dev Tunnel, executes real Clerk and E2E tests, runs migration smoke, and cleans up every process it started.

## One-time setup

1. Install Docker Desktop, Go, Node.js/npm, and `migrate`. Then run:

   ```powershell
   npm run test:setup
   ```

   The command downloads the Dev Tunnel CLI to the current Windows user's local tools directory if needed, opens its login flow if needed, creates an anonymous persistent tunnel and HTTP port `8001`, then writes its ID and URL to `.env`. It does not need `winget`.
   Empty values and template values such as `replace_me` are both treated as not yet configured.
   If a previously configured tunnel was deleted, `test:setup` provisions a replacement and updates `.env`.

2. Copy the URL printed by `test:setup` into Clerk Dashboard for `user.created`, `user.updated`, and `user.deleted`:

   ```text
   <CLERK_WEBHOOK_PUBLIC_URL>/webhooks/clerk/user
   ```

3. Validate after saving the Clerk setting:

   ```powershell
   npm run test:setup
   ```

The tunnel must allow anonymous access because Clerk webhook delivery cannot complete an interactive Dev Tunnel login. The runner hosts the configured tunnel only during Clerk tests and never changes Clerk Dashboard settings. Dev Tunnel CLI management commands are documented in the [Microsoft reference](https://learn.microsoft.com/en-us/azure/developer/dev-tunnels/cli-commands).

## Environment contract

| Group | Variables | Used by |
| --- | --- | --- |
| Runtime | `DATABASE_URL` or `POSTGRES_*`, `CLERK_SECRET_KEY`, `CLERK_AUTHORIZED_PARTIES`, `CLERK_WEBHOOK_SIGNING_SECRET`, CORS/storage settings | normal backend and isolated test backend |
| Local DB | `TEST_DATABASE_URL`, `MIGRATION_SMOKE_DATABASE_URL` | every test command; smoke DB is reset destructively |
| Live E2E | `E2E_TEST1_MAIL`, `E2E_TEST2_MAIL` | `test:e2e`, `test:all` |
| Live webhook | `TEST_BACKEND_PORT`, `DEVTUNNEL_TUNNEL_ID`, `CLERK_WEBHOOK_PUBLIC_URL` | `test:clerk`, `test:all` |
| Optional | `CLERK_WEBHOOK_WAIT_TIMEOUT` | real webhook delivery timeout |

Do not set `E2E_AUTH_TOKEN`, `E2E_SECOND_AUTH_TOKEN`, E2E passwords, `E2E_TESTS`, `CLERK_INTEGRATION_TESTS`, `MIGRATION_SMOKE_TESTS`, `E2E_BASE_URL`, or `CLERK_LOCAL_BACKEND_URL`. They are internal runner state, not user configuration.

`TEST_DATABASE_URL` must be a loopback database ending in `_test`. `MIGRATION_SMOKE_DATABASE_URL` must be different and end in `_migration_smoke_test`; the runner drops and recreates it before smoke tests. Supabase and production URLs are rejected before destructive SQL runs.

## Commands

| Command | Runs | Required env |
| --- | --- | --- |
| `npm run test:setup` | installs/checks CLI, login, persistent tunnel, port, DB configuration | runtime, local DB, E2E emails, `TEST_BACKEND_PORT`; tunnel ID/URL are created when empty |
| `npm run test:local` | Docker test DB preparation and `go test ./...` | local DB values |
| `npm run test:e2e` | isolated backend and live HTTP/multi-user/WebSocket E2E | local DB, Clerk runtime, E2E emails, test port |
| `npm run test:clerk` | isolated backend, hosted tunnel, real Clerk webhook tests | local DB, Clerk runtime, tunnel values |
| `npm run test:migrations` | reset smoke DB and migration up/down/up tests | local DB values |
| `npm run test:all` | every suite in the order below | all values |

`npm run test:all` order:

1. Validate tools and environment.
2. Start `postgres-test`, apply migrations, seed `TEST_DATABASE_URL`, and reset the smoke DB.
3. Run `go test ./...` with all live guards disabled.
4. Start the isolated backend on `TEST_BACKEND_PORT` with `DATABASE_URL=TEST_DATABASE_URL` and temporary portrait storage.
5. Host the persistent Dev Tunnel, run real Clerk webhook tests, then run live E2E.
6. Run migration smoke and stop the tunnel/backend process trees. Docker PostgreSQL remains available for later local work.

## Fail-fast and cleanup

- A missing required `.env` value, Docker engine, failed tunnel login/provisioning, occupied test backend port, or invalid DB URL stops before the relevant suite starts. `test:setup` downloads the Dev Tunnel CLI if it is absent.
- The runner never stops an existing process on the test port; choose a dedicated unused port.
- Clerk integration removes created Clerk users and local rows. E2E revokes its temporary Clerk sessions. The runner removes its temporary portrait storage and stops only its own backend/tunnel process trees.
- `test:all` prints compact stage summaries on success and expands the captured command/test output only on failure. `npm run test:pretty` remains the optional per-test detailed local-output helper; it does not prepare external dependencies.
