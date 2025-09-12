# Build the cards application
build-cards:
	@echo "Building cards application..."
	go build ./cmd/cards

# Run the cards application
run-cards:
	@echo "Running cards application..."
	go run ./cmd/cards

# Run the tests for the cards application
test-cards:
	@echo "Running cards tests..."
	go test ./cmd/cards

# Build basics application
build-basics:
	@echo "Building basics application..."
	go build ./cmd/basics

# Run the basics application
run-basics:
	@echo "Running basics application..."
	go run ./cmd/basics

# Run the tests for the basics application
test-basics:
	@echo "Running basics tests..."
	go test ./cmd/basics

# Build the rest api application
build-rest-api:
	@echo "Building rest api application..."
	go build ./cmd/rest-api

# Run the rest api application
run-rest-api:
	@echo "Running rest api application..."
	go run ./cmd/rest-api

# Run the tests for the rest api application
test-rest-api:
	@echo "Running rest api application tests..."
	go test ./cmd/rest-api

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  build-cards: 	Build the cards application"
	@echo "  run-cards: 	Run the cards application"
	@echo "  test-cards:	Run the cards application tests"
	@echo "  build-basics:	Build the basics application"
	@echo "  run-basics: 	Run the basics application"
	@echo "  test-basics: 	Run the basics application tests"
	@echo "  build-rest-api:	Build the rest api application"
	@echo "  run-rest-api: 	Run the rest api application"
	@echo "  test-rest-api: 	Run the rest api application tests"