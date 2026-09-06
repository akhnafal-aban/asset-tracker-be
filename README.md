# asset-tracker-be

REST API backend for a **Personal Asset Tracker** — track cash, physical items, stocks, and crypto in one place. Built with Go, the standard library `net/http` router (Go 1.22+ path patterns), and SQLite, in a clean `handler → service → repository` layering.

> The companion iOS client (SwiftUI) is in [Asset-Tracker](https://github.com/akhnafal-aban/Asset-Tracker). This repo is the **backend only**.

## Highlights

- **Clean architecture** — `handler → service → repository` layering with constructor-based dependency injection. Each layer has a single responsibility and is independently testable.
- **Repository pattern** — database logic isolated behind `AssetRepository` interface; swapping SQLite for PostgreSQL is a one-struct change.
- **Go 1.22+ routing** — `GET /api/v1/assets/{id}` style patterns via `http.ServeMux`, no third-party router.
- **Graceful shutdown** — listens for `SIGINT`/`SIGTERM`, drains in-flight requests with a 15s timeout, logs cleanly on exit.
- **Structured logging** — `log/slog` with request-scoped fields (method, path, status, duration, ip).
- **Request logging middleware** — wraps `http.ResponseWriter` to capture the status code.
- **Auto-schema + seed** — on first run the DB schema and default categories are created automatically.
- **Ownership checks** — service layer enforces that an asset belongs to the requesting user before returning/mutating it.

## Tech

| | |
|---|---|
| Language | Go 1.22+ (uses `r.PathValue`) |
| HTTP | `net/http` (stdlib) |
| DB | SQLite via [`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) |
| Logging | `log/slog` (stdlib) |
| Config | env vars with sane defaults |

## Project structure

```
cmd/
  api/main.go              # Composition root: wires deps, starts server
internal/
  config/                  # Env-based config
  database/                # SQLite init, schema, seed
  handler/                 # HTTP handlers + JSON response helpers
  middleware/              # Request logger
  model/                   # Domain types (Asset, AssetCategory)
  repository/              # AssetRepository interface + SQLite impl
  server/                  # Server struct, routing, graceful shutdown
  service/                 # Business logic, validation, ownership
```

## API

Base URL: `http://localhost:8080`

| Method | Route | Description |
|---|---|---|
| `GET` | `/api/v1/categories` | List asset categories |
| `GET` | `/api/v1/assets` | List assets for the user |
| `POST` | `/api/v1/assets` | Create an asset |
| `GET` | `/api/v1/assets/{id}` | Get an asset |
| `PUT` | `/api/v1/assets/{id}` | Update an asset |
| `DELETE` | `/api/v1/assets/{id}` | Delete an asset |

### Example

```bash
# List categories
curl http://localhost:8080/api/v1/categories

# Create an asset
curl -X POST http://localhost:8080/api/v1/assets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Emergency Fund",
    "category_id": 1,
    "description": "Cash in savings account",
    "purchase_price": 5000000.00,
    "current_value": 5000000.00
  }'
```

## Getting started

```bash
go run ./cmd/api
```

The server starts on `:8080` and creates `asset_tracker.db` (with schema + seed data) on first run.

### Configuration (env vars)

| Var | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP port |
| `ENV` | `development` | Environment name |
| `DB_PATH` | `asset_tracker.db` | SQLite file path |

## Build

```bash
go build -o bin/api ./cmd/api
./bin/api
```

## Notes

- `defaultUserID = 1` is hardcoded in the handler layer for the MVP — a seeded default user exists. Auth is intentionally out of scope for this iteration.
- `UpdateAsset` uses pointer fields (`*string`, `*float64`) so omitted fields are not zeroed — supports PATCH-style partial updates over PUT.
- The SQLite schema uses `LEFT JOIN asset_categories` so `GET` responses include the category object inline.

## Status

MVP — feature-complete for the tracked scope (CRUD assets + list categories). See `walkthroughMVP.md` for the original build notes.
