package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	orderpb "go-learning/pkg/grpc/order"
	userpb "go-learning/pkg/grpc/user"

	"google.golang.org/grpc"
)

var mu sync.Mutex // protect concurrent access

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &userServer{})
	orderpb.RegisterOrderServiceServer(grpcServer, &orderServer{})

	fmt.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
