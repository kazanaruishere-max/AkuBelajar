.PHONY: dev-api run-api build test lint migrate-up migrate-down seed db-check clean

# ---- Development ----

# Run API server with hot-reload (requires air: go install github.com/air-verse/air@latest)
dev-api:
	air -c .air.toml

# Run API server without hot-reload
run-api:
	go run ./cmd/api

# Build binary
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server ./cmd/api

# ---- Testing ----

# Run all tests with coverage
test:
	go test ./... -v -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out

# Run tests in short mode
test-short:
	go test ./... -short -count=1

# ---- Linting ----

lint:
	golangci-lint run ./...

# ---- Database ----

# Run all migrations up
migrate-up:
	migrate -path migrations -database "$${DATABASE_URL}" up

# Rollback last migration
migrate-down:
	migrate -path migrations -database "$${DATABASE_URL}" down 1

# Create new migration
# Usage: make migrate-create name=create_users
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Seed development data
seed:
	psql "$${DATABASE_URL}" -f migrations/seed.sql

# Check database connection
db-check:
	@psql "$${DATABASE_URL}" -c "SELECT 1 AS connected;" && echo "✅ Database connected" || echo "❌ Database NOT connected"

# ---- Docker ----

# Start infrastructure (PostgreSQL + Redis)
infra-up:
	docker compose up -d postgres redis

# Stop infrastructure
infra-down:
	docker compose down

# ---- Cleanup ----

clean:
	rm -rf bin/ coverage.out tmp/
