.PHONY: help build run test clean docker-build docker-run docker-stop

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the backend application"
	@echo "  run           - Run the backend application"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker images"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker Compose services"
	@echo "  dev           - Run in development mode"

# Build backend
build:
	cd backend && go build -o main cmd/main.go

# Run backend
run:
	cd backend && go run cmd/main.go

# Run tests
test:
	cd backend && go test ./...

# Clean build artifacts
clean:
	rm -f backend/main
	rm -f backend/main.exe
	rm -f backend/todo.db

# Build Docker images
docker-build:
	docker-compose build

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose services
docker-stop:
	docker-compose down

# Development mode
dev:
	@echo "Starting development environment..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Production deployment
deploy:
	@echo "Deploying to production..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
