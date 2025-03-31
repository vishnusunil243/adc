package payment_client

// StripeClient interface defines methods for interacting with Stripe.
type PaymentClient interface {
	CreateCheckoutSession(orderID string, totalAmount int64) (string, error)
	HandleWebhook(payload []byte, sigHeader string) (*PaymentIntentRes, error)
}
