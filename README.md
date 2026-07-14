<p align="center">
  <img src="assets/readme-logo-esoteric.webp" alt="Miskatonic Lab" width="720">
</p>

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.26.3-00ADD8?logo=go&logoColor=white">
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white">
  <img alt="Docker" src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white">
  <img alt="Clerk" src="https://img.shields.io/badge/Auth-Clerk-6C47FF?logo=clerk&logoColor=white">
</p>

# Miskatonic Lab Backend

Go backend для листов персонажей Call of Cthulhu, комнат, бросков кубов и WebSocket-чата. PostgreSQL хранит данные, Clerk отвечает за аутентификацию и пользовательские webhook-события, а `sqlc` генерирует типизированный слой запросов.

## Возможности

- Листы персонажей: характеристики, навыки, состояния, предыстория, финансы, заметки и портреты.
- Броски кубов с историей и необязательным контекстом комнаты.
- Комнаты: участники, роли, выбранные персонажи, история событий и WebSocket-чат.
- Защищённые API-маршруты и Clerk webhook `POST /webhooks/clerk/user`.

## Требования

- Go `1.26.3`
- Node.js и npm
- Docker Desktop
- `migrate` CLI
- `.env`, созданный на основе [.env.example](.env.example)

`npm run test:setup` самостоятельно скачивает Dev Tunnel CLI при его отсутствии; `winget` не требуется.

## Запуск приложения

### Supabase / внешний PostgreSQL

```powershell
Copy-Item .env.example .env
# Заполни DATABASE_URL и Clerk runtime values в .env.
npm run migrate:up:all
go run ./cmd
```

Для Supabase используй `sslmode=require`. Runtime не должен использовать Transaction Pooler: `pgx`/`sqlc` работают с prepared statements.

### Локальный PostgreSQL

```powershell
Copy-Item .env.example .env
docker compose up -d postgres
npm run migrate:up:all
go run ./cmd
```

По умолчанию API слушает `http://localhost:8000`.

## Тесты

Главная команда полного прогона:

```powershell
npm run test:all
```

Она валидирует `.env`, поднимает Docker test DB, пересоздаёт disposable migration-smoke DB, запускает локальные Go-тесты, временный backend и Dev Tunnel, выполняет реальные Clerk/E2E проверки, migration smoke и останавливает созданные backend/tunnel процессы.

### Однократная настройка live-тестов

1. Заполни runtime, test DB и два `E2E_TEST*_MAIL` в `.env`.
2. Выполни:

   ```powershell
   npm run test:setup
   ```

   Команда запросит Dev Tunnel login при необходимости, создаст persistent anonymous tunnel на `8001` и заполнит `DEVTUNNEL_TUNNEL_ID` с `CLERK_WEBHOOK_PUBLIC_URL`.

3. Один раз укажи напечатанный адрес `<CLERK_WEBHOOK_PUBLIC_URL>/webhooks/clerk/user` в Clerk Dashboard для `user.created`, `user.updated` и `user.deleted`.

После этого для полного набора достаточно одной команды: `npm run test:all`.

### Команды

| Команда | Назначение |
| --- | --- |
| `npm run test:setup` | Подготовить и проверить Dev Tunnel, test DB и обязательный env-контракт. |
| `npm run test:local` | Поднять/подготовить Docker test DB и запустить `go test ./...`. |
| `npm run test:e2e` | Запустить изолированный backend и live HTTP/WebSocket E2E. |
| `npm run test:clerk` | Запустить реальную Clerk webhook integration suite через Dev Tunnel. |
| `npm run test:migrations` | Сбросить disposable smoke DB и проверить migration `up/down/up`. |
| `npm run test:all` | Запустить все suites в воспроизводимом порядке. |
| `npm run test:pretty` | Детальный вывод по отдельным локальным тестам; `test:all` намеренно выводит только краткий итог этапов. |

`TEST_DATABASE_URL` разрешён только для loopback DB с именем `*_test`. `MIGRATION_SMOKE_DATABASE_URL` должен быть отдельной loopback БД с именем `*_migration_smoke_test`; перед smoke-тестом она удаляется и создаётся заново. Основной Supabase/runtime database runner не затрагивает.

Полный env-контракт, порядок, cleanup и разбор ошибок: [docs/testing.md](docs/testing.md).

## Остальные команды

```powershell
npm run migrate:up:all   # применить все migrations
npm run migrate:version  # показать версию migrations
npm run migrate:down -- 1 # откатить одну migration
npm run sqlc:generate    # обновить сгенерированный sqlc-код
npm run test:pretty      # читаемый вывод тестов
```

## API

Публичный маршрут:

- `POST /webhooks/clerk/user`

Защищённые маршруты находятся под `/api`:

- `GET /me`
- `/characters`
- `/dice-roll/{characterID}`
- `/rooms`
- `GET /rooms/{roomID}/ws`

Подробности: [docs/testing.md](docs/testing.md), [docs/room-realtime.md](docs/room-realtime.md), [docs/errors/index.md](docs/errors/index.md).
