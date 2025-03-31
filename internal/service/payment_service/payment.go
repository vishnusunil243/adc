package payment_service

import (
	"log"

	"main.go/common/payment_client"
	"main.go/internal/repository/payment"
)

// PaymentService defines the interface for the payment service.
type PaymentService interface {
	CreateCheckoutSessionService(orderID string, totalAmount int64) (string, error)
	HandleWebhookService(payload []byte, sigHeader string) (*payment_client.PaymentIntentRes, error)
}

type paymentService struct {
	paymentRepo *payment.PaymentRepo
}

// NewPaymentService creates a new instance of PaymentService.
func NewPaymentService() PaymentService {
	return &paymentService{
		paymentRepo: payment.NewPaymentRepo(),
	}
}

// CreateCheckoutSessionService is a wrapper for CreateCheckoutSession in the client layer.
// It creates a new checkout session and returns the session ID.
func (ps *paymentService) CreateCheckoutSessionService(orderID string, totalAmount int64) (string, error) {
	// Call the client's CreateCheckoutSession method
	sessionID, err := payment_client.GetStripeClient().CreateCheckoutSession(orderID, totalAmount)
	if err != nil {
		log.Printf("Error in creating checkout session: %v", err)
		return "", err
	}
	// Return the session ID if creation was successful
	return sessionID, nil
}

// HandleWebhookService is a wrapper for HandleWebhook in the client layer.
// It processes the Stripe webhook and returns the payment intent response.
func (ps *paymentService) HandleWebhookService(payload []byte, sigHeader string) (*payment_client.PaymentIntentRes, error) {
	// Call the client's HandleWebhook method
	result, err := payment_client.GetStripeClient().HandleWebhook(payload, sigHeader)
	if err != nil {
		log.Printf("Error in handling webhook: %v", err)
		return nil, err
	}

	// Return the result containing the payment status
	return result, nil
}
