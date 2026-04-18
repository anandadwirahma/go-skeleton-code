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
Delivery → Usecase → Repository → Database
```

- **Entity** represents the database schema.
- **Model** defines the payload for requests and responses.
- **Repository** implements persistence using GORM.
- **Usecase** orchestrates business rules.
- **Delivery** handles HTTP concerns using Fiber (controllers).
- **Config** wires concrete implementations together and initializes libraries.

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
