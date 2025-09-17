package main

import (
	"context"
	"fmt"
	common "go-learning/pkg/grpc/common"
	userpb "go-learning/pkg/grpc/user"
)

var userStore = make(map[string]*userpb.GetUserReply)

type userServer struct {
	userpb.UnimplementedUserServiceServer
}

func (s *userServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserReply, error) {
	mu.Lock()
	defer mu.Unlock()

	user, exists := userStore[req.GetId()]
	if !exists {
		return &userpb.GetUserReply{
			Status: &common.ResponseStatus{Code: 404, Message: "User not found"},
		}, nil
	}
	return user, nil
}

func (s *userServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserReply, error) {
	mu.Lock()
	defer mu.Unlock()

	// fake ID generation
	id := fmt.Sprintf("u%d", len(userStore)+1)
	user := &userpb.GetUserReply{
		Id:    id,
		Name:  req.GetName(),
		Email: req.GetEmail(),
		Status: &common.ResponseStatus{
			Code:    201,
			Message: "User created",
		},
	}
	userStore[id] = user

	return &userpb.CreateUserReply{
		Id:     id,
		Status: &common.ResponseStatus{Code: 201, Message: "User created successfully"},
	}, nil
}
