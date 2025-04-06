package main

import (
	"context"
	"log"
	"time"

	"create-order-saga/internal/orchestrator"
	"create-order-saga/pkg/grpc_clients"
	commonpb "create-order-saga/proto/common"
)

const (
	orderServiceAddr    = "localhost:50051"
	paymentServiceAddr  = "localhost:50052"
	shippingServiceAddr = "localhost:50053"
)

func main() {
	log.Println("Starting Saga Orchestrator...")

	// Connect to downstream services
	clients, err := grpc_clients.NewServiceClients(orderServiceAddr, paymentServiceAddr, shippingServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create service clients: %v", err)
	}
	// Note: Connections are not closed in this simple example.

	// Create the orchestrator instance
	sagaOrchestrator := orchestrator.NewOrchestrator(clients)

	// --- Simulate an incoming order request ---
	// In a real application, this might come from an API gateway or message queue.
	log.Println("Simulating incoming order request...")
	orderDetails := &commonpb.OrderDetails{
		UserId: "user-123",
		Items: []*commonpb.Item{
			{ProductId: "prod-A", Quantity: 2, Price: 10.50},
			{ProductId: "prod-B", Quantity: 1, Price: 25.00},
		},
	}
	paymentInfo := &commonpb.PaymentInfo{
		CardNumber: "xxxx-xxxx-xxxx-1234", // Dummy data
		ExpiryDate: "12/26",
		Cvv:        "123",
		Amount:     46.00, // 2*10.50 + 25.00
	}
	shippingAddress := &commonpb.ShippingAddress{
		Street:  "123 Saga Lane",
		City:    "Orchestration City",
		State:   "Workflow",
		ZipCode: "98765",
		Country: "GoLand",
	}

	// Execute the saga
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Set a deadline for the saga
	defer cancel()

	err = sagaOrchestrator.ExecuteCreateOrderSaga(ctx, orderDetails, paymentInfo, shippingAddress)
	if err != nil {
		log.Printf("Saga Execution Failed: %v", err)
	} else {
		log.Println("Saga Execution Completed Successfully.")
	}

	log.Println("Orchestrator finished.")
}
