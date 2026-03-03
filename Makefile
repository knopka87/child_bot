# Child Bot Makefile
# ====================

# Load .env file if exists
-include .env
export

# Default values
POSTGRES_PASSWORD ?= root
DATABASE_URL ?= postgres://childbot:$(POSTGRES_PASSWORD)@localhost:5432/childbot?sslmode=disable
TEST_DATABASE_URL ?= postgres://childbot:$(POSTGRES_PASSWORD)@localhost:5432/childbot_test?sslmode=disable
MIGRATIONS_DIR ?= api/migrations
LLM_SERVER_URL ?= http://localhost:8081

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
db-up: ## Start PostgreSQL in Docker
	@echo "$(GREEN)Starting PostgreSQL...$(NC)"
	docker compose -f deploy/child-bot.compose.yml up -d db
	@echo "$(GREEN)Waiting for database to be ready...$(NC)"
	@sleep 3
	@docker compose -f deploy/child-bot.compose.yml exec db pg_isready -U childbot -d childbot || (echo "$(RED)Database not ready$(NC)" && exit 1)
	@echo "$(GREEN)Database is ready!$(NC)"

.PHONY: db-down
db-down: ## Stop PostgreSQL
	@echo "$(YELLOW)Stopping PostgreSQL...$(NC)"
	docker compose -f deploy/child-bot.compose.yml down db

.PHONY: db-logs
db-logs: ## Show database logs
	docker compose -f deploy/child-bot.compose.yml logs -f db

.PHONY: db-shell
db-shell: ## Open psql shell
	docker compose -f deploy/child-bot.compose.yml exec db psql -U childbot -d childbot

.PHONY: db-test-create
db-test-create: ## Create test database
	@echo "$(GREEN)Creating test database...$(NC)"
	docker compose -f deploy/child-bot.compose.yml exec db psql -U childbot -d postgres -c "DROP DATABASE IF EXISTS childbot_test;"
	docker compose -f deploy/child-bot.compose.yml exec db psql -U childbot -d postgres -c "CREATE DATABASE childbot_test;"
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
build: ## Build the application
	@echo "$(GREEN)Building...$(NC)"
	cd api && go build -o ../bin/server ./cmd/bot
	@echo "$(GREEN)Build complete: bin/server$(NC)"

.PHONY: run
run: ## Run the application locally
	@echo "$(GREEN)Running server...$(NC)"
	cd api && go run ./cmd/bot

# ====================
# Tests
# ====================

.PHONY: test
test: ## Run all short tests
	@echo "$(GREEN)Running short tests...$(NC)"
	cd api && go test -short -v ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	cd api && go test -short -coverprofile=coverage.out ./...
	cd api && go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: api/coverage.html$(NC)"

.PHONY: test-telegram
test-telegram: ## Run telegram package tests
	@echo "$(GREEN)Running telegram tests...$(NC)"
	cd api && go test -short -v ./internal/v2/telegram/...

.PHONY: test-e2e
test-e2e: ## Run E2E tests (requires LLM proxy and test DB)
	@echo "$(GREEN)Running E2E tests...$(NC)"
	@echo "$(YELLOW)Make sure LLM proxy is running at $(LLM_SERVER_URL)$(NC)"
	@echo "$(YELLOW)Make sure test database is ready$(NC)"
	@echo "$(YELLOW)Test images in: api/test/e2e/testdata/$(NC)"
	cd api && TEST_LLM_PROXY_URL=$(LLM_SERVER_URL) TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test -v ./test/e2e/... -timeout 30m

.PHONY: test-e2e-setup
test-e2e-setup: db-test-create migrate-test-up ## Setup test database for E2E tests
	@echo "$(GREEN)E2E test environment ready!$(NC)"

.PHONY: test-all
test-all: test test-e2e ## Run all tests (short + E2E)

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
docker-up: ## Start all services with Docker Compose
	@echo "$(GREEN)Starting all services...$(NC)"
	docker compose -f deploy/child-bot.compose.yml up -d

.PHONY: docker-down
docker-down: ## Stop all services
	@echo "$(YELLOW)Stopping all services...$(NC)"
	docker compose -f deploy/child-bot.compose.yml down

.PHONY: docker-logs
docker-logs: ## Show all logs
	docker compose -f deploy/child-bot.compose.yml logs -f

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
