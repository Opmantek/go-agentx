# Makefile for go-agentx Docker operations

.PHONY: help build test test-verbose shell clean

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the Docker image"
	@echo "  test         - Run all tests in Docker"
	@echo "  test-verbose - Run tests with verbose output"
	@echo "  shell        - Start a shell in the container"
	@echo "  clean        - Remove the Docker image"

# Build the Docker image
build:
	docker build -t go-agentx-test .

# Run tests
test: build
	docker run --rm \
		--cap-add NET_ADMIN \
		--cap-add SYS_ADMIN \
		-v "$$(pwd)":/app \
		go-agentx-test

# Run tests with verbose output
test-verbose: build
	docker run --rm \
		--cap-add NET_ADMIN \
		--cap-add SYS_ADMIN \
		-v "$$(pwd)":/app \
		go-agentx-test \
		go test -v -race ./...

# Start an interactive shell in the container
shell: build
	docker run --rm -it \
		--cap-add NET_ADMIN \
		--cap-add SYS_ADMIN \
		-v "$$(pwd)":/app \
		go-agentx-test \
		/bin/sh

# Clean up
clean:
	docker rmi go-agentx-test || true