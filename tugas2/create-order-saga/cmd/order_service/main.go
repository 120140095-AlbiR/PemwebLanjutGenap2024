package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	orderservice "create-order-saga/internal/order"
	orderpb "create-order-saga/proto/order"
)

const (
	port = ":50051" // Port for the Order service
)

func main() {
	log.Printf("Starting Order Service on port %s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Create an instance of our Order service implementation
	orderServer := orderservice.NewServer()

	// Register the Order service with the gRPC server
	orderpb.RegisterOrderServiceServer(s, orderServer)

	log.Printf("Order Service listening at %v", lis.Addr())
	// Start serving requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
