package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderpb "go-learning/pkg/grpc/order"
	userpb "go-learning/pkg/grpc/user"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create clients
	userClient := userpb.NewUserServiceClient(conn)
	orderClient := orderpb.NewOrderServiceClient(conn)

	// Context purpose is to prevent:
	// - Resource leaks
	// - Orphaned operations
	// - Cascading delays
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Create a user
	createUserResp, err := userClient.CreateUser(ctx, &userpb.CreateUserRequest{
		Name:  "Jorge",
		Email: "jorge@example.com",
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	log.Printf("Created User: ID=%s, Status=%s", createUserResp.GetId(), createUserResp.GetStatus().GetMessage())

	// Fetch the user
	userResp, err := userClient.GetUser(ctx, &userpb.GetUserRequest{Id: createUserResp.GetId()})
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	log.Printf("Fetched User: %s (%s)", userResp.GetName(), userResp.GetEmail())

	// Create an order
	createOrderResp, err := orderClient.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId:     createUserResp.GetId(),
		ProductIds: []string{"p1", "p2", "p3"},
	})
	if err != nil {
		log.Fatalf("CreateOrder failed: %v", err)
	}
	log.Printf("Created Order: ID=%s, Status=%s", createOrderResp.GetId(), createOrderResp.GetStatus().GetMessage())

	// Fetch the order
	orderResp, err := orderClient.GetOrder(ctx, &orderpb.GetOrderRequest{Id: createOrderResp.GetId()})
	if err != nil {
		log.Fatalf("GetOrder failed: %v", err)
	}
	log.Printf("Fetched Order: %s, Amount=%.2f, Products=%v", orderResp.GetId(), orderResp.GetAmount(), orderResp.GetProductIds())
}
