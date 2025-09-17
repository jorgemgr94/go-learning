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