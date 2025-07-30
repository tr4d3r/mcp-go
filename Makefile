# MCP Go Server Makefile

.PHONY: help init deps build run dev test lint format clean docker-build docker-run tools health

# Default target
help:
	@echo "Available targets:"
	@echo "  init        - Initialize Go module (run once)"
	@echo "  deps        - Install dependencies"
	@echo "  build       - Build the binary"
	@echo "  run         - Run the server"
	@echo "  dev         - Run with hot reload using Air"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linter"
	@echo "  format      - Format code"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run Docker container"

# Build the binary
build:
	mkdir -p bin
	go build -o bin/mcp-server main.go

# Run the server
run:
	go run main.go

# Run with hot reload
dev:
	air

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Format code
format:
	goimports -w .
	go fmt ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf tmp/
	go clean

# Install dependencies
deps:
	go mod tidy
	go mod download

# Initialize a new Go module (only run once)
init:
	go mod init mcp-go-server
	go mod tidy

# Install development tools
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/air-verse/air@latest

# Build Docker image
docker-build:
	docker build -t mcp-go-server .

# Run Docker container
docker-run:
	docker run -p 8080:8080 mcp-go-server

# Check health endpoint
health:
	curl -s http://localhost:8080/health | jq .
