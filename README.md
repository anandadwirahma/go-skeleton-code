# go-skeleton-code

A **production-ready Golang REST API skeleton** using **Clean Architecture**, built with Fiber v3, GORM, PostgreSQL, Viper, and Logrus.

---

## Folder Structure

```
go-skeleton-code/
├── internal/
│   ├── cmd/                             # Entry point (main.go)
│   ├── config/                          # Viper, Fiber, GORM, Logrus, and Validator setup
│   ├── delivery/                        # Delivery Layer (Fiber controllers, router)
│   │   └── http/
│   │       ├── router/                  # Route registration
│   │       └── example_controller.go
│   ├── entity/                          # Domain entities (GORM structs)
│   ├── gateway/                         # Outbound HTTP clients (go-resty)
│   │   └── http/                        # HTTPGateway — GET/POST/PUT/PATCH/DELETE
│   ├── model/                           # Request/Response models and converters
│   │   └── converter/
│   ├── repository/                      # Database operations (GORM implementations)
│   └── usecase/                         # Business logic orchestration
├── test/                                # Testing files (e.g., HTTP request tests)
├── config.json                          # Application configuration file
├── Makefile                             # Make commands for build, run, and test
└── go.mod
```

---

## Tech Stack

| Layer          | Technology                      |
|----------------|---------------------------------|
| Language       | Go 1.25+                        |
| HTTP Framework | Fiber v3                        |
| ORM            | GORM                            |
| Database       | PostgreSQL                      |
| Logging        | Logrus                          |
| Config         | Viper                           |
| Validation     | Validator v10                   |
| HTTP Client    | go-resty v2                     |

---

## Prerequisites

- Go 1.25+
- PostgreSQL running locally (or via Docker)
- `make` (optional but recommended)

---

## Quick Start

```bash
# 1. Clone the repo
git clone https://github.com/yourusername/go-skeleton-code.git
cd go-skeleton-code

# 2. Configure environment
# Edit config.json with your PostgreSQL credentials and app settings

# 3. Create the database (if it doesn't exist yet)
psql -h localhost -U postgres -c "CREATE DATABASE skeleton_db"

# Or, if PostgreSQL is running inside Docker:
# docker exec -it <postgres-container-name> psql -U postgres -c "CREATE DATABASE skeleton_db;"

# 4. Install dependencies
make tidy

# 5. Run the application
make run
# Or manually: go run ./internal/cmd/main.go
```

The server will start on `http://localhost:3000`.

---

## API Endpoints

Base URL: `http://localhost:3000/api`

| Method | Path             | Description           |
|--------|------------------|-----------------------|
| POST   | `/example`       | Create an example     |

---

## Curl Examples

```bash
# Create an example
curl -X POST http://localhost:3000/api/example \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Nanda",
    "email": "nanda@gmail.com"
  }'
```

---

## Configuration

This project uses `config.json` managed by Viper. Ensure `config.json` is present in the root directory.

```json
{
    "app": {
        "name": "go-skeleton-code"
    },
    "web": {
        "prefork": false,
        "port": 3000
    },
    "log": {
        "level": 6
    },
    "database": {
        "username": "postgres",
        "password": "postgres",
        "host": "localhost",
        "port": 5432,
        "name": "skeleton_db",
        "pool": {
            "idle": 10,
            "max": 100,
            "lifetime": 300
        }
    }
}
```

---

## Architecture

This project follows **Clean Architecture**:

```
HTTP Request → Controller → Usecase → Repository → PostgreSQL
                                ↓
                          Gateway (outbound HTTP)
```

- **Entity** represents the database schema.
- **Model** defines the payload for requests and responses.
- **Repository** implements persistence using GORM. A generic `Repository[T]` base provides `Create`, `Update`, `Delete`, and `FindById` — domain repositories embed it and add domain-specific queries.
- **Usecase** orchestrates business rules.
- **Delivery** handles HTTP concerns using Fiber (controllers).
- **Gateway** wraps outbound HTTP calls via go-resty (`internal/gateway/http`, package `httpgateway`).
- **Config** wires concrete implementations together and initializes libraries.

---

## HTTP Gateway

`HTTPGateway` is an outbound HTTP client built on go-resty. Use it to call external services from your usecases.

### Construction

```go
gw := httpgateway.New(
    httpgateway.WithBaseURL("https://api.example.com"),
    httpgateway.WithTimeout(10 * time.Second),
    httpgateway.WithRetryCount(3),
    httpgateway.WithRetryWaitTime(2 * time.Second),
    httpgateway.WithHeaders(map[string]string{"Authorization": "Bearer <token>"}),
)
```

Defaults (applied when no option overrides): 30s timeout, 3 retries, 2s retry wait, `Content-Type: application/json` and `Accept: application/json`.

### Methods

All methods share the same signature pattern and return `(*Response, error)`.

```go
// GET
resp, err := gw.Get(ctx, "/users/1", headers, queryParams, &result)

// POST
resp, err := gw.Post(ctx, "/users", headers, body, &result)

// PUT
resp, err := gw.Put(ctx, "/users/1", headers, body, &result)

// PATCH
resp, err := gw.Patch(ctx, "/users/1", headers, body, &result)

// DELETE
resp, err := gw.Delete(ctx, "/users/1", headers, &result)
```

Pass `nil` for `headers`, `queryParams`, or `result` when not needed.

### Error Handling

Non-2xx responses return `*HTTPError`. Use `httpgateway.AsHTTPError` to inspect the status code:

```go
resp, err := gw.Get(ctx, "/resource", nil, nil, &result)
if err != nil {
    if httpErr, ok := httpgateway.AsHTTPError(err); ok {
        // httpErr.StatusCode, httpErr.Body, httpErr.URL
    }
    return err
}
```

---

## Extending the Project

To add a new feature (e.g., `user`), follow this workflow:

1. **Entity**: Create `internal/entity/user_entity.go` for the database schema.
2. **Model**: Create `internal/model/user_model.go` for request/response bodies.
3. **Converter**: Create `internal/model/converter/user_converter.go` to convert between Entity and Model.
4. **Repository**: Create `internal/repository/user_repository.go` for database operations.
5. **Usecase**: Create `internal/usecase/user_usecase.go` for business logic.
6. **Controller**: Create `internal/delivery/http/user_controller.go` to handle HTTP requests.
7. **Route**: Register the route group in `internal/delivery/http/router/route/route.go`.
8. **Wiring**: Wire the dependencies and register the controller in `internal/config/app.go`.

To call an external service, inject an `*httpgateway.HTTPGateway` into your usecase and use the appropriate method.
