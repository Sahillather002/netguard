.PHONY: build run test clean docker-build docker-run help

APP_NAME=netguard-api
VERSION=1.0.0
BUILD_DIR=build

help:
	@echo "NetGuard API - Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make lint         - Run linter"
	@echo "  make deps         - Download dependencies"

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME).exe .
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME).exe"

build-linux:
	@echo "Building $(APP_NAME) for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run:
	@echo "Running $(APP_NAME)..."
	@go run .

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
	@if exist $(APP_NAME).exe del $(APP_NAME).exe
	@echo "Clean complete"

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):$(VERSION) .
	@docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo "Docker image built: $(APP_NAME):$(VERSION)"

docker-run:
	@echo "Running Docker container..."
	@docker run -d -p 8080:8080 --name $(APP_NAME) $(APP_NAME):latest
	@echo "Container started: $(APP_NAME)"

docker-stop:
	@docker stop $(APP_NAME)
	@docker rm $(APP_NAME)

lint:
	@echo "Running linter..."
	@golangci-lint run

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

install:
	@echo "Installing $(APP_NAME)..."
	@go install .
	@echo "Installation complete"
