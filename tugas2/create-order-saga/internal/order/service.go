package order

import (
	"context"
	"log"

	commonpb "create-order-saga/proto/common"
	orderpb "create-order-saga/proto/order"
	"sync" // For safe concurrent map access

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the OrderServiceServer interface.
type Server struct {
	orderpb.UnimplementedOrderServiceServer // Embed for forward compatibility
	orders                                  map[string]*orderpb.Order
	mu                                      sync.RWMutex // Mutex to protect the orders map
}

// NewServer creates a new Order service server.
func NewServer() *Server {
	return &Server{
		orders: make(map[string]*orderpb.Order),
	}
}

// CreateOrder handles the creation of a new order.
// In a real implementation, this would persist the order to a database.
func (s *Server) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	log.Printf("Received CreateOrder request for user: %s", req.Details.UserId)

	// 1. Generate a unique order ID (e.g., using UUID)
	//    For simplicity, we'll use a placeholder.
	orderID := "order-" + req.Details.UserId // Replace with actual ID generation

	// 2. Create the order object (in memory for now)
	newOrder := &orderpb.Order{
		Id:     orderID,
		UserId: req.Details.UserId,
		Items:  req.Details.Items,
		// Calculate total amount based on items
		TotalAmount: calculateTotal(req.Details.Items),
		Status:      orderpb.OrderStatus_PENDING, // Initial status
	}

	// 3. Persist the order
	s.mu.Lock()
	s.orders[orderID] = newOrder
	s.mu.Unlock()
	log.Printf("Order %s created and stored with status PENDING", orderID)

	// 4. Return the response
	return &orderpb.CreateOrderResponse{
		OrderId: &commonpb.OrderID{Id: orderID},
		Status:  newOrder.Status,
	}, nil
}

// CancelOrder handles the compensation action for cancelling an order.
// In a real implementation, this would update the order status in the database.
func (s *Server) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*commonpb.CompensationResponse, error) {
	orderID := req.OrderId.Id
	log.Printf("Received CancelOrder request for order ID: %s", orderID)

	// 1. Find the order (e.g., order, exists := s.orders[orderID])
	// 1. Find the order
	s.mu.Lock()
	order, exists := s.orders[orderID]
	if !exists {
		s.mu.Unlock()
		log.Printf("CancelOrder failed: Order %s not found", orderID)
		return nil, status.Errorf(codes.NotFound, "Order %s not found", orderID)
	}

	// 2. Check if cancellation is possible (e.g., already cancelled?)
	if order.Status == orderpb.OrderStatus_CANCELLED {
		s.mu.Unlock()
		log.Printf("CancelOrder skipped: Order %s already cancelled", orderID)
		// Return success as the desired state is achieved (idempotency)
		return &commonpb.CompensationResponse{Success: true, Message: "Order already cancelled"}, nil
	}

	// 3. Update the order status to CANCELLED
	order.Status = orderpb.OrderStatus_CANCELLED
	s.mu.Unlock() // Unlock before logging potentially slow operations
	log.Printf("Order %s status updated to CANCELLED", orderID)

	// 4. Return success response
	return &commonpb.CompensationResponse{
		Success: true,
		Message: "Order cancelled successfully",
	}, nil

	// Example error handling:
	// if !exists {
	// 	return nil, status.Errorf(codes.NotFound, "Order %s not found", orderID)
	// }
	// return nil, status.Errorf(codes.Internal, "Failed to cancel order %s", orderID)
}

// CompleteOrder marks an order as completed in the storage.
func (s *Server) CompleteOrder(ctx context.Context, req *orderpb.CompleteOrderRequest) (*commonpb.CompensationResponse, error) {
	orderID := req.OrderId.Id
	log.Printf("Received CompleteOrder request for order ID: %s", orderID)

	s.mu.Lock()
	order, exists := s.orders[orderID]
	if !exists {
		s.mu.Unlock()
		log.Printf("CompleteOrder failed: Order %s not found", orderID)
		// This might indicate an issue if the orchestrator thinks it succeeded but the record is gone
		return nil, status.Errorf(codes.NotFound, "Order %s not found", orderID)
	}

	// Update status only if it makes sense (e.g., was PENDING)
	if order.Status == orderpb.OrderStatus_PENDING {
		order.Status = orderpb.OrderStatus_COMPLETED
		log.Printf("Order %s status updated to COMPLETED", orderID)
	} else {
		log.Printf("CompleteOrder skipped: Order %s status was %s, not PENDING", orderID, order.Status)
	}
	s.mu.Unlock()

	return &commonpb.CompensationResponse{
		Success: true,
		Message: "Order completion processed", // Indicate processed, even if status wasn't PENDING
	}, nil
}

// Helper function to calculate total amount (replace with actual logic)
func calculateTotal(items []*commonpb.Item) float32 {
	var total float32 = 0.0
	for _, item := range items {
		total += item.Price * float32(item.Quantity)
	}
	return total
}
