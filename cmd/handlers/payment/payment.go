package payment_handlers

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"main.go/common"
	"main.go/internal/service"
	"main.go/internal/service/payment_service"
)

type PaymentHandler struct {
	paymentService payment_service.PaymentService
}

func (p *PaymentHandler) CreateCheckoutSession(c echo.Context) error {
	// Extract orderID and totalAmount from the request body
	orderId := c.Param("order_id")
	// Call the service to create a checkout session
	sessionRes, err := p.paymentService.CreateCheckoutSessionService(c.Request().Context(), orderId)
	if err != nil {
		return common.NewResponse(c, false, err, err.StatusCode)
	}

	// Return the session ID in the response
	return common.NewResponse(c, true, sessionRes, http.StatusOK)
}

func (p *PaymentHandler) HandleWebhook(c echo.Context) error {
	// Ensure the body is closed after reading by using defer
	defer c.Request().Body.Close()

	// Read the raw body for webhook verification and signature
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return common.NewResponse(c, false,
			service.NewServiceError("Failed to read request body", common.ErrCodeInternal, http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// Read the Stripe-Signature header for webhook verification
	signatureHeader := c.Request().Header.Get("Stripe-Signature")

	// Call the service to handle the webhook
	serErr := p.paymentService.HandleWebhookService(c.Request().Context(), payload, signatureHeader)
	if serErr != nil {
		return common.NewResponse(c, false, serErr, serErr.StatusCode)
	}

	// Return success response
	return common.NewResponse(c, true, "Webhook processed successfully", http.StatusOK)
}

// NewPaymentHandler creates a new instance of PaymentHandler
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService: payment_service.NewPaymentService(),
	}
}
