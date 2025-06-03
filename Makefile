.PHONY: build clean

# Build the binary
build:
	@echo "Building laravel-setup..."
	@go build -o laravel-setup ./cmd/laravel-setup

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f laravel-setup

# Install the binary to /usr/local/bin
install: build
	@echo "Installing laravel-setup to /usr/local/bin..."
	@sudo cp laravel-setup /usr/local/bin/

# Run the binary
run: build
	@echo "Running laravel-setup..."
	@./laravel-setup

# Help message
help:
	@echo "Laravel Setup - A tool to set up Laravel production servers"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - Build the binary"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make install  - Install the binary to /usr/local/bin"
	@echo "  make run      - Run the binary"
	@echo "  make help     - Show this help message"

# Default target
.DEFAULT_GOAL := build