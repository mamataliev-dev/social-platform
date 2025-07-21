include .env
export

# Settings
APP_NAME := user-service
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

# Docker
docker-build:
	@echo "Building Docker image for user-service..."
	@docker build -t user-service:latest -f Dockerfile ../..

docker-run:
	@echo "Stopping and removing existing container if it exists..."
	-@docker stop user-service && docker rm user-service
	@echo "Running new user-service container..."
	@docker run -d \
	  --name user-service \
	  --env-file .env \
	  -p 100:100 \
	  -p 50051:50051 \
	  user-service:latest

docker-stop:
	@echo "Stopping user-service container..."
	@docker stop user-service

docker-compose-up:
	docker compose up -d

docker-compose-down:
	docker compose down