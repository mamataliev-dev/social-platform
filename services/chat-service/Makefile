include .env
export

# Settings
APP_NAME := chat-service
MAIN := cmd/server/main.go
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_DIR := migrations

# Run
run:
	go run $(MAIN)

# Build
build:
	go build -o bin/$(APP_NAME) $(MAIN)

# Migrations
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
	# Example: make migrate-create name=create_user_table

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

# Tests
test:
	go test ./internal/... -v

# Formatting and tidy
fmt:
	go fmt ./...
	go vet ./...