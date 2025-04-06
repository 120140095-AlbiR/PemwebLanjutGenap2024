package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	paymentservice "create-order-saga/internal/payment"
	paymentpb "create-order-saga/proto/payment"
)

const (
	port = ":50052" // Port for the Payment service (different from Order service)
)

func main() {
	log.Printf("Starting Payment Service on port %s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Create an instance of our Payment service implementation
	paymentServer := paymentservice.NewServer()

	// Register the Payment service with the gRPC server
	paymentpb.RegisterPaymentServiceServer(s, paymentServer)

	log.Printf("Payment Service listening at %v", lis.Addr())
	// Start serving requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
