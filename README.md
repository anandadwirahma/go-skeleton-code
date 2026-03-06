# go-skeleton-code

A **production-ready Golang REST API skeleton** using **Clean Architecture**, built with Gin, GORM, PostgreSQL, and Zap.

---

## Folder Structure

```
go-skeleton-code/
├── cmd/
│   └── main.go                          # Entry point — manual DI & graceful shutdown
├── internal/
│   ├── domain/                          # Domain Layer (entities, interfaces, errors)
│   │   └── contact/
│   │       ├── entity.go
│   │       ├── repository.go            # Repository interface
│   │       ├── usecase.go               # Usecase interface + input types
│   │       └── errors.go                # Sentinel errors
│   ├── repository/                      # Repository Layer (GORM implementations)
│   │   └── contact/
│   │       └── postgres_repository.go
│   ├── usecase/                         # Usecase Layer (business logic)
│   │   └── contact/
│   │       └── usecase.go
│   ├── delivery/                        # Delivery Layer (Gin handlers, DTOs, router)
│   │   └── http/
│   │       ├── dto/
│   │       │   └── contact_dto.go
│   │       ├── handler/
│   │       │   └── contact_handler.go
│   │       ├── middleware/
│   │       │   └── logger.go
│   │       └── router/
│   │           └── router.go
│   └── infrastructure/                 # Infrastructure Layer
│       ├── config/
│       │   └── config.go               # Env config loader
│       ├── database/
│       │   └── postgres.go             # GORM + auto-migrate
│       ├── logger/
│       │   └── logger.go               # Zap logger
│       └── httpclient/
│           └── client.go               # External HTTP client
├── .env.example
├── Makefile
└── go.mod
```

---

## Tech Stack

| Layer          | Technology                      |
|----------------|---------------------------------|
| Language       | Go 1.22+                        |
| HTTP Framework | Gin                             |
| ORM            | GORM                            |
| Database       | PostgreSQL                      |
| Logging        | Uber Zap (structured)           |
| Config         | godotenv + os.Getenv            |
| HTTP Client    | net/http (wrapped)              |

---

## Prerequisites

- Go 1.22+
- PostgreSQL running locally (or via Docker)
- `make` (optional but recommended)

---

## Quick Start

```bash
# 1. Clone the repo
git clone https://github.com/yourusername/go-skeleton-code.git
cd go-skeleton-code

# 2. Copy and configure environment
cp .env.example .env
# Edit .env with your PostgreSQL credentials

# 3. Create the database (if it doesn't exist yet)
psql -h localhost -U postgres -c "CREATE DATABASE skeleton_db"


# Or, if PostgreSQL is running inside Docker:
# docker exec -it <postgres-container-name> psql -U postgres -c "CREATE DATABASE skeleton_db;"

# 4. Install dependencies
make tidy

# 5. Run the application
make run
```

The server will start on `http://localhost:8080`.

---

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

| Method | Path             | Description           |
|--------|------------------|-----------------------|
| GET    | `/health`        | Health check          |
| POST   | `/contacts`      | Create a contact      |
| GET    | `/contacts`      | List all contacts     |
| GET    | `/contacts/:id`  | Get contact by ID     |
| PUT    | `/contacts/:id`  | Update a contact      |
| DELETE | `/contacts/:id`  | Delete a contact      |

---

## Curl Examples

```bash
# Health check
curl http://localhost:8080/health

# Create a contact
curl -X POST http://localhost:8080/api/v1/contacts \
  -H "Content-Type: application/json" \
  -d '{
    "name":    "Alice Smith",
    "email":   "alice@example.com",
    "message": "Hello, I would like to get in touch with your team."
  }'

# List all contacts
curl http://localhost:8080/api/v1/contacts

# Get contact by ID
curl http://localhost:8080/api/v1/contacts/1

# Update a contact
curl -X PUT http://localhost:8080/api/v1/contacts/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name":    "Alice Updated",
    "email":   "alice.updated@example.com",
    "message": "Updated message with sufficient length here."
  }'

# Delete a contact
curl -X DELETE http://localhost:8080/api/v1/contacts/1
```

---

## Response Format

All endpoints return a consistent JSON envelope:

**Success:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error:**
```json
{
  "success": false,
  "error": "descriptive error message"
}
```

---

## Architecture

This project follows **Clean Architecture** (also known as Hexagonal / Ports & Adapters):

```
Delivery → Usecase → Domain ← Repository
              ↑                    ↑
         (interface)          (interface)
```

- **Domain** knows nothing about HTTP, databases, or frameworks.
- **Usecase** orchestrates business rules using domain interfaces.
- **Repository** implements persistence using GORM.
- **Delivery** handles HTTP concerns (validation, serialization).
- **Infrastructure** wires concrete implementations together.

Dependencies always point **inward** — outer layers depend on inner layers, never the reverse.

---

## Environment Variables

| Variable                | Default                               | Description                    |
|-------------------------|---------------------------------------|--------------------------------|
| `APP_ENV`               | `development`                         | `development` or `production`  |
| `SERVER_PORT`           | `8080`                                | HTTP server port               |
| `DB_HOST`               | `localhost`                           | PostgreSQL host                |
| `DB_PORT`               | `5432`                                | PostgreSQL port                |
| `DB_USER`               | `postgres`                            | Database user                  |
| `DB_PASSWORD`           | *(empty)*                             | Database password              |
| `DB_NAME`               | `skeleton_db`                         | Database name                  |
| `DB_SSLMODE`            | `disable`                             | SSL mode                       |
| `DB_MAX_OPEN_CONNS`     | `25`                                  | Max open DB connections        |
| `DB_MAX_IDLE_CONNS`     | `5`                                   | Max idle DB connections        |
| `LOG_LEVEL`             | `info`                                | `debug` / `info` / `warn` / `error` |
| `EXTERNAL_API_BASE_URL` | `https://jsonplaceholder.typicode.com`| External API base URL          |
| `EXTERNAL_API_TIMEOUT_SEC` | `10`                               | External API timeout (seconds) |

---

## Extending the Project

To add a new feature (e.g. `user`), follow this checklist:

1. Create `internal/domain/user/entity.go` — struct + table name
2. Create `internal/domain/user/repository.go` — interface
3. Create `internal/domain/user/usecase.go` — interface + inputs
4. Create `internal/domain/user/errors.go` — sentinel errors
5. Create `internal/repository/user/postgres_repository.go` — GORM impl
6. Create `internal/usecase/user/usecase.go` — business logic
7. Create `internal/delivery/http/dto/user_dto.go` — request/response DTOs
8. Create `internal/delivery/http/handler/user_handler.go` — Gin handlers
9. Register route group in `internal/delivery/http/router/router.go`
10. Register model in `internal/infrastructure/database/postgres.go` → `autoMigrate`
11. Wire dependencies in `cmd/main.go`
