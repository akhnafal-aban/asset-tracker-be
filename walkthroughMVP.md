# Walkthrough: Personal Asset Tracker MVP

The backend foundation for the Personal Asset Tracker has been successfully implemented according to our plan!

## What Was Completed

1. **Database & Schema Initialization**
   - Configured `mattn/go-sqlite3` and added a `database` package to handle the SQLite connection.
   - Automatically initializes the database schema (`asset_tracker.db`) if it doesn't exist upon server startup.
   - Seeded a default user (`user_id = 1`) and default asset categories (`Cash`, `Physical`, `Stocks`, `Crypto`, `Other`) so you can test the API immediately.

2. **Domain Models & Repository Layer**
   - Defined the core data structures in `internal/model/asset.go`.
   - Created `internal/repository/asset_repository.go` implementing the Repository pattern with SQLite. This decouples the database logic from the rest of the application, making a future migration to PostgreSQL very simple.

3. **Service & Business Logic**
   - Implemented `internal/service/asset_service.go` which enforces basic validation rules (e.g., prices cannot be negative, assets must belong to the hardcoded default user).

4. **HTTP Handlers & Routing**
   - Created RESTful JSON handlers in `internal/handler/asset.go`.
   - Updated `internal/server/server.go` to register these new API routes.
   - Wired all the dependencies up in `cmd/api/main.go`.

## Implemented Endpoints

The following REST APIs are now available under `http://localhost:8080`:

| Method | Route | Description |
|---|---|---|
| `GET` | `/health` | Application healthcheck |
| `GET` | `/api/v1/categories` | List seeded asset categories |
| `GET` | `/api/v1/assets` | List all assets for the user |
| `POST` | `/api/v1/assets` | Create a new asset |
| `GET` | `/api/v1/assets/{id}` | Retrieve a specific asset |
| `PUT` | `/api/v1/assets/{id}` | Update a specific asset |
| `DELETE` | `/api/v1/assets/{id}` | Delete a specific asset |

## Verification

- The project builds cleanly with `go build -o bin/api ./cmd/api`.
- You can now run the application using: `go run ./cmd/api`
- Once running, you can test it by opening a new terminal and curling the endpoints. For example, to list categories:
  ```bash
  curl http://localhost:8080/api/v1/categories
  ```

And to create a test asset:
```bash
curl -X POST http://localhost:8080/api/v1/assets \
-H "Content-Type: application/json" \
-d '{
  "name": "My Emergency Fund",
  "category_id": 1,
  "description": "Cash in savings account",
  "purchase_price": 5000.00,
  "current_value": 5000.00
}'
```

The foundation is fully built! You can now start experimenting with these APIs in your SwiftUI application. Let me know if you would like me to adjust any models or endpoints!
