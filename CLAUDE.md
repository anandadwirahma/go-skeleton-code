# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Scope Discipline
- When asked to generate a single file (e.g., one usecase, one handler), generate ONLY that file.
- Do NOT proactively create supporting files, tests, or related boilerplate unless explicitly requested.
- If you think additional files are needed, list them and ask before creating. File to generate: [describe here]
- Ask before expanding scope beyond the literal request.

## Commands

```bash
make run       # Start the application (go run ./internal/cmd/main.go)
make build     # Compile binary to bin/app
make test      # Run all tests with race detector and coverage
make lint      # Run golangci-lint
make tidy      # Sync go.mod/go.sum
```

Run a single test:
```bash
go test -v -run TestFunctionName ./internal/...
```

## Architecture

This project follows **Clean Architecture** with these layers (outer → inner):

```
HTTP Request → Controller → Usecase → Repository → PostgreSQL
```

| Layer | Path | Responsibility |
|---|---|---|
| Delivery | `internal/delivery/http/` | HTTP handlers (Fiber v3), request binding, response formatting |
| Usecase | `internal/usecase/` | Business logic; depends only on repository interfaces |
| Repository | `internal/repository/` | Data access via GORM; generic base `Repository[T]` for reusable CRUD |
| Entity | `internal/entity/` | GORM structs (database schema) |
| Model | `internal/model/` | Request/Response DTOs, validators, `WebResponse[T]` generic wrapper |
| Config | `internal/config/` | Dependency wiring (DI container); initializes Viper, Logrus, GORM, Fiber, Validator |
| Gateway | `internal/gateway/http/` | Outbound HTTP client (go-resty); `HTTPGateway` struct with GET/POST/PUT/PATCH/DELETE methods |

**Entry point:** `internal/cmd/main.go` calls `config/app.go` which wires all dependencies.

**Generic Repository:** `internal/repository/repository.go` provides a base `Repository[T any]` struct with `Create`, `Update`, `Delete`, `FindById` — domain repositories embed this and add domain-specific queries.

**Converters:** `internal/model/converter/` contains functions that transform GORM entities into response models — keep entity↔model mapping here.

**All HTTP responses** use `model.WebResponse[T]` for consistency. Error handling is centralized in `internal/config/fiber.go`.

**HTTP Gateway** (`internal/gateway/http`, package `httpgateway`): outbound HTTP client built on go-resty. Construct with `httpgateway.New(opts...)`. Options: `WithBaseURL`, `WithTimeout`, `WithRetryCount`, `WithRetryWaitTime`, `WithHeaders`. All methods return `(*Response, error)`; non-2xx responses return `*HTTPError` — inspect with `httpgateway.AsHTTPError(err)`.

## Go Project Conventions
- This project uses Go with clean architecture (usecase/handler/repository layering).
- Use go-resty for HTTP clients.
- When scaffolding, follow the existing domain structure.

## Configuration

Runtime config is read from `config.json` in the project root via Viper. Key fields:

```json
{
  "app": { "name": "..." },
  "web": { "port": 3000, "prefork": false },
  "log": { "level": 5 },
  "database": { "host", "port", "user", "password", "name", "pool" }
}
```

Viper also checks the parent directory for `config.json`, so it works when invoked from `internal/cmd/`.

## Adding a New Domain

To scaffold a new domain (e.g., `User`), create these files following the `example_*` naming pattern:

1. `internal/entity/user_entity.go` — GORM struct; add to `AutoMigrate` in `internal/config/gorm.go`
2. `internal/model/user_model.go` — Request/Response DTOs with `validate` tags
3. `internal/model/converter/user_converter.go` — Entity → Model conversion
4. `internal/repository/user_repository.go` — Embed `Repository[entity.User]`, add domain queries
5. `internal/usecase/user_usecase.go` — Define interface + implementation
6. `internal/delivery/http/user_controller.go` — Fiber handlers
7. Register routes in `internal/delivery/http/router/route/route.go`
8. Wire dependencies in `internal/config/app.go`