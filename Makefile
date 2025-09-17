# Cards application
build-cards:
	@echo "Building cards application..."
	go build ./cmd/cards
run-cards:
	@echo "Running cards application..."
	go run ./cmd/cards
test-cards:
	@echo "Running cards tests..."
	go test ./cmd/cards

# Basics application
build-basics:
	@echo "Building basics application..."
	go build ./cmd/basics
run-basics:
	@echo "Running basics application..."
	go run ./cmd/basics
test-basics:
	@echo "Running basics tests..."
	go test ./cmd/basics

# Rest api application
build-rest-api:
	@echo "Building rest api application..."
	go build ./cmd/rest-api
run-rest-api:
	@echo "Running rest api application..."
	go run ./cmd/rest-api
test-rest-api:
	@echo "Running rest api application tests..."
	go test ./cmd/rest-api

# Grpc application
compile-grpc:
	@echo "Compiling grpc proto files..."
	make clean-grpc
	protoc --go_out=. --go-grpc_out=. \
		--go_opt=module=go-learning \
		--go-grpc_opt=module=go-learning \
		api/proto/common/*.proto \
		api/proto/user/*.proto \
		api/proto/order/*.proto

clean-grpc:
	@echo "Cleaning generated grpc files..."
	rm -rf pkg/grpc/

run-grpc-server:
	@echo "Running grpc server..."
	go run ./cmd/grpc/server

run-grpc-client:
	@echo "Running grpc client..."
	go run ./cmd/grpc/client

# Development and testing
test-all:
	@echo "Running all tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Running linter..."
	golangci-lint run

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/ coverage.out coverage.html
	make clean-grpc

help:
	@echo "Go Learning Project - Available Commands"
	@echo "======================================"
	@echo ""
	@echo "📦 Cards Application:"
	@echo "  build-cards     - Build the cards application"
	@echo "  run-cards       - Run the cards application"
	@echo "  test-cards      - Run cards application tests"
	@echo ""
	@echo "🎯 Basics Application:"
	@echo "  build-basics    - Build the basics application"
	@echo "  run-basics      - Run the basics application"
	@echo "  test-basics     - Run basics application tests"
	@echo ""
	@echo "🌐 REST API Application:"
	@echo "  build-rest-api  - Build the REST API application"
	@echo "  run-rest-api    - Run the REST API application"
	@echo "  test-rest-api   - Run REST API application tests"
	@echo ""
	@echo "🔌 gRPC Application:"
	@echo "  compile-grpc    - Compile gRPC proto files to Go code"
	@echo "  clean-grpc      - Remove generated gRPC files"
	@echo "  run-grpc-server - Run the gRPC server"
	@echo "  run-grpc-client - Run the gRPC client"
	@echo ""
	@echo "🔧 Development & Testing:"
	@echo "  test-all        - Run all tests with coverage report"
	@echo "  lint            - Run Go linter (golangci-lint)"
	@echo "  clean           - Clean build artifacts and generated files"
	@echo ""
	@echo "❓ Help:"
	@echo "  help            - Show this help message"
	@echo ""
	@echo "Usage: make <command>"
	@echo "Example: make run-rest-api"

# Semantic versioning
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# Build flags for version info
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Build with version info
build-with-version:
	@echo "Building with version $(VERSION)..."
	go build $(LDFLAGS) ./cmd/cards
	go build $(LDFLAGS) ./cmd/basics
	go build $(LDFLAGS) ./cmd/rest-api
	go build $(LDFLAGS) ./cmd/grpc/server

# Check conventional commits
check-commits:
	@echo "Checking commit messages..."
	@git log --pretty=format:"%s" origin/main..HEAD | while read commit; do \
		if [[ ! "$$commit" =~ ^(feat|fix|perf|refactor|docs|style|test|chore|ci|build)(\(.+\))?(!)?: .+ ]]; then \
			echo "❌ Invalid commit: $$commit"; \
			echo "   Must follow: type(scope): description"; \
			exit 1; \
		else \
			echo "✅ Valid commit: $$commit"; \
		fi; \
	done

# Predict next version
predict-version:
	@echo "Analyzing commits for version prediction..."
	@COMMITS=$$(git log --pretty=format:"%s" origin/main..HEAD); \
	BREAKING=$$(echo "$$COMMITS" | grep -E "^(feat|fix|perf|refactor)(\(.+\))?!: .+" | wc -l); \
	FEATURES=$$(echo "$$COMMITS" | grep -E "^feat(\(.+\))?: .+" | wc -l); \
	FIXES=$$(echo "$$COMMITS" | grep -E "^fix(\(.+\))?: .+" | wc -l); \
	echo "Breaking changes: $$BREAKING"; \
	echo "New features: $$FEATURES"; \
	echo "Bug fixes: $$FIXES"; \
	if [ $$BREAKING -gt 0 ]; then \
		echo "🚨 MAJOR version bump expected"; \
	elif [ $$FEATURES -gt 0 ]; then \
		echo "🆕 MINOR version bump expected"; \
	elif [ $$FIXES -gt 0 ]; then \
		echo "🐛 PATCH version bump expected"; \
	else \
		echo "📝 NO version bump expected"; \
	fi