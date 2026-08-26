# Variables
BINARY_NAME=lem-in
GO_FILES=$(shell find . -name '*.go')

.PHONY: all build test clean fmt vet help

# Default target: check code format, run tests, and build binary
all: fmt vet test build

# Compile the project binary
build:
	@echo "Building binary: $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) main.go

# Run all unit tests within subdirectories
test:
	@echo "Running all tests..."
	@go test -v ./...

# Run code coverage analysis
cover:
	@echo "Running test coverage analysis..."
	@go test -cover ./...

# Clean up build artifacts and binary outputs
clean:
	@echo "Cleaning up build artifacts..."
	@rm -f $(BINARY_NAME)
	@go clean

# Auto-format all Go source files according to language specifications
fmt:
	@echo "Formatting Go files..."
	@go fmt ./...

# Run Go static analysis tool to inspect common issues
vet:
	@echo "Vetting source files for software bugs..."
	@go vet ./...

# Display help information regarding all possible execution options
help:
	@echo "Available commands in this Makefile:"
	@echo "  make all      - Format, vet, test, and build the program binary"
	@echo "  make build    - Compile main entry file into execution binary"
	@echo "  make test     - Run complete verification matrix suites"
	@echo "  make cover    - Inspect code test coverage percentage values"
	@echo "  make clean    - Remove local binary artifacts and dependency cache"
	@echo "  make fmt      - Enforce uniform code styling guidelines"
	@echo "  make vet      - Run basic compilation bug inspections"
