# Child Bot Makefile
# ====================

# Load .env file if exists
-include .env
export

# Default values
POSTGRES_PASSWORD ?= dev_secret
DATABASE_URL ?= postgres://child_bot:$(POSTGRES_PASSWORD)@localhost:5432/child_bot?sslmode=disable
TEST_DATABASE_URL ?= postgres://child_bot:$(POSTGRES_PASSWORD)@localhost:5432/child_bot_test?sslmode=disable
MIGRATIONS_DIR ?= api/migrations
LLM_SERVER_URL ?= http://138.124.55.145:8000

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
NC     := \033[0m # No Color

.PHONY: help
help: ## Show this help
	@echo "Child Bot - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*## .*$$' $(MAKEFILE_LIST) | sed 's/^Makefile://' | sort | awk 'BEGIN {FS = ":.*## "}; {printf "  \033[0;32m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ====================
# Docker / Database
# ====================

.PHONY: db-up
db-up: ## Start PostgreSQL and Redis (development)
	@echo "$(GREEN)Starting PostgreSQL and Redis...$(NC)"
	docker compose -f docker/docker-compose.dev.yml up -d postgres redis
	@echo "$(GREEN)Waiting for database to be ready...$(NC)"
	@sleep 3
	@docker compose -f docker/docker-compose.dev.yml exec postgres pg_isready -U child_bot -d child_bot || (echo "$(RED)Database not ready$(NC)" && exit 1)
	@echo "$(GREEN)Database is ready!$(NC)"

.PHONY: db-down
db-down: ## Stop PostgreSQL and Redis
	@echo "$(YELLOW)Stopping database services...$(NC)"
	docker compose -f docker/docker-compose.dev.yml down postgres redis

.PHONY: db-logs
db-logs: ## Show database logs
	docker compose -f docker/docker-compose.dev.yml logs -f postgres

.PHONY: db-shell
db-shell: ## Open psql shell
	docker compose -f docker/docker-compose.dev.yml exec postgres psql -U child_bot -d child_bot

.PHONY: db-test-create
db-test-create: ## Create test database
	@echo "$(GREEN)Creating test database...$(NC)"
	docker compose -f docker/docker-compose.dev.yml exec postgres psql -U child_bot -d postgres -c "DROP DATABASE IF EXISTS child_bot_test;"
	docker compose -f docker/docker-compose.dev.yml exec postgres psql -U child_bot -d postgres -c "CREATE DATABASE child_bot_test;"
	@echo "$(GREEN)Test database created!$(NC)"

# ====================
# Migrations
# ====================

.PHONY: migrate-install
migrate-install: ## Install golang-migrate CLI
	@echo "$(GREEN)Installing golang-migrate...$(NC)"
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

.PHONY: migrate-up
migrate-up: ## Run all pending migrations
	@echo "$(GREEN)Running migrations...$(NC)"
	migrate -source "file://$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" up
	@echo "$(GREEN)Migrations completed!$(NC)"

.PHONY: migrate-down
migrate-down: ## Rollback last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	migrate -source "file://$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" down 1

.PHONY: migrate-down-all
migrate-down-all: ## Rollback all migrations (DANGER!)
	@echo "$(RED)Rolling back ALL migrations...$(NC)"
	@read -p "Are you sure? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	migrate -source "file://$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" down -all

.PHONY: migrate-status
migrate-status: ## Show migration status
	@echo "$(GREEN)Migration status:$(NC)"
	migrate -source "file://$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" version

.PHONY: migrate-force
migrate-force: ## Force set migration version (use: make migrate-force VERSION=N)
	@echo "$(YELLOW)Forcing migration version to $(VERSION)...$(NC)"
	migrate -source "file://$(MIGRATIONS_DIR)" -database "$(DATABASE_URL)" force $(VERSION)

.PHONY: migrate-create
migrate-create: ## Create new migration (use: make migrate-create NAME=description)
	@echo "$(GREEN)Creating migration: $(NAME)$(NC)"
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

# Test database migrations
.PHONY: migrate-test-up
migrate-test-up: ## Run migrations on test database
	@echo "$(GREEN)Running migrations on test database...$(NC)"
	migrate -source "file://$(MIGRATIONS_DIR)" -database "$(TEST_DATABASE_URL)" up
	@echo "$(GREEN)Test migrations completed!$(NC)"

.PHONY: migrate-test-down-all
migrate-test-down-all: ## Rollback all migrations on test database
	@echo "$(YELLOW)Rolling back all test migrations...$(NC)"
	migrate -source "file://$(MIGRATIONS_DIR)" -database "$(TEST_DATABASE_URL)" down -all

# ====================
# Build
# ====================

.PHONY: build
build: ## Build REST API server
	@echo "$(GREEN)Building REST API server...$(NC)"
	cd api && go build -o ../bin/server ./cmd/server
	@echo "$(GREEN)Build complete: bin/server$(NC)"

.PHONY: run
run: ## Run REST API server locally
	@echo "$(GREEN)Running REST API server...$(NC)"
	cd api && go run ./cmd/server

# ====================
# Tests
# ====================

.PHONY: test
test: ## Run all unit tests (short mode)
	@echo "$(GREEN)Running unit tests...$(NC)"
	cd api && go test -short -v ./...

.PHONY: test-unit
test-unit: ## Run unit tests for specific package (use: make test-unit PKG=handler)
	@echo "$(GREEN)Running unit tests for $(PKG)...$(NC)"
	cd api && go test -short -v ./internal/api/$(PKG)/...

.PHONY: test-handlers
test-handlers: ## Run handler unit tests
	@echo "$(GREEN)Running handler tests...$(NC)"
	cd api && go test -short -v ./internal/api/handler/...

.PHONY: test-middleware
test-middleware: ## Run middleware unit tests
	@echo "$(GREEN)Running middleware tests...$(NC)"
	cd api && go test -short -v ./internal/api/middleware/...

.PHONY: test-service
test-service: ## Run service integration tests (requires TEST_DATABASE_URL)
	@echo "$(GREEN)Running service integration tests...$(NC)"
	@echo "$(YELLOW)Make sure test database is ready: make db-test-create migrate-test-up$(NC)"
	cd api && TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test -v ./internal/service/... -timeout 5m

.PHONY: test-cover
test-cover: ## Run tests with coverage
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	cd api && go test -short -coverprofile=coverage.out ./...
	cd api && go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: api/coverage.html$(NC)"

.PHONY: test-race
test-race: ## Run tests with race detector
	@echo "$(GREEN)Running tests with race detector...$(NC)"
	cd api && go test -short -race -v ./...


.PHONY: test-e2e-rest
test-e2e-rest: ## Run REST API E2E tests (requires test DB, fast with mock LLM)
	@echo "$(GREEN)Running REST API E2E tests...$(NC)"
	@echo "$(YELLOW)Using mock LLM for fast tests$(NC)"
	cd api && TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test -v ./test/e2e/rest_api_test.go -timeout 10m

.PHONY: test-e2e-rest-real
test-e2e-rest-real: ## Run REST API E2E tests with real LLM (requires LLM proxy)
	@echo "$(GREEN)Running REST API E2E tests with real LLM...$(NC)"
	@echo "$(YELLOW)Make sure LLM proxy is running at $(LLM_SERVER_URL)$(NC)"
	cd api && USE_REAL_LLM=true LLM_PROXY_URL=$(LLM_SERVER_URL) TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test -v ./test/e2e/rest_api_test.go -timeout 30m


.PHONY: test-e2e-setup
test-e2e-setup: db-test-create migrate-test-up ## Setup test database for E2E tests
	@echo "$(GREEN)E2E test environment ready!$(NC)"
	@echo "Run: make test-e2e-rest (fast) or make test-e2e-rest-real (with LLM)"

.PHONY: test-e2e
test-e2e: test-e2e-rest ## Run all E2E tests (REST API only, fast)

.PHONY: test-e2e-all
test-e2e-all: test-e2e-rest ## Run all E2E tests (REST API only)

.PHONY: test-integration
test-integration: test-service ## Run all integration tests

.PHONY: test-all
test-all: test test-integration test-e2e ## Run all tests (unit + integration + E2E)

.PHONY: test-ci
test-ci: test-race test-cover ## Run tests for CI (with race detector and coverage)

# ====================
# Development
# ====================

.PHONY: lint
lint: ## Run linter
	@echo "$(GREEN)Running linter...$(NC)"
	cd api && golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	cd api && go fmt ./...

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "$(GREEN)Tidying modules...$(NC)"
	go mod tidy

.PHONY: vendor
vendor: ## Vendor dependencies
	@echo "$(GREEN)Vendoring dependencies...$(NC)"
	go mod vendor

# ====================
# Docker
# ====================

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t child-bot:local .

.PHONY: docker-up
docker-up: ## Start all services with Docker Compose (alias for dev)
	@$(MAKE) dev

.PHONY: docker-down
docker-down: ## Stop all services (alias for dev-down)
	@$(MAKE) dev-down

.PHONY: docker-logs
docker-logs: ## Show all logs (alias for dev-logs)
	@$(MAKE) dev-logs

# ====================
# Cleanup
# ====================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -rf bin/
	rm -rf api/coverage.out api/coverage.html
	rm -rf api/test/e2e/results/*.json

.DEFAULT_GOAL := help

# ====================
# Development (Full Stack)
# ====================

.PHONY: dev
dev: ## Start full stack in development mode
	@echo "$(GREEN)Starting development environment...$(NC)"
	docker compose -f docker/docker-compose.dev.yml up -d
	@echo "$(GREEN)Development environment started!$(NC)"
	@echo "Frontend: http://localhost:5173"
	@echo "Backend:  http://localhost:8080"
	@echo "Postgres: localhost:5432"
	@echo "Redis:    localhost:6379"

.PHONY: dev-down
dev-down: ## Stop development environment
	@echo "$(YELLOW)Stopping development environment...$(NC)"
	docker compose -f docker/docker-compose.dev.yml down

.PHONY: dev-logs
dev-logs: ## Show development logs
	docker compose -f docker/docker-compose.dev.yml logs -f

.PHONY: dev-backend
dev-backend: ## Run backend only (requires DB)
	@echo "$(GREEN)Starting backend in development mode...$(NC)"
	docker compose -f docker/docker-compose.dev.yml up -d postgres redis
	@sleep 2
	cd api && go run ./cmd/server

.PHONY: dev-frontend
dev-frontend: ## Run frontend only (requires backend)
	@echo "$(GREEN)Starting frontend in development mode...$(NC)"
	cd frontend && npm run dev

# ====================
# Production (Full Stack)
# ====================

.PHONY: prod-up
prod-up: ## Start production environment
	@echo "$(GREEN)Starting production environment...$(NC)"
	docker compose up -d
	@echo "$(GREEN)Production environment started!$(NC)"

.PHONY: prod-down
prod-down: ## Stop production environment
	@echo "$(YELLOW)Stopping production environment...$(NC)"
	docker compose down

.PHONY: prod-logs
prod-logs: ## Show production logs
	docker compose logs -f

.PHONY: prod-restart
prod-restart: ## Restart production services
	@echo "$(YELLOW)Restarting production services...$(NC)"
	docker compose restart

# ====================
# Utilities
# ====================

.PHONY: docker-clean
docker-clean: ## Remove all containers and volumes
	@echo "$(RED)Cleaning up Docker resources...$(NC)"
	@read -p "This will remove all data. Are you sure? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	docker compose -f docker/docker-compose.dev.yml down -v
	docker compose down -v
	@echo "$(GREEN)Cleanup complete!$(NC)"

.PHONY: frontend-shell
frontend-shell: ## Open shell in frontend container
	docker compose exec frontend sh

.PHONY: backend-shell
backend-shell: ## Open shell in backend container
	docker compose exec backend sh

.PHONY: health
health: ## Check health of all services
	@echo "$(GREEN)Checking services health...$(NC)"
	@echo -n "Backend:  " && curl -s http://localhost:8080/health | grep -q ok && echo "OK" || echo "DOWN"
	@docker compose exec postgres pg_isready -U child_bot 2>/dev/null && echo "Postgres: OK" || echo "Postgres: DOWN"
	@docker compose exec redis redis-cli --pass dev_redis_secret ping 2>/dev/null | grep -q PONG && echo "Redis: OK" || echo "Redis: DOWN"
