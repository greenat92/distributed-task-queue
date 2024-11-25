# Variables
APP_NAME=distributed-task-queue
BUILD_DIR=bin
SRC_WORKER_DIR=cmd/worker
SRC_API_DIR=cmd/api
WORKER_MAIN=$(SRC_WORKER_DIR)/main.go
API_MAIN=$(SRC_API_DIR)/main.go

# Default goal
.DEFAULT_GOAL := run-worker

# Run the worker
run-worker:
	@echo "Running the worker..."
	@go run $(WORKER_MAIN)

# Run the API server
run-api:
	@echo "Running the API server..."
	@go run $(API_MAIN)

# Build the worker
build-worker:
	@echo "Building the worker..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/worker $(WORKER_MAIN)

# Build the API server
build-api:
	@echo "Building the API server..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/api $(API_MAIN)

# Build both worker and API
build-all: build-worker build-api
	@echo "Build completed for both worker and API."

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete!"

# Format the code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Check for linting issues
lint:
	@echo "Running golint..."
	@golint ./...
