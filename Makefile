# Makefile for getboy

.PHONY: test build run clean help

# Default target
all: test build

# Run tests (with gotestsum if available, otherwise plain go test)
test:
	@if command -v gotestsum > /dev/null 2>&1; then \
		gotestsum --format testname ./...; \
	elif [ -f ~/go/bin/gotestsum ]; then \
		~/go/bin/gotestsum --format testname ./...; \
	else \
		echo "Running tests... (install gotestsum for colorized output: go install gotest.tools/gotestsum@latest)"; \
		go test -v ./...; \
	fi

# Run tests with coverage
test-coverage:
	@if command -v gotestsum > /dev/null 2>&1; then \
		gotestsum --format testname -- -cover ./...; \
	elif [ -f ~/go/bin/gotestsum ]; then \
		~/go/bin/gotestsum --format testname -- -cover ./...; \
	else \
		echo "Running tests with coverage..."; \
		go test -cover ./...; \
	fi

# Build the application
build:
	@echo "Building getboy..."
	@go build -o getboy ./cmd/getboy

# Run the application (without building)
run:
	@go run ./cmd/getboy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f getboy
	@go clean

# Show help
help:
	@echo "Available targets:"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application (no build)"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make help           - Show this help message"
