package shipping

import (
	"context"
	"log"
	"math/rand" // For simulating success/failure

	commonpb "create-order-saga/proto/common"
	shippingpb "create-order-saga/proto/shipping"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the ShippingServiceServer interface.
type Server struct {
	shippingpb.UnimplementedShippingServiceServer // Embed for forward compatibility
	shipments                                     map[string]*shippingpb.Shipment
	mu                                            sync.RWMutex
}

// NewServer creates a new Shipping service server.
func NewServer() *Server {
	return &Server{
		shipments: make(map[string]*shippingpb.Shipment),
	}
}

// ArrangeShipping handles arranging shipping for an order.
// Simulates success or failure.
func (s *Server) ArrangeShipping(ctx context.Context, req *shippingpb.ArrangeShippingRequest) (*shippingpb.ArrangeShippingResponse, error) {
	orderID := req.OrderId.Id
	log.Printf("Received ArrangeShipping request for order ID: %s, Address: %s", orderID, req.Address.City)

	// 1. Generate a unique shipment ID
	shipmentID := "ship-" + orderID // Replace with actual ID generation

	// 2. Simulate shipping arrangement (e.g., call a carrier API)
	//    Randomly succeed or fail for demonstration purposes.
	succeeded := rand.Intn(10) > 1 // 80% chance of success

	if !succeeded {
		log.Printf("Failed to arrange shipping for order %s (simulated failure)", orderID)
		// Return a gRPC error to signal failure to the orchestrator
		return nil, status.Errorf(codes.Internal, "Failed to arrange shipping for order %s: Carrier unavailable", orderID)
	}

	// 3. Create and persist shipment record (in memory for now)
	newShipment := &shippingpb.Shipment{
		Id:      shipmentID,
		OrderId: req.OrderId,
		Address: req.Address,
		Status:  shippingpb.ShippingStatus_PENDING, // Initial status
		// TrackingNumber: // Get from carrier API if successful
	}
	// --- Modified Logic ---
	// Set status directly to SHIPPED on success
	newShipment.Status = shippingpb.ShippingStatus_SHIPPED

	// Persist
	s.mu.Lock()
	s.shipments[shipmentID] = newShipment
	s.mu.Unlock()
	log.Printf("Shipment %s created and stored for order %s with status SHIPPED. Record: %+v", shipmentID, orderID, newShipment)

	// 4. Return response with SHIPPED status
	return &shippingpb.ArrangeShippingResponse{
		ShipmentId: shipmentID,
		Status:     newShipment.Status, // Should be SHIPPED
	}, nil
}

// CancelShipping handles the compensation action for cancelling shipping.
func (s *Server) CancelShipping(ctx context.Context, req *shippingpb.CancelShippingRequest) (*commonpb.CompensationResponse, error) {
	orderID := req.OrderId.Id
	shipmentID := req.ShipmentId
	log.Printf("Received CancelShipping request for order ID: %s, Shipment ID: %s", orderID, shipmentID)

	// 1. Find the shipment record (e.g., shipment, exists := s.shipments[shipmentID])
	//    Ensure it belongs to the correct orderID.
	// 1. Find the shipment record
	s.mu.Lock()
	shipment, exists := s.shipments[shipmentID]
	if !exists {
		s.mu.Unlock()
		log.Printf("CancelShipping failed: Shipment %s not found", shipmentID)
		return nil, status.Errorf(codes.NotFound, "Shipment %s not found", shipmentID)
	}
	// Optional: Verify order ID
	if shipment.OrderId.Id != orderID {
		s.mu.Unlock()
		log.Printf("CancelShipping failed: Shipment %s does not belong to order %s", shipmentID, orderID)
		return nil, status.Errorf(codes.InvalidArgument, "Shipment %s does not belong to order %s", shipmentID, orderID)
	}

	// 2. Check if cancellation is possible
	if shipment.Status == shippingpb.ShippingStatus_CANCELLED {
		s.mu.Unlock()
		log.Printf("CancelShipping skipped: Shipment %s already cancelled", shipmentID)
		return &commonpb.CompensationResponse{Success: true, Message: "Shipment already cancelled"}, nil
	}
	// In a real system, you might prevent cancelling if already SHIPPED,
	// but for this example, we allow setting to CANCELLED from SHIPPED.
	// if shipment.Status == shippingpb.ShippingStatus_SHIPPED {
	// 	 s.mu.Unlock()
	// 	 log.Printf("CancelShipping failed: Shipment %s already shipped", shipmentID)
	// 	 return nil, status.Errorf(codes.FailedPrecondition, "Cannot cancel already shipped shipment %s", shipmentID)
	// }

	// 3. Perform cancellation action (simulation)
	// Assume cancellation is successful for this example.

	// 4. Update shipment status to CANCELLED
	shipment.Status = shippingpb.ShippingStatus_CANCELLED
	s.mu.Unlock() // Unlock before logging
	log.Printf("Shipment %s for order %s status updated to CANCELLED.", shipmentID, orderID)

	// 5. Return success response
	return &commonpb.CompensationResponse{
		Success: true,
		Message: "Shipping cancelled successfully",
	}, nil

	// Example error handling:
	// if !exists {
	// 	return nil, status.Errorf(codes.NotFound, "Shipment %s not found", shipmentID)
	// }
	// if shipment.Status == shippingpb.ShippingStatus_SHIPPED {
	//  return nil, status.Errorf(codes.FailedPrecondition, "Cannot cancel already shipped shipment %s", shipmentID)
	// }
	// return nil, status.Errorf(codes.Internal, "Failed to cancel shipment %s", shipmentID)
}
