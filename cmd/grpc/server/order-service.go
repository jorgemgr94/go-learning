package main

import (
	"context"
	"fmt"
	common "go-learning/pkg/grpc/common"
	orderpb "go-learning/pkg/grpc/order"
)

var orderStore = make(map[string]*orderpb.GetOrderReply)

type orderServer struct {
	orderpb.UnimplementedOrderServiceServer
}

func (s *orderServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderReply, error) {
	mu.Lock()
	defer mu.Unlock()

	order, exists := orderStore[req.GetId()]
	if !exists {
		return &orderpb.GetOrderReply{
			Status: &common.ResponseStatus{Code: 404, Message: "Order not found"},
		}, nil
	}
	return order, nil
}

func (s *orderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderReply, error) {
	mu.Lock()
	defer mu.Unlock()

	// fake ID generation
	id := fmt.Sprintf("o%d", len(orderStore)+1)
	order := &orderpb.GetOrderReply{
		Id:         id,
		Amount:     float64(len(req.GetProductIds())) * 100.0, // fake pricing
		ProductIds: req.GetProductIds(),
		Status: &common.ResponseStatus{
			Code:    201,
			Message: "Order created",
		},
	}
	orderStore[id] = order

	return &orderpb.CreateOrderReply{
		Id:     id,
		Status: &common.ResponseStatus{Code: 201, Message: "Order created successfully"},
	}, nil
}
