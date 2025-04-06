package payment

import (
	"context"
	"log"
	"math/rand" // For simulating success/failure

	commonpb "create-order-saga/proto/common"
	paymentpb "create-order-saga/proto/payment"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the PaymentServiceServer interface.
type Server struct {
	paymentpb.UnimplementedPaymentServiceServer // Embed for forward compatibility
	payments                                    map[string]*paymentpb.Payment
	mu                                          sync.RWMutex
}

// NewServer creates a new Payment service server.
func NewServer() *Server {
	return &Server{
		payments: make(map[string]*paymentpb.Payment),
	}
}

// ProcessPayment handles processing a payment for an order.
// Simulates success or failure.
func (s *Server) ProcessPayment(ctx context.Context, req *paymentpb.ProcessPaymentRequest) (*paymentpb.ProcessPaymentResponse, error) {
	orderID := req.OrderId.Id
	log.Printf("Received ProcessPayment request for order ID: %s, Amount: %.2f", orderID, req.PaymentInfo.Amount)

	// 1. Generate a unique payment ID
	paymentID := "pay-" + orderID // Replace with actual ID generation

	// 2. Simulate payment processing (e.g., call a payment gateway)
	//    Randomly succeed or fail for demonstration purposes.
	succeeded := rand.Intn(10) > 2 // 70% chance of success

	paymentStatus := paymentpb.PaymentStatus_FAILED
	message := "Payment failed due to insufficient funds." // Example failure message
	if succeeded {
		paymentStatus = paymentpb.PaymentStatus_SUCCESS
		message = "Payment processed successfully."
		log.Printf("Payment %s for order %s succeeded.", paymentID, orderID)
	} else {
		log.Printf("Payment %s for order %s failed.", paymentID, orderID)
	}

	// 3. Create and persist payment record (in memory for now)
	newPayment := &paymentpb.Payment{
		Id:      paymentID,
		OrderId: req.OrderId,
		Amount:  req.PaymentInfo.Amount,
		Status:  paymentStatus,
		// TransactionId: // Get from gateway if successful
	}
	// Persist
	s.mu.Lock()
	s.payments[paymentID] = newPayment
	s.mu.Unlock()
	log.Printf("Payment record stored: %+v", newPayment)

	// 4. Return response
	return &paymentpb.ProcessPaymentResponse{
		PaymentId: paymentID,
		Status:    paymentStatus,
		Message:   message,
	}, nil

	// Note: In a real scenario, errors from the gateway should be handled
	// and potentially returned as gRPC errors.
	// return nil, status.Errorf(codes.Internal, "Payment gateway error")
}

// RefundPayment handles the compensation action for refunding a payment.
func (s *Server) RefundPayment(ctx context.Context, req *paymentpb.RefundPaymentRequest) (*commonpb.CompensationResponse, error) {
	orderID := req.OrderId.Id
	paymentID := req.PaymentId
	log.Printf("Received RefundPayment request for order ID: %s, Payment ID: %s", orderID, paymentID)

	// 1. Find the payment record (e.g., payment, exists := s.payments[paymentID])
	//    Ensure it belongs to the correct orderID.
	// 1. Find the payment record
	s.mu.Lock()
	payment, exists := s.payments[paymentID]
	if !exists {
		s.mu.Unlock()
		log.Printf("RefundPayment failed: Payment %s not found", paymentID)
		return nil, status.Errorf(codes.NotFound, "Payment %s not found", paymentID)
	}
	// Optional: Verify it belongs to the correct orderID
	if payment.OrderId.Id != orderID {
		s.mu.Unlock()
		log.Printf("RefundPayment failed: Payment %s does not belong to order %s", paymentID, orderID)
		return nil, status.Errorf(codes.InvalidArgument, "Payment %s does not belong to order %s", paymentID, orderID)
	}

	// 2. Check if refund is possible
	if payment.Status == paymentpb.PaymentStatus_REFUNDED {
		s.mu.Unlock()
		log.Printf("RefundPayment skipped: Payment %s already refunded", paymentID)
		return &commonpb.CompensationResponse{Success: true, Message: "Payment already refunded"}, nil
	}
	if payment.Status == paymentpb.PaymentStatus_FAILED {
		s.mu.Unlock()
		log.Printf("RefundPayment skipped: Payment %s originally failed", paymentID)
		// Arguably, this should still be success from orchestrator's perspective
		return &commonpb.CompensationResponse{Success: true, Message: "Payment originally failed, no refund needed"}, nil
	}

	// 3. Perform refund action (simulation)
	// Assume refund is successful for this example.

	// 4. Update payment status to REFUNDED
	payment.Status = paymentpb.PaymentStatus_REFUNDED
	s.mu.Unlock() // Unlock before logging
	log.Printf("Payment %s for order %s status updated to REFUNDED.", paymentID, orderID)

	// 5. Return success response
	return &commonpb.CompensationResponse{
		Success: true,
		Message: "Payment refunded successfully",
	}, nil

	// Example error handling:
	// if !exists {
	// 	return nil, status.Errorf(codes.NotFound, "Payment %s not found", paymentID)
	// }
	// return nil, status.Errorf(codes.Internal, "Failed to refund payment %s", paymentID)
}
