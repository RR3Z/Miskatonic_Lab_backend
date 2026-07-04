<p align="center">
  <img src="assets/readme-logo-esoteric.webp" alt="Miskatonic Lab" width="720">
</p>

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.26.3-00ADD8?logo=go&logoColor=white">
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white">
  <img alt="Docker" src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white">
  <img alt="Clerk" src="https://img.shields.io/badge/Auth-Clerk-6C47FF?logo=clerk&logoColor=white">
  <img alt="chi" src="https://img.shields.io/badge/Router-chi-111111">
  <img alt="sqlc" src="https://img.shields.io/badge/SQL-sqlc-2F855A">
  <img alt="WebSocket" src="https://img.shields.io/badge/Realtime-WebSocket-0F766E">
</p>

# Miskatonic Lab Backend

Backend for a Call of Cthulhu character and room-management app.

It provides a small HTTP API for users, characters, dice rolls, rooms, room events, and room WebSocket chat. PostgreSQL stores the domain data, Clerk handles authentication, and sqlc generates the typed repository layer.

## Features

- Character sheets with characteristics, skills, states, backstory, finances, and notes.
- Dice rolls with persisted history and optional room context.
- Rooms with members, selected characters, event history, and WebSocket delivery.
- Clerk user webhooks and protected API routes.
- Focused Go tests for handlers, services, integrations, WebSocket flow, and migrations.

## Requirements

- Go `1.26.3`
- Node.js and npm for project scripts
- Docker for local PostgreSQL
- `migrate` CLI for database migrations

## Setup

```powershell
cp .env.example .env
docker compose up -d
npm run migrate:up:all
go run ./cmd
```

The server listens on `http://localhost:8000` by default.

## Environment

See [.env.example](.env.example) for the local template.

Required for the app:

- `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_SSLMODE`
- `CLERK_SECRET_KEY`
- `CLERK_WEBHOOK_SIGNING_SECRET`

Optional:

- `PORT`
- `CORS_ALLOWED_ORIGINS`

## Scripts

```powershell
npm run migrate:up:all   # apply all migrations
npm run migrate:down -- 1 # roll back one migration
npm run sqlc:generate    # regenerate sqlc repository code
npm run test:pretty      # run tests under ./tests with readable output
npm run test:all         # run local, Clerk, E2E, and migration smoke suites
```

## API Surface

Public:

- `POST /webhooks/clerk/user`

Protected under `/api`:

- `GET /me`
- `/characters`
- `/dice-roll/{characterID}`
- `/rooms`
- `GET /rooms/{roomID}/ws`

More detail lives in [docs/testing.md](docs/testing.md), [docs/room-realtime.md](docs/room-realtime.md), and [docs/errors/index.md](docs/errors/index.md).
