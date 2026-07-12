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

## Testing

### Стек тестирования

| Инструмент | Назначение |
|---|---|
| `testing` (stdlib) | Каркас тестов |
| `testify` (require/assert) | Утверждения |
| `httptest` (stdlib) | API-тесты хендлеров |
| `gotestsum` | Форматированный вывод тестов |
| Hand-written fakes | Заглушки на границах пакетов (вместо mock-библиотек) |
| PostgreSQL (Docker) | Реальная БД для интеграционных тестов |

### Проведённое тестирование

| Уровень | Что покрывает | Кол-во тестов |
|---|---|---|
| **Модульные** | Бизнес-логика сервисов, хендлеры, валидация, мапперы, парсеры, middleware, утилиты, WebSocket-хаб | ~450 |
| **Интеграционные** | `sqlc`-запросы, constraints, foreign keys, upserts, каскады, owner-scoping, миграции — на реальной PostgreSQL | ~250 |
| **End-to-End** | HTTP + WebSocket поверх реального сервера и БД (опционально, с реальным Clerk-токеном) | 7 |
| **Clerk Integration** | Реальный Clerk API → webhook → локальная БД (опционально) | 1 |
| **Migration Smoke** | `migrate up/down` против disposable-БД (опционально) | 1 |

**Всего: 710 тестовых сценариев** — все `pass`. Детальная трасса по доменам:

| Домен | Кол-во |
|---|---|
| Персонажи | 385 |
| Комнаты | 102 |
| Броски кубов | 73 |
| Пользователи | 49 |
| События | 21 |
| Модели | 20 |
| WebSocket | 11 |
| Слушатели | 11 |
| Утилиты | 9 |
| Сквозные проверки | 7 |
| Middleware | 7 |
| Наблюдаемость | 7 |
| HTTP-адаптер | 4 |
| Конфигурация | 3 |
| Миграции | 1 |

Полная тестовая трасса — [docs/test-trace.md](docs/test-trace.md). Подробнее о стратегии — [docs/testing.md](docs/testing.md).

## Requirements

- Go `1.26.3`
- Node.js and npm for project scripts
- Docker for local PostgreSQL
- `migrate` CLI for database migrations

## Setup

### Supabase

```powershell
cp .env.example .env
# Set DATABASE_URL in .env to your Supabase Direct connection or Session pooler URL.
npm run migrate:up:all
go run ./cmd
```

Use `sslmode=require` in the Supabase URL. Do not use the Transaction pooler URL for the backend runtime; transaction pooling does not support prepared statements, while the Go repository layer uses `pgx`/`sqlc`.

### Local PostgreSQL fallback

```powershell
cp .env.example .env
docker compose up -d
npm run migrate:up:all
go run ./cmd
```

### Local database tests

```powershell
npm run test:db
```

This starts isolated PostgreSQL on port `5433`, applies migrations, seeds deterministic data, then runs database tests. Tests use `TEST_DATABASE_URL` only; Supabase and database names without `_test` are rejected.

The server listens on `http://localhost:8000` by default.

## Environment

See [.env.example](.env.example) for the local template.

Required for the app:

- `DATABASE_URL` for Supabase or `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_SSLMODE` for local PostgreSQL
- `CLERK_SECRET_KEY`
- `CLERK_WEBHOOK_SIGNING_SECRET`

Optional:

- `PORT`
- `CORS_ALLOWED_ORIGINS`

## Scripts

```powershell
npm run migrate:up:all   # apply all migrations
npm run migrate:version  # print current migration version
npm run migrate:down -- 1 # roll back one migration
npm run testdb:prepare   # start, migrate, and seed local test DB
npm run test:db          # prepare local DB and run database tests
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

More detail lives in [docs/testing.md](docs/testing.md), [docs/test-trace.md](docs/test-trace.md), [docs/room-realtime.md](docs/room-realtime.md), and [docs/errors/index.md](docs/errors/index.md).
