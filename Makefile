.PHONY: run build tidy lint test

# Install / tidy Go modules
tidy:
	go mod tidy

# Build the binary
build:
	go build -o bin/app ./internal/cmd/main.go

# Run the application
run:
	go run ./internal/cmd/main.go

# Run all tests
test:
	go test ./... -v -race -cover

# Lint (requires golangci-lint)
lint:
	golangci-lint run ./...
