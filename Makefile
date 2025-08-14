.PHONY: help build run test clean docker-build docker-up docker-down docker-logs docker-shell db-migrate db-seed

# Variables
APP_NAME = gonexttask
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_DEV = docker-compose -f docker-compose.dev.yml
GO = go
GOFLAGS = -v

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Local development
build: ## Build the application locally
	$(GO) build $(GOFLAGS) -o bin/api cmd/api/main.go

run: ## Run the application locally
	$(GO) run cmd/api/main.go

test: ## Run tests
	$(GO) test $(GOFLAGS) ./...

clean: ## Clean build artifacts
	rm -rf bin/
	$(GO) clean

mod-download: ## Download Go modules
	$(GO) mod download

mod-tidy: ## Tidy Go modules
	$(GO) mod tidy

# Docker commands for production
docker-build: ## Build Docker images
	$(DOCKER_COMPOSE) build

docker-up: ## Start all services in production mode
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop all services
	$(DOCKER_COMPOSE) down

docker-restart: ## Restart all services
	$(DOCKER_COMPOSE) restart

docker-logs: ## View logs from all services
	$(DOCKER_COMPOSE) logs -f

docker-logs-api: ## View logs from API service only
	$(DOCKER_COMPOSE) logs -f api

docker-ps: ## Show running containers
	$(DOCKER_COMPOSE) ps

docker-clean: ## Remove all containers and volumes
	$(DOCKER_COMPOSE) down -v

# Docker commands for development
dev-up: ## Start all services in development mode with hot reload
	$(DOCKER_COMPOSE_DEV) up

dev-up-d: ## Start all services in development mode (detached)
	$(DOCKER_COMPOSE_DEV) up -d

dev-down: ## Stop development services
	$(DOCKER_COMPOSE_DEV) down

dev-restart: ## Restart development services
	$(DOCKER_COMPOSE_DEV) restart

dev-logs: ## View development logs
	$(DOCKER_COMPOSE_DEV) logs -f

dev-build: ## Rebuild development images
	$(DOCKER_COMPOSE_DEV) build

dev-shell: ## Enter the API container shell
	$(DOCKER_COMPOSE_DEV) exec api sh

# Database commands
db-shell: ## Enter PostgreSQL shell
	$(DOCKER_COMPOSE) exec postgres psql -U postgres -d gonexttask

db-migrate: ## Run database migrations
	$(DOCKER_COMPOSE) exec postgres psql -U postgres -d gonexttask -f /docker-entrypoint-initdb.d/001_create_tables.sql

db-backup: ## Backup database
	@mkdir -p backups
	$(DOCKER_COMPOSE) exec postgres pg_dump -U postgres gonexttask > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql

db-restore: ## Restore database from latest backup
	@read -p "Enter backup filename: " file; \
	$(DOCKER_COMPOSE) exec -T postgres psql -U postgres -d gonexttask < backups/$$file

# Quick start commands
quick-start: ## Quick start for development (build and run with Docker)
	@echo "Starting GoNexttask in development mode..."
	@make dev-up

quick-start-prod: ## Quick start for production
	@echo "Starting GoNexttask in production mode..."
	@make docker-build
	@make docker-up
	@echo "Application is running at http://localhost:8080"
	@echo "To view logs, run: make docker-logs"

status: ## Check service status
	@echo "=== Docker Services Status ==="
	@$(DOCKER_COMPOSE) ps
	@echo ""
	@echo "=== API Health Check ==="
	@curl -s http://localhost:8080/health || echo "API is not responding"

stop-all: ## Stop all services (both dev and prod)
	@$(DOCKER_COMPOSE) down
	@$(DOCKER_COMPOSE_DEV) down