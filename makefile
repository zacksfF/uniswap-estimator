.PHONY: build run test benchmark clean help

APP_NAME=uniswap-estimator

# Default - show help
help:
	@echo "Uniswap Estimator Commands:"
	@echo "  make run       - Run the application"
	@echo "  make build     - Build the application"
	@echo "  make test      - Run tests"
	@echo "  make benchmark - Run math benchmarks (1inch requirement)"
	@echo "  make clean     - Clean build files"

# Install dependencies and run
run:
	@echo "Starting server..."
	@go mod tidy
	@go run cmd/main.go

# Build the app
build:
	@echo "Building..."
	@go build -o $(APP_NAME) cmd/main.go

# Run all tests
test:
	@echo "Running tests..."
	@go test ./...

# Benchmark math functions (important for 1inch!)
benchmark:
	@echo "Running math benchmarks..."
	@go test -bench=. -benchmem ./test/

# Clean up
clean:
	@echo "Cleaning..."
	@rm -f $(APP_NAME)
	@go clean

