# Makefile for routing-manager

# Docker image and container variables
IMAGE_NAME := routing-manager
IMAGE_TAG := latest
CONTAINER_NAME := routing-manager

# Version information
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date +'%Y-%m-%d_%H:%M:%S%Z')

# Configuration file path
CONFIG_FILE := ./config.yaml

# Optional overrides for testing/development
# These will only be used if explicitly passed as arguments to make
MQTT_BROKER_URL ?=
WORKER_BASE_URL ?=
LOG_FORMAT ?=
LOG_LEVEL ?=
NETWORK ?=

# Go build variables
GO := go
BINARY_NAME := routing-manager
BUILD_DIR := build
GO_BUILD_FLAGS := -v

# Default target
.PHONY: help
help:
	@echo "Routing Manager Makefile"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make build              Build the Docker image"
	@echo "  make run                Run the Docker container using config.yaml"
	@echo "  make stop               Stop and remove the Docker container"
	@echo "  make restart            Stop, remove, and run the Docker container"
	@echo "  make logs               View container logs"
	@echo "  make clean              Remove the Docker image and container"
	@echo "  make all                Build and run the Docker container"
	@echo ""
	@echo "Go Application Commands:"
	@echo "  make go-build           Build the Go application locally"
	@echo "  make go-run             Build and run the Go application locally"
	@echo "  make go-clean           Clean Go build artifacts"
	@echo "  make go-test            Run Go tests"
	@echo ""
	@echo "Configuration:"
	@echo "  make build IMAGE_TAG=dev                Build with custom tag"
	@echo "  make run CONFIG_FILE=/path/to/config.yaml  Run with specific config file"
	@echo "  make run MQTT_BROKER_URL=tcp://custom:1883  Override MQTT URL (optional)"
	@echo "  make run LOG_LEVEL=debug                Override log level (optional)"
	@echo "  make run NETWORK=custom-network         Run on a specific network"
	@echo ""
	@echo "Docker Compose:"
	@echo "  make up                 Start services with Docker Compose"
	@echo "  make down               Stop services with Docker Compose"
	@echo ""
	@echo "Development:"
	@echo "  make dev                Run a development container with mounted source"

# Build the Docker image
.PHONY: build
build:
	@echo "Building Docker image $(IMAGE_NAME):$(IMAGE_TAG)..."
	docker build \
	 --build-arg VERSION=$(VERSION) \
	 --build-arg COMMIT=$(COMMIT) \
	 --build-arg DATE=$(DATE) \
	 -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "Build complete."

# Run the Docker container
.PHONY: run
run: stop verify-config
	@echo "Running Docker container $(CONTAINER_NAME) with config from $(CONFIG_FILE)..."
	@docker_run_cmd="docker run -d --name $(CONTAINER_NAME)"; \
	if [ -n "$(NETWORK)" ]; then \
		docker_run_cmd="$$docker_run_cmd --network $(NETWORK)"; \
	fi; \
	if [ -n "$(MQTT_BROKER_URL)" ]; then \
		docker_run_cmd="$$docker_run_cmd -e MQTT_BROKER_URL=$(MQTT_BROKER_URL)"; \
	fi; \
	if [ -n "$(WORKER_BASE_URL)" ]; then \
		docker_run_cmd="$$docker_run_cmd -e WORKER_BASE_URL=$(WORKER_BASE_URL)"; \
	fi; \
	if [ -n "$(LOG_FORMAT)" ]; then \
		docker_run_cmd="$$docker_run_cmd -e LOG_FORMAT=$(LOG_FORMAT)"; \
	fi; \
	if [ -n "$(LOG_LEVEL)" ]; then \
		docker_run_cmd="$$docker_run_cmd -e LOG_LEVEL=$(LOG_LEVEL)"; \
	fi; \
	docker_run_cmd="$$docker_run_cmd -v $(CONFIG_FILE):/app/config.yaml"; \
	docker_run_cmd="$$docker_run_cmd $(IMAGE_NAME):$(IMAGE_TAG) --config /app/config.yaml"; \
	echo "Executing: $$docker_run_cmd"; \
	eval $$docker_run_cmd
	@echo "Container started."

# Stop and remove the Docker container
.PHONY: stop
stop:
	@echo "Stopping and removing container $(CONTAINER_NAME)..."
	@docker stop $(CONTAINER_NAME) 2>/dev/null || true
	@docker rm $(CONTAINER_NAME) 2>/dev/null || true
	@echo "Container stopped and removed."

# Restart the Docker container
.PHONY: restart
restart: stop run

# View container logs
.PHONY: logs
logs:
	@docker logs -f $(CONTAINER_NAME)

# Remove the Docker image and container
.PHONY: clean
clean: stop
	@echo "Removing Docker image $(IMAGE_NAME):$(IMAGE_TAG)..."
	@docker rmi $(IMAGE_NAME):$(IMAGE_TAG) 2>/dev/null || true
	@echo "Cleanup complete."

# Build and run the Docker container
.PHONY: all
all: build run

# Docker Compose commands
.PHONY: up
up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "Services started."

.PHONY: down
down:
	@echo "Stopping services with Docker Compose..."
	@docker-compose down
	@echo "Services stopped."

# Development targets
.PHONY: dev
dev: build
	@echo "Running development container..."
	@docker run -it --rm \
		-v $(PWD):/app \
		-w /app \
		$(IMAGE_NAME):$(IMAGE_TAG) \
		sh

# Verify config file exists
.PHONY: verify-config
verify-config:
	@if [ ! -f "$(CONFIG_FILE)" ]; then \
		echo "Error: Configuration file $(CONFIG_FILE) not found!"; \
		exit 1; \
	fi
	@echo "Configuration file $(CONFIG_FILE) exists."

# Go application targets
.PHONY: go-build
go-build:
	@echo "Building Go application..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: go-run
go-run: verify-config go-build
	@echo "Running Go application with config from $(CONFIG_FILE)..."
	$(BUILD_DIR)/$(BINARY_NAME) --config $(CONFIG_FILE)

.PHONY: go-clean
go-clean:
	@echo "Cleaning Go build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

.PHONY: go-test
go-test:
	@echo "Running tests..."
	$(GO) test -v ./...
	@echo "Tests complete." 