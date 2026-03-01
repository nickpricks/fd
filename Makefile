.PHONY: help all tidy clean fmt vet test cover build install upgrade build-all

# The name of our output binary
BINARY_NAME=ft

# Default target
.DEFAULT_GOAL := help

# OS detection for cleaning
ifeq ($(OS),Windows_NT)
    CLEAN_CMD := del /q /f $(BINARY_NAME) 2>nul & del /q /f $(BINARY_NAME).exe 2>nul & rmdir /s /q bin 2>nul
else
    CLEAN_CMD := rm -rf $(BINARY_NAME) $(BINARY_NAME).exe bin/
endif

help:
	@echo "Available commands:"
	@echo "  make help      - Show this help message"
	@echo "  make all       - Run tidy, fmt, vet, test, and build"
	@echo "  make tidy      - Run go mod tidy"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make fmt       - Format code"
	@echo "  make vet       - Run go vet"
	@echo "  make test      - Run tests"
	@echo "  make cover     - Run tests with coverage"
	@echo "  make build     - Build the ft binary"
	@echo "  make build-all - Build for Windows, Linux, and macOS (in bin/)"
	@echo "  make install   - Install the binary to GOPATH/bin"
	@echo "  make upgrade   - Upgrade dependencies and reinstall"

all: tidy fmt vet test build

tidy:
	@echo "=> Running go mod tidy..."
	go mod tidy

clean:
	@echo "=> Cleaning build artifacts..."
	go clean
	-@$(CLEAN_CMD)

fmt:
	@echo "=> Formatting code..."
	go fmt ./...

vet:
	@echo "=> Running go vet..."
	go vet ./...

test:
	@echo "=> Running tests..."
	go test ./... -v

cover:
	@echo "=> Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "=> Coverage report generated at coverage.html"

build: clean
	@echo "=> Building binary..."
	go build -o $(BINARY_NAME) ./cmd/ft
	@echo "=> Build complete: $(BINARY_NAME)"

build-all: clean
	@echo "=> Building for multiple OS..."
ifeq ($(OS),Windows_NT)
	set GOOS=linux&& set GOARCH=amd64&& go build -o bin/ft-linux-amd64 ./cmd/ft
	set GOOS=darwin&& set GOARCH=amd64&& go build -o bin/ft-darwin-amd64 ./cmd/ft
	set GOOS=windows&& set GOARCH=amd64&& go build -o bin/ft-windows-amd64.exe ./cmd/ft
else
	GOOS=linux GOARCH=amd64 go build -o bin/ft-linux-amd64 ./cmd/ft
	GOOS=darwin GOARCH=amd64 go build -o bin/ft-darwin-amd64 ./cmd/ft
	GOOS=windows GOARCH=amd64 go build -o bin/ft-windows-amd64.exe ./cmd/ft
endif
	@echo "=> Cross-platform builds complete in bin/ directory"

install: build
	@echo "=> Installing ft to GOPATH/bin..."
	go install ./cmd/ft

upgrade:
	@echo "=> Upgrading dependencies and reinstalling..."
	go get -u ./...
	go mod tidy
	go install ./cmd/ft
