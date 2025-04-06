package orchestrator

import (
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc/status"

	"create-order-saga/pkg/grpc_clients"
	commonpb "create-order-saga/proto/common"
	orderpb "create-order-saga/proto/order"
	paymentpb "create-order-saga/proto/payment"
	shippingpb "create-order-saga/proto/shipping"
)

// Orchestrator manages the execution of the Create Order Saga.
type Orchestrator struct {
	clients *grpc_clients.ServiceClients
}

// NewOrchestrator creates a new saga orchestrator.
func NewOrchestrator(clients *grpc_clients.ServiceClients) *Orchestrator {
	return &Orchestrator{clients: clients}
}

// SagaState holds the intermediate results during saga execution.
type SagaState struct {
	OrderID    *commonpb.OrderID
	PaymentID  string
	ShipmentID string
}

// ExecuteCreateOrderSaga runs the distributed transaction for creating an order.
func (o *Orchestrator) ExecuteCreateOrderSaga(ctx context.Context, details *commonpb.OrderDetails, paymentInfo *commonpb.PaymentInfo, shippingAddr *commonpb.ShippingAddress) error {
	log.Println("Starting Create Order Saga...")
	state := &SagaState{}
	var err error

	// --- Step 1: Create Order ---
	log.Println("Step 1: Creating Order...")
	createOrderResp, err := o.clients.Order.CreateOrder(ctx, &orderpb.CreateOrderRequest{Details: details})
	if err != nil {
		log.Printf("Saga Failed: Step 1 (CreateOrder) failed: %v", err)
		// --- Modified Logic ---
		// Attempt compensation for consistency, even though order likely wasn't created
		o.compensateCreateOrder(state.OrderID) // state.OrderID will be nil here
		return errors.New("failed to create order")
	}
	state.OrderID = createOrderResp.OrderId // ID assigned *after* successful call
	log.Printf("Step 1 Success: Order created with ID: %s", state.OrderID.Id)

	// --- Step 2: Process Payment ---
	log.Println("Step 2: Processing Payment...")
	processPaymentReq := &paymentpb.ProcessPaymentRequest{
		OrderId:     state.OrderID,
		PaymentInfo: paymentInfo, // Use the provided payment info
	}
	processPaymentResp, err := o.clients.Payment.ProcessPayment(ctx, processPaymentReq)
	// Check for gRPC error OR explicit failure status in response
	paymentFailed := err != nil || (processPaymentResp != nil && processPaymentResp.Status == paymentpb.PaymentStatus_FAILED)

	if paymentFailed {
		log.Printf("Saga Failed: Step 2 (ProcessPayment) failed. Error: %v, Response Status: %s", err, processPaymentResp.GetStatus()) // GetStatus() is safe even if processPaymentResp is nil
		// --- Modified Logic ---
		// Also attempt to compensate the failed payment step itself
		o.compensateProcessPayment(state.OrderID, state.PaymentID) // PaymentID might be empty here

		// Compensate preceding successful steps (as before)
		o.compensateCreateOrder(state.OrderID) // Compensate Step 1
		return errors.New("failed to process payment")
	}
	// If successful:
	state.PaymentID = processPaymentResp.PaymentId // ID is assigned *after* successful call
	log.Printf("Step 2 Success: Payment processed with ID: %s", state.PaymentID)

	// --- Step 3: Arrange Shipping ---
	log.Println("Step 3: Arranging Shipping...")
	arrangeShippingReq := &shippingpb.ArrangeShippingRequest{
		OrderId: state.OrderID,
		Address: shippingAddr, // Use the provided shipping address
	}
	arrangeShippingResp, err := o.clients.Shipping.ArrangeShipping(ctx, arrangeShippingReq)
	if err != nil {
		// Check if the error is a gRPC status error (indicating service-level failure)
		grpcStatus, ok := status.FromError(err)
		if ok {
			log.Printf("Saga Failed: Step 3 (ArrangeShipping) failed with gRPC status: %s - %s", grpcStatus.Code(), grpcStatus.Message())
		} else {
			log.Printf("Saga Failed: Step 3 (ArrangeShipping) failed with non-gRPC error: %v", err)
		}
		// --- Modified Logic ---
		// Also attempt to compensate the failed shipping step itself
		o.compensateArrangeShipping(state.OrderID, state.ShipmentID) // ShipmentID might be empty here

		// Compensate preceding successful steps (as before)
		o.compensateProcessPayment(state.OrderID, state.PaymentID) // Compensate Step 2
		o.compensateCreateOrder(state.OrderID)                     // Compensate Step 1
		return errors.New("failed to arrange shipping")
	}
	state.ShipmentID = arrangeShippingResp.ShipmentId // ID is assigned *after* successful call
	log.Printf("Step 3 Success: Shipping arranged with ID: %s", state.ShipmentID)

	// --- Saga Success ---
	log.Printf("Saga Completed Successfully for Order ID: %s", state.OrderID.Id)

	// Final step: Mark the order as completed in the Order service
	log.Printf("Marking Order %s as COMPLETED...", state.OrderID.Id)
	completeCtx, completeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer completeCancel()
	_, completeErr := o.clients.Order.CompleteOrder(completeCtx, &orderpb.CompleteOrderRequest{OrderId: state.OrderID})
	if completeErr != nil {
		// Log this failure, but the core saga succeeded. Might need monitoring/alerting.
		log.Printf("WARNING: Saga succeeded, but failed to mark Order %s as COMPLETED: %v", state.OrderID.Id, completeErr)
	} else {
		log.Printf("Order %s successfully marked as COMPLETED.", state.OrderID.Id)
	}

	return nil // Return success even if the final CompleteOrder call failed (core transaction was okay)
}

// --- Compensation Functions ---

func (o *Orchestrator) compensateCreateOrder(orderID *commonpb.OrderID) {
	// Handle cases where CreateOrder failed before generating an ID
	if orderID == nil || orderID.Id == "" {
		log.Printf("Attempting Order compensation, but OrderID was not generated (step failed early). Skipping CancelOrder call.")
		return // Skip compensation if no ID was generated
	}

	log.Printf("Compensating: Cancelling Order %s", orderID.Id)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Use a background context for compensation
	defer cancel()

	_, err := o.clients.Order.CancelOrder(ctx, &orderpb.CancelOrderRequest{OrderId: orderID})
	if err != nil {
		// Log critical error: Compensation failed! Manual intervention might be needed.
		log.Printf("CRITICAL: Failed to compensate CreateOrder for Order ID %s: %v", orderID.Id, err)
	} else {
		log.Printf("Compensation Success: Order %s cancelled.", orderID.Id)
	}
}

// Note: compensateProcessPayment is now also called if ProcessPayment itself fails.
func (o *Orchestrator) compensateProcessPayment(orderID *commonpb.OrderID, paymentID string) {
	// Handle cases where ProcessPayment failed before generating an ID
	if paymentID == "" {
		log.Printf("Attempting Payment compensation for Order %s, but PaymentID was not generated (step failed early). Skipping specific RefundPayment call.", orderID.Id)
		// Depending on PaymentService implementation, RefundPayment might handle lookup by OrderID if PaymentID is empty.
		return // Skip compensation if no ID was generated
	}

	log.Printf("Compensating: Refunding Payment %s for Order %s", paymentID, orderID.Id)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.clients.Payment.RefundPayment(ctx, &paymentpb.RefundPaymentRequest{OrderId: orderID, PaymentId: paymentID})
	if err != nil {
		log.Printf("CRITICAL: Failed to compensate ProcessPayment for Order ID %s, Payment ID %s: %v", orderID.Id, paymentID, err)
	} else {
		log.Printf("Compensation Success: Payment %s refunded.", paymentID)
	}
}

// Note: compensateArrangeShipping is now also called if ArrangeShipping itself fails.
func (o *Orchestrator) compensateArrangeShipping(orderID *commonpb.OrderID, shipmentID string) {
	// Handle cases where ArrangeShipping failed before generating an ID
	if shipmentID == "" {
		log.Printf("Attempting Shipping compensation for Order %s, but ShipmentID was not generated (step failed early). Skipping specific CancelShipping call.", orderID.Id)
		// Depending on ShippingService implementation, a different compensation might be needed,
		// or CancelShipping might handle lookup by OrderID if ShipmentID is empty.
		return // Skip compensation if no ID was generated
	}

	log.Printf("Compensating: Cancelling Shipping %s for Order %s", shipmentID, orderID.Id)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.clients.Shipping.CancelShipping(ctx, &shippingpb.CancelShippingRequest{OrderId: orderID, ShipmentId: shipmentID})
	if err != nil {
		log.Printf("CRITICAL: Failed to compensate ArrangeShipping for Order ID %s, Shipment ID %s: %v", orderID.Id, shipmentID, err)
	} else {
		log.Printf("Compensation Success: Shipment %s cancelled.", shipmentID)
	}
}
