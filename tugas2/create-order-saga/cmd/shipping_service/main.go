package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	shippingservice "create-order-saga/internal/shipping"
	shippingpb "create-order-saga/proto/shipping"
)

const (
	port = ":50053" // Port for the Shipping service (different from others)
)

func main() {
	log.Printf("Starting Shipping Service on port %s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Create an instance of our Shipping service implementation
	shippingServer := shippingservice.NewServer()

	// Register the Shipping service with the gRPC server
	shippingpb.RegisterShippingServiceServer(s, shippingServer)

	log.Printf("Shipping Service listening at %v", lis.Addr())
	// Start serving requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
