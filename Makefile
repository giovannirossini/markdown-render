# Markdown Render

# Variables
BINARY_NAME=markdown-render
BUILD_DIR=bin
MAIN_PATH=main.go
GO_FILES=$(shell find . -name '*.go')
OUTPUT=$(BUILD_DIR)/$(BINARY_NAME)

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD)fmt

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test run tidy vet fmt-check help

all: clean tidy test build

## build: Build the binary (use OUTPUT=path/to/binary to override output path)
build:
	@echo "Building $(BINARY_NAME)..."
	@OUTPUT_DIR="$$(dirname $(OUTPUT))"; \
	if [ "$$OUTPUT_DIR" != "." ]; then \
		mkdir -p "$$OUTPUT_DIR"; \
	fi
	$(GOBUILD) $(LDFLAGS) -o $(OUTPUT) $(MAIN_PATH)

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

## test: Run unit tests with coverage
test:
	@echo "Running tests with coverage..."
	@$(GOTEST) -v ./render
	@echo ""
	@echo "Coverage summary:"
	@$(GOTEST) -coverprofile=coverage.out ./render > /dev/null 2>&1 && \
		go tool cover -func=coverage.out | tail -1 || true
	@rm -f coverage.out

## tidy: Clean up go.mod and go.sum
tidy:
	@echo "Tidying up modules..."
	$(GOMOD) tidy

## vet: Run go vet on all packages
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

## fmt-check: Check if code is properly formatted
fmt-check:
	@echo "Checking code formatting..."
	@if [ "$$($(GOFMT) -s -l . | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted. Run 'go fmt ./...'"; \
		$(GOFMT) -s -d .; \
		exit 1; \
	fi
	@echo "Code is properly formatted."

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' Makefile | sed -e 's/## //g' | column -t -s ':'
