.DEFAULT_GOAL := help

CONTAINER         = restaurant-assignment-postgres
POSTGRES_IMAGE    ?= postgres:17
POSTGRES_PORT     ?= 5432
POSTGRES_USER     ?= samu
POSTGRES_PASSWORD ?= samu
POSTGRES_DB       ?= samu
DATABASE_SCRIPT   ?= scripts/database.sql

.PHONY: help build migrate unit-test run all

help:
	@echo "Restaurant Assignment"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  help       Show this help message"
	@echo "  build      Create and start the PostgreSQL Docker container"
	@echo "  migrate    Load the database schema into PostgreSQL"
	@echo "  run        Start the application on :8080"
	@echo "  unit-test  Run unit tests with coverage report"
	@echo "  all        Run tests, setup the database, and start the application"

build:
	@if docker ps -aq -f name=^$(CONTAINER)$$ | grep -q .; then \
		echo "Container $(CONTAINER) already exists. Remove it first with: docker rm -f $(CONTAINER)"; \
		exit 1; \
	fi
	@echo "Pulling $(POSTGRES_IMAGE)..."
	@docker pull $(POSTGRES_IMAGE)
	@echo "Creating and starting container $(CONTAINER)..."
	@docker run -d \
		--name $(CONTAINER) \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p $(POSTGRES_PORT):5432 \
		$(POSTGRES_IMAGE) >/dev/null
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker exec $(CONTAINER) pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) >/dev/null 2>&1; do sleep 1; done
	@echo "Container $(CONTAINER) is running."

migrate:
	@if ! docker ps -aq -f name=^$(CONTAINER)$$ | grep -q .; then \
		echo "Container $(CONTAINER) does not exist. Run 'make build' first."; \
		exit 1; \
	fi
	@if [ "$$(docker inspect -f '{{.State.Running}}' $(CONTAINER))" != "true" ]; then \
		echo "Container $(CONTAINER) is not running. Start it with: docker start $(CONTAINER)"; \
		exit 1; \
	fi
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker exec $(CONTAINER) pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) >/dev/null 2>&1; do sleep 1; done
	@echo "Loading $(DATABASE_SCRIPT)..."
	@docker exec -i $(CONTAINER) psql -v ON_ERROR_STOP=1 -U $(POSTGRES_USER) -d $(POSTGRES_DB) < $(DATABASE_SCRIPT)
	@echo "Database schema loaded."

run:
	@echo "Starting application on :8080..."
	@go run ./cmd/

unit-test:
	@set -e; \
	status=0; \
	go test ./... -count=1 -race -v -coverprofile=coverage.out -covermode=atomic || status=$$?; \
	if [ -f coverage.out ]; then \
		echo ""; \
		echo "Coverage report (below 100%):"; \
		go tool cover -func=coverage.out | grep -Ev '[[:space:]]100\.0%$$' || true; \
		rm -f coverage.out coverage.html; \
	fi; \
	echo "Unit tests completed."; \
	exit $$status

all: unit-test build migrate run