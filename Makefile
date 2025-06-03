.PHONY: build build-linux build-x64 build-arm64 build-x86-64 build-x86 clean

# Build the binary
build:
	@echo "Building laravel-setup..."
	@go build -o laravel-setup ./cmd/laravel-setup

# Build the binary for Linux
build-linux:
	@echo "Building laravel-setup for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o laravel-setup-linux ./cmd/laravel-setup

# Build the binary for x64 architecture
build-x64:
	@echo "Building laravel-setup for x64 architecture..."
	@GOOS=darwin GOARCH=amd64 go build -o laravel-setup-x64 ./cmd/laravel-setup

# Build the binary for ARM64 architecture (Apple Silicon)
build-arm64:
	@echo "Building laravel-setup for ARM64 architecture..."
	@GOOS=darwin GOARCH=arm64 go build -o laravel-setup-arm64 ./cmd/laravel-setup

# Build the binary for x86 architecture (32-bit)
build-x86:
	@echo "Building laravel-setup for x86 architecture..."
	@GOOS=darwin GOARCH=386 go build -o laravel-setup-x86 ./cmd/laravel-setup

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f laravel-setup laravel-setup-linux laravel-setup-x64 laravel-setup-arm64 laravel-setup-x86

# Install the binary to /usr/local/bin
install: build
	@echo "Installing laravel-setup to /usr/local/bin..."
	@sudo cp laravel-setup /usr/local/bin/

# Run the binary
run: build
	@echo "Running laravel-setup..."
	@./laravel-setup


upload: build-linux
	@echo "Uploading to server"
	@scp ./laravel-setup-linux ./examples/config.example.toml skillshare:~/
	@rm -rf laravel-setup-linux

# Help message
help:
	@echo "Laravel Setup - A tool to set up Laravel production servers"
	@echo ""
	@echo "Usage:"
	@echo "  make build        - Build the binary"
	@echo "  make build-linux  - Build the binary for Linux (cross-compilation)"
	@echo "  make build-x64    - Build the binary for x64 architecture (macOS Intel)"
	@echo "  make build-arm64  - Build the binary for ARM64 architecture (Apple Silicon)"
	@echo "  make build-x86    - Build the binary for x86 architecture (32-bit)"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make install      - Install the binary to /usr/local/bin"
	@echo "  make run          - Run the binary"
	@echo "  make upload       - Upload the Linux binary and example config to server"
	@echo "  make help         - Show this help message"

# Default target
.DEFAULT_GOAL := build
