.PHONY: run build tidy lint test migrate copy-env

# Copy .env.example to .env if .env does not exist
copy-env:
	@if [ ! -f .env ]; then cp .env.example .env && echo ".env created from .env.example"; fi

# Install / tidy Go modules
tidy:
	go mod tidy

# Build the binary
build:
	go build -o bin/app ./cmd/main.go

# Run the application (copies .env first)
run: copy-env
	go run ./cmd/main.go

# Run all tests
test:
	go test ./... -v -race -cover

# Lint (requires golangci-lint)
lint:
	golangci-lint run ./...
