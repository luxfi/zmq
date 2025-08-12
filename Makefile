# zmq4 Makefile

.PHONY: all build test test-verbose test-race test-cover test-czmq test-czmq-verbose test-czmq-race bench lint fmt clean help

# Default target
all: build

# Build the project
build:
	GO111MODULE=on go build -mod=readonly .

# Run tests
test:
	GO111MODULE=on go test -mod=readonly ./...

# Run tests with verbose output
test-verbose:
	GO111MODULE=on go test -v -mod=readonly ./...

# Run tests with race detector
test-race:
	GO111MODULE=on go test -race -mod=readonly ./...

# Run tests with coverage
test-cover:
	GO111MODULE=on go test -coverprofile=coverage.out -mod=readonly ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests with czmq compatibility layer (requires czmq/v4)
test-czmq:
	GO111MODULE=on go test -tags czmq4 -mod=readonly ./...

# Run tests with czmq compatibility layer and verbose output
test-czmq-verbose:
	GO111MODULE=on go test -tags czmq4 -v -mod=readonly ./...

# Run tests with czmq compatibility layer and race detector
test-czmq-race:
	GO111MODULE=on go test -tags czmq4 -race -mod=readonly ./...

# Run benchmarks
bench:
	GO111MODULE=on go test -bench=. -benchmem -mod=readonly ./...

# Run linter (requires golangci-lint)
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/"; exit 1; }
	golangci-lint run

# Format code
fmt:
	GO111MODULE=on go fmt ./...
	GO111MODULE=on go mod tidy

# Clean build artifacts
clean:
	rm -f coverage.out coverage.html
	go clean -cache

# Run examples
example-hwserver:
	GO111MODULE=on go run -mod=readonly example/hwserver.go

example-hwclient:
	GO111MODULE=on go run -mod=readonly example/hwclient.go

example-psenvpub:
	GO111MODULE=on go run -mod=readonly example/psenvpub.go

example-psenvsub:
	GO111MODULE=on go run -mod=readonly example/psenvsub.go

# Help
help:
	@echo "Available targets:"
	@echo "  make build        - Build the project"
	@echo "  make test         - Run tests"
	@echo "  make test-verbose - Run tests with verbose output"
	@echo "  make test-race    - Run tests with race detector"
	@echo "  make test-cover   - Run tests with coverage report"
	@echo "  make test-czmq    - Run tests with czmq compatibility layer"
	@echo "  make test-czmq-verbose - Run tests with czmq (verbose)"
	@echo "  make test-czmq-race    - Run tests with czmq (race detector)"
	@echo "  make bench        - Run benchmarks"
	@echo "  make lint         - Run linter (requires golangci-lint)"
	@echo "  make fmt          - Format code and tidy modules"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make help         - Show this help message"
	@echo ""
	@echo "Example targets:"
	@echo "  make example-hwserver - Run hello world server"
	@echo "  make example-hwclient - Run hello world client"