# Build the cards application
build-cards:
	@echo "Building cards application..."
	cd cmd/cards && go build .

# Run the cards application
run-cards:
	@echo "Running cards application..."
	cd cmd/cards && go run .

# Run the tests for the cards application
test-cards:
	@echo "Running cards tests..."
	cd cmd/cards && go test .

# Build basics application
build-basics:
	@echo "Building cards application..."
	cd cmd/basics && go build .

# Run the basics application
run-basics:
	@echo "Running cards application..."
	cd cmd/basics && go run .

# Run the tests for the basics application
test-basics:
	@echo "Running basics tests..."
	cd cmd/basics && go test .

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  build-cards: 	Build the cards application"
	@echo "  run-cards: 	Run the cards application"
	@echo "  test-cards:	Run the cards application tests"
	@echo "  build-basics:	Build the basics application"
	@echo "  run-basics: 	Run the basics application"
	@echo "  test-basics: 	Run the basics application tests"
