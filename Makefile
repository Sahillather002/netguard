.PHONY: help install dev build test lint clean docker-build docker-up docker-down deploy

# Default target
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

##@ General

help: ## Display this help message
	@echo "$(BLUE)NetGuard - Enterprise Security Platform$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make $(GREEN)<target>$(NC)\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

install: ## Install all dependencies
	@echo "$(BLUE)Installing dependencies...$(NC)"
	npm install
	cd apps/web && npm install
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

dev: ## Start development servers
	@echo "$(BLUE)Starting development servers...$(NC)"
	npm run dev

dev-web: ## Start frontend only
	@echo "$(BLUE)Starting web app...$(NC)"
	cd apps/web && npm run dev

dev-services: ## Start backend services only
	@echo "$(BLUE)Starting backend services...$(NC)"
	docker-compose up api-gateway auth-service threat-detector

setup-db: ## Set up databases
	@echo "$(BLUE)Setting up databases...$(NC)"
	docker-compose up -d postgres redis influxdb elasticsearch
	@echo "$(GREEN)✓ Databases started$(NC)"

migrate: ## Run database migrations
	@echo "$(BLUE)Running migrations...$(NC)"
	npm run migrate
	@echo "$(GREEN)✓ Migrations complete$(NC)"

seed: ## Seed database with test data
	@echo "$(BLUE)Seeding database...$(NC)"
	npm run seed
	@echo "$(GREEN)✓ Database seeded$(NC)"

##@ Building

build: ## Build all applications
	@echo "$(BLUE)Building applications...$(NC)"
	npm run build
	@echo "$(GREEN)✓ Build complete$(NC)"

build-web: ## Build web app
	@echo "$(BLUE)Building web app...$(NC)"
	cd apps/web && npm run build
	@echo "$(GREEN)✓ Web app built$(NC)"

build-services: ## Build all services
	@echo "$(BLUE)Building services...$(NC)"
	cd services/api-gateway && go build
	cd services/auth-service && go build
	cd services/threat-detector && pip install -r requirements.txt
	cd services/network-monitor && cargo build --release
	cd services/firewall-service && cargo build --release
	@echo "$(GREEN)✓ Services built$(NC)"

##@ Testing

test: ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	npm run test
	@echo "$(GREEN)✓ Tests complete$(NC)"

test-web: ## Run frontend tests
	@echo "$(BLUE)Running web tests...$(NC)"
	cd apps/web && npm run test
	@echo "$(GREEN)✓ Web tests complete$(NC)"

test-services: ## Run backend tests
	@echo "$(BLUE)Running service tests...$(NC)"
	cd services/api-gateway && go test ./...
	cd services/auth-service && go test ./...
	cd services/threat-detector && pytest
	@echo "$(GREEN)✓ Service tests complete$(NC)"

test-e2e: ## Run end-to-end tests
	@echo "$(BLUE)Running E2E tests...$(NC)"
	npm run test:e2e
	@echo "$(GREEN)✓ E2E tests complete$(NC)"

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	npm run test:coverage
	@echo "$(GREEN)✓ Coverage report generated$(NC)"

##@ Code Quality

lint: ## Run linters
	@echo "$(BLUE)Running linters...$(NC)"
	npm run lint
	@echo "$(GREEN)✓ Linting complete$(NC)"

lint-fix: ## Fix linting issues
	@echo "$(BLUE)Fixing linting issues...$(NC)"
	npm run lint:fix
	@echo "$(GREEN)✓ Linting fixed$(NC)"

format: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	npm run format
	@echo "$(GREEN)✓ Code formatted$(NC)"

typecheck: ## Run type checking
	@echo "$(BLUE)Running type check...$(NC)"
	npm run typecheck
	@echo "$(GREEN)✓ Type check complete$(NC)"

##@ Docker

docker-build: ## Build Docker images
	@echo "$(BLUE)Building Docker images...$(NC)"
	docker-compose build
	@echo "$(GREEN)✓ Docker images built$(NC)"

docker-up: ## Start Docker containers
	@echo "$(BLUE)Starting Docker containers...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)✓ Containers started$(NC)"

docker-down: ## Stop Docker containers
	@echo "$(BLUE)Stopping Docker containers...$(NC)"
	docker-compose down
	@echo "$(GREEN)✓ Containers stopped$(NC)"

docker-logs: ## View Docker logs
	docker-compose logs -f

docker-ps: ## List running containers
	docker-compose ps

docker-clean: ## Remove all containers and volumes
	@echo "$(YELLOW)Removing all containers and volumes...$(NC)"
	docker-compose down -v
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

##@ Deployment

deploy-staging: ## Deploy to staging
	@echo "$(BLUE)Deploying to staging...$(NC)"
	./scripts/deploy/staging.sh
	@echo "$(GREEN)✓ Deployed to staging$(NC)"

deploy-prod: ## Deploy to production
	@echo "$(YELLOW)Deploying to production...$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		./scripts/deploy/production.sh; \
		echo "$(GREEN)✓ Deployed to production$(NC)"; \
	else \
		echo "$(RED)✗ Deployment cancelled$(NC)"; \
	fi

k8s-deploy: ## Deploy to Kubernetes
	@echo "$(BLUE)Deploying to Kubernetes...$(NC)"
	kubectl apply -f infrastructure/kubernetes/
	@echo "$(GREEN)✓ Deployed to Kubernetes$(NC)"

helm-install: ## Install with Helm
	@echo "$(BLUE)Installing with Helm...$(NC)"
	helm install netguard ./infrastructure/helm/netguard
	@echo "$(GREEN)✓ Helm installation complete$(NC)"

helm-upgrade: ## Upgrade Helm release
	@echo "$(BLUE)Upgrading Helm release...$(NC)"
	helm upgrade netguard ./infrastructure/helm/netguard
	@echo "$(GREEN)✓ Helm upgrade complete$(NC)"

##@ Monitoring

logs: ## View application logs
	docker-compose logs -f

logs-web: ## View web app logs
	docker-compose logs -f web

logs-api: ## View API logs
	docker-compose logs -f api-gateway

metrics: ## Open Prometheus
	@echo "$(BLUE)Opening Prometheus...$(NC)"
	open http://localhost:9090

dashboard: ## Open Grafana
	@echo "$(BLUE)Opening Grafana...$(NC)"
	open http://localhost:3001

##@ Database

db-console: ## Open database console
	docker-compose exec postgres psql -U netguard

db-backup: ## Backup database
	@echo "$(BLUE)Backing up database...$(NC)"
	./scripts/maintenance/backup-db.sh
	@echo "$(GREEN)✓ Database backed up$(NC)"

db-restore: ## Restore database
	@echo "$(BLUE)Restoring database...$(NC)"
	./scripts/maintenance/restore-db.sh
	@echo "$(GREEN)✓ Database restored$(NC)"

db-reset: ## Reset database (WARNING: deletes all data)
	@echo "$(RED)WARNING: This will delete all data!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		docker-compose up -d postgres redis influxdb; \
		sleep 5; \
		make migrate; \
		echo "$(GREEN)✓ Database reset$(NC)"; \
	else \
		echo "$(YELLOW)✗ Reset cancelled$(NC)"; \
	fi

##@ Utilities

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	rm -rf node_modules
	rm -rf apps/web/.next
	rm -rf apps/web/node_modules
	rm -rf services/*/target
	rm -rf services/*/node_modules
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

clean-all: clean docker-clean ## Clean everything including Docker

update-deps: ## Update dependencies
	@echo "$(BLUE)Updating dependencies...$(NC)"
	npm update
	cd apps/web && npm update
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

security-audit: ## Run security audit
	@echo "$(BLUE)Running security audit...$(NC)"
	npm audit
	cd apps/web && npm audit
	@echo "$(GREEN)✓ Security audit complete$(NC)"

check-health: ## Check service health
	@echo "$(BLUE)Checking service health...$(NC)"
	curl -f http://localhost:3000/health || echo "$(RED)✗ Web app not healthy$(NC)"
	curl -f http://localhost:8080/health || echo "$(RED)✗ API Gateway not healthy$(NC)"
	@echo "$(GREEN)✓ Health check complete$(NC)"

##@ Documentation

docs-serve: ## Serve documentation locally
	@echo "$(BLUE)Serving documentation...$(NC)"
	cd docs && python -m http.server 8000

docs-build: ## Build documentation
	@echo "$(BLUE)Building documentation...$(NC)"
	npm run docs:build
	@echo "$(GREEN)✓ Documentation built$(NC)"

api-docs: ## Generate API documentation
	@echo "$(BLUE)Generating API docs...$(NC)"
	npm run api-docs
	@echo "$(GREEN)✓ API docs generated$(NC)"