package payment_client

import (
	"encoding/json"
	"log"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/webhook"
	"main.go/internal/models"
)

var (
	SUCCESSURL = "localhost:8001/api/v1/health"
	CANCELURL  = "localhost:8001/api/v1/cancel"
)

type PaymentIntentRes struct {
	Status          models.PaymentStatus
	PaymentIntentId string
	SessionId       string
}

var stripeClient *Client

// Client struct holds the Stripe secret key and webhook secret key.
type Client struct {
	StripeSecretKey     string
	StripeWebhookSecret string
}

// NewClient creates a new StripeClient instance with provided keys.
func NewStripeClient(secretKey, webhookSecret string) PaymentClient {
	if stripeClient == nil {
		stripeClient = &Client{
			StripeSecretKey:     secretKey,
			StripeWebhookSecret: webhookSecret,
		}
	}
	return stripeClient
}

func GetStripeClient() *Client {
	return stripeClient
}

// CreateCheckoutSession creates a Stripe Checkout session and returns the session ID.
func (c *Client) CreateCheckoutSession(orderID string, totalAmount int64) (string, error) {
	// Initialize Stripe client with the secret key
	stripe.Key = c.StripeSecretKey

	// Create Checkout Session parameters
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String(string(stripe.CurrencyUSD)),
					UnitAmount: stripe.Int64(totalAmount), // Amount in cents
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Order " + orderID),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(SUCCESSURL),
		CancelURL:  stripe.String(CANCELURL),
	}

	// Create session with the provided parameters
	session, err := session.New(params)
	if err != nil {
		log.Printf("Error creating Checkout session: %v", err)
		return "", err
	}

	// Return the session ID
	return session.ID, nil
}

// HandleWebhook processes the Stripe webhook event by verifying the signature
// and handling the event types like `checkout.session.completed`.
func (c *Client) HandleWebhook(payload []byte, sigHeader string) (*PaymentIntentRes, error) {
	// Verify the webhook signature using the Stripe secret key
	event, err := webhook.ConstructEvent(payload, sigHeader, c.StripeWebhookSecret)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		return nil, err
	}

	// Handle specific event types
	switch event.Type {
	case "checkout.session.completed":
		// Handle successful payment session completion
		res, err := c.handleCheckoutSessionCompleted(&event)
		if err != nil {
			log.Printf("Error handling checkout.session.completed: %v", err)
			return nil, err
		}
		return res, nil

	// You can handle more event types here (like payment_intent.succeeded)
	default:
		log.Printf("Unhandled event type: %s", event.Type)
		return nil, nil
	}

}

func handlePaymentIntentForStripe(paymentIntent *stripe.PaymentIntent) *PaymentIntentRes {
	res := PaymentIntentRes{}
	switch paymentIntent.Status {
	case stripe.PaymentIntentStatusCanceled:
		res.Status = models.PaymentCancelled
	case stripe.PaymentIntentStatusSucceeded:
		res.Status = models.PaymentCompleted
	case stripe.PaymentIntentStatusProcessing:
		res.Status = models.PaymentProcessing
	}
	res.PaymentIntentId = paymentIntent.ID
	return &res
}

// handleCheckoutSessionCompleted processes the event for a completed checkout session.
func (c *Client) handleCheckoutSessionCompleted(event *stripe.Event) (*PaymentIntentRes, error) {
	// Extract the Checkout Session from the event
	dataBytes, err := json.Marshal(event.Data.Object)
	if err != nil {
		return nil, err
	}
	var session stripe.CheckoutSession
	if err := json.Unmarshal(dataBytes, &session); err != nil {
		return nil, err
	}

	// Log or process session completion, e.g., update order status, trigger email notifications
	log.Printf("Checkout session completed: %s", session.ID)

	// Optionally, you can fetch the PaymentIntent to check if payment was successful
	paymentIntent := session.PaymentIntent
	if err != nil {
		log.Printf("Error retrieving payment intent: %v", err)
		return nil, err
	}

	// Log payment status
	log.Printf("PaymentIntent status: %s", paymentIntent.Status)

	// Further processing based on payment intent status (e.g., update order, send email, etc.)
	if paymentIntent.Status == stripe.PaymentIntentStatusSucceeded {
		log.Printf("Payment succeeded for session: %s", session.ID)
		// Optionally: Update order status in the database to 'paid', notify user, etc.
	} else {
		log.Printf("Payment failed for session: %s", session.ID)
		// Optionally: Update order status to 'failed' and notify user
	}
	res := handlePaymentIntentForStripe(paymentIntent)
	res.SessionId = session.ID
	return res, nil
}
