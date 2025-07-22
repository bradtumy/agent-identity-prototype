# Makefile for Agent Identity PoC

APP_NAME := agent-identity-poc
COMPOSE_FILE := docker-compose.yml

# Targets

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build         Build Go services"
	@echo "  make run           Run Go services (locally)"
	@echo "  make lint          Run golangci-lint"
	@echo "  make test          Run unit tests"
	@echo "  make docker-up     Start all containers"
	@echo "  make docker-down   Stop all containers and volumes"
	@echo "  make restart       Rebuild and restart Docker environment"
	@echo "  make logs          Show logs from all containers"

broker:
	go run ./broker

# Go targets
build:
	@echo "Building Go services..."
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	@echo "Running locally..."
	go run ./cmd/server

lint:
	@echo "Linting..."
	golangci-lint run

test:
	@echo "Running tests..."
	go test ./... -v

# Docker targets
docker-up:
	docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down -v

restart:
	docker compose -f $(COMPOSE_FILE) down -v
	docker compose -f $(COMPOSE_FILE) up --build -d

logs:
	docker compose logs -f --tail=100

