# GoFiver Makefile - Easy commands like Laravel's artisan

.PHONY: help up down build logs api-logs frontend-logs shell-api shell-frontend

# Default target
help:
	@echo "GoFiver - Go + Fiber + React"
	@echo ""
	@echo "Usage:"
	@echo "  make up              Start all services"
	@echo "  make down            Stop all services"
	@echo "  make build           Build all containers"
	@echo "  make logs            View all logs"
	@echo "  make api-logs        View API logs"
	@echo "  make frontend-logs   View frontend logs"
	@echo "  make shell-api       Open shell in API container"
	@echo "  make shell-frontend  Open shell in frontend container"
	@echo "  make restart         Restart all services"
	@echo ""

# Start services
up:
	docker compose up -d

# Stop services
down:
	docker compose down

# Build containers
build:
	docker compose build

# View all logs
logs:
	docker compose logs -f

# View API logs
api-logs:
	docker logs -f gofiver-api

# View frontend logs
frontend-logs:
	docker logs -f gofiver-frontend

# Shell into API container
shell-api:
	docker exec -it gofiver-api sh

# Shell into frontend container
shell-frontend:
	docker exec -it gofiver-frontend sh

# Restart services
restart:
	docker compose restart

# Fresh start (rebuild everything)
fresh:
	docker compose down -v
	docker compose build --no-cache
	docker compose up -d
