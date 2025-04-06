package grpc_clients

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // Use insecure for example only

	orderpb "create-order-saga/proto/order"
	paymentpb "create-order-saga/proto/payment"
	shippingpb "create-order-saga/proto/shipping"
)

// ServiceClients holds clients for all required services.
type ServiceClients struct {
	Order    orderpb.OrderServiceClient
	Payment  paymentpb.PaymentServiceClient
	Shipping shippingpb.ShippingServiceClient
}

// NewServiceClients creates and returns gRPC clients for the saga services.
func NewServiceClients(orderAddr, paymentAddr, shippingAddr string) (*ServiceClients, error) {
	// Establish connection to Order Service
	orderConn, err := grpc.Dial(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to Order Service at %s: %v", orderAddr, err)
		return nil, err
	}
	orderClient := orderpb.NewOrderServiceClient(orderConn)
	log.Printf("Connected to Order Service at %s", orderAddr)

	// Establish connection to Payment Service
	paymentConn, err := grpc.Dial(paymentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to Payment Service at %s: %v", paymentAddr, err)
		// Consider closing orderConn here if needed
		return nil, err
	}
	paymentClient := paymentpb.NewPaymentServiceClient(paymentConn)
	log.Printf("Connected to Payment Service at %s", paymentAddr)

	// Establish connection to Shipping Service
	shippingConn, err := grpc.Dial(shippingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to Shipping Service at %s: %v", shippingAddr, err)
		// Consider closing orderConn and paymentConn here if needed
		return nil, err
	}
	shippingClient := shippingpb.NewShippingServiceClient(shippingConn)
	log.Printf("Connected to Shipping Service at %s", shippingAddr)

	return &ServiceClients{
		Order:    orderClient,
		Payment:  paymentClient,
		Shipping: shippingClient,
	}, nil

	// Note: Connections should ideally be closed gracefully when the application shuts down.
	// This basic example doesn't include connection closing logic.
}
