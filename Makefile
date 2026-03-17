BINARY := shopify-admin-mcp-server
BUILD_DIR := build
VERSION := $(shell cat .version 2>/dev/null | tr -d '[:space:]')

LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

# Detect host OS
OS := $(shell uname -s 2>/dev/null | tr '[:upper:]' '[:lower:]')
ifeq ($(findstring mingw,$(OS)),mingw)
  HOST_OS := windows
else ifeq ($(findstring msys,$(OS)),msys)
  HOST_OS := windows
else
  HOST_OS := $(OS)
endif

all: build

# dev: run backend in development mode
dev:
	@echo "Starting in development mode..."
	@DEV_MODE=true go run -C backend $(LDFLAGS) ./cmd/$(BINARY)

# build: compile binary to build/
build: clean
	@echo "Building $(BINARY)..."
	CGO_ENABLED=0 go build -C backend $(LDFLAGS) -o ../$(BUILD_DIR)/$(BINARY) ./cmd/$(BINARY)
	@echo "Build complete. Output in $(BUILD_DIR)/"

# build-win: build Windows exe
# - On Linux/macOS: cross-compiles with GOOS=windows
# - On Windows (Git Bash/MSYS2): builds natively without GOOS override
build-win:
	@echo "Building $(BINARY).exe for Windows (host: $(HOST_OS))..."
ifeq ($(HOST_OS),windows)
	CGO_ENABLED=0 go build -C backend $(LDFLAGS) -o ../$(BUILD_DIR)/$(BINARY).exe ./cmd/$(BINARY)
else
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -C backend $(LDFLAGS) -o ../$(BUILD_DIR)/$(BINARY).exe ./cmd/$(BINARY)
endif
	@echo "Build complete. Output in $(BUILD_DIR)/$(BINARY).exe"

# docker-build-images: build Docker image with auto-versioned tags (latest, vX.Y, vX.Y.Z)
# Override version with: make docker-build-images VERSION=v1.2.3
docker-build-images:
	@VERSION=$(VERSION) utils/docker-build-images.sh

# docker-push-images: push all tags to Docker Hub (run docker-build-images first)
docker-push-images:
	@utils/docker-push.sh

# clean: remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# test: run all tests
test:
	go test -C backend -v ./...

# test-coverage: run tests with coverage report
test-coverage:
	go test -C backend -v -coverprofile=../coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# lint: run linter
lint:
	golangci-lint run ./backend/...

# fmt: format code
fmt:
	go fmt -C backend ./...

# tidy: tidy Go modules
tidy:
	go mod -C backend tidy

# help: show this help message
help:
	@echo "Available targets:"
	@echo "  dev            - Run backend server in development mode"
	@echo "  build          - Build binary to $(BUILD_DIR)/ (version from .version)"
	@echo "  build-win      - Cross-compile Windows exe to $(BUILD_DIR)/"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy Go modules"
	@echo "  docker-build-images - Build Docker image (auto-increment patch, or VERSION=vX.Y.Z)"
	@echo "  docker-push-images  - Push all tags to Docker Hub"
	@echo "  help           - Show this help message"

.PHONY: all dev build build-win clean test test-coverage lint fmt tidy docker-build-images docker-push-images help
