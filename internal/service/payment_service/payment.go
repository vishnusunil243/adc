package payment_service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"main.go/common"
	"main.go/common/payment_client"
	"main.go/internal/models"
	"main.go/internal/repository/order"
	"main.go/internal/repository/payment"
	"main.go/internal/service"
)

// PaymentService defines the interface for the payment service.
type PaymentService interface {
	CreateCheckoutSessionService(ctx context.Context, orderID string) (*CheckoutSessionRes, *service.ServiceError)
	HandleWebhookService(ctx context.Context, payload []byte, sigHeader string) *service.ServiceError
}

type paymentService struct {
	paymentRepo *payment.PaymentRepo
	orderRepo   *order.OrderRepo
}

// NewPaymentService creates a new instance of PaymentService.
func NewPaymentService() PaymentService {
	return &paymentService{
		paymentRepo: payment.NewPaymentRepo(),
		orderRepo:   order.NewOrderRepo(),
	}
}

// CreateCheckoutSessionService is a wrapper for CreateCheckoutSession in the client layer.
// It creates a new checkout session and returns the session ID.
func (ps *paymentService) CreateCheckoutSessionService(ctx context.Context, orderId string) (*CheckoutSessionRes, *service.ServiceError) {
	// Call the client's CreateCheckoutSession method
	order, repoErr := ps.orderRepo.Get(ctx, &order.GetOrderRequest{
		Id: orderId,
	})
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to get order")
	}
	totalAmount := order.Total
	sessionID, err := payment_client.GetStripeClient().CreateCheckoutSession(orderId, totalAmount)
	if err != nil {
		log.Printf("Error in creating checkout session: %v", err)
		return nil, service.NewServiceError("failed to create checkout session", common.ErrCodeInternal, http.StatusBadRequest)
	}

	order.UpdatePaymentSessionId(sessionID)
	_, repoErr = ps.orderRepo.Save(ctx, order)
	if repoErr != nil {
		return nil, service.HandleRepoErr(repoErr, "Failed to update order")
	}

	// Return the session ID if creation was successful
	return NewCheckoutSessionRes(sessionID), nil
}

// HandleWebhookService is a wrapper for HandleWebhook in the client layer.
// It processes the Stripe webhook and returns the payment intent response.
func (ps *paymentService) HandleWebhookService(ctx context.Context, payload []byte, sigHeader string) *service.ServiceError {
	// Call the client's HandleWebhook method
	result, err := payment_client.GetStripeClient().HandleWebhook(payload, sigHeader)
	if err != nil {
		log.Printf("Error in handling webhook: %v", err)
		return service.NewServiceError(fmt.Sprintf("Failed to handle webhook with err : %s", err.Error()), common.ErrCodeInternal, http.StatusInternalServerError)
	}
	order, repoErr := ps.orderRepo.Get(ctx, &order.GetOrderRequest{
		PaymentSessionId: result.SessionId,
	})
	if repoErr != nil {
		return service.HandleRepoErr(repoErr, "Failed to get order")
	}
	paymentObj, repoErr := ps.paymentRepo.Get(ctx, &payment.GetPaymentRequest{
		OrderId: order.Id,
	})
	if repoErr != nil && repoErr.ErrorCode != common.ErrCodeNotFound {
		return service.HandleRepoErr(repoErr, "Failed to get payment obj")
	}
	if paymentObj != nil {
		paymentObj.PaymentIntentId = result.PaymentIntentId
		paymentObj.PaymentStatus = result.Status
		_, repoErr = ps.paymentRepo.Save(ctx, paymentObj)
		if repoErr != nil {
			return service.HandleRepoErr(repoErr, "Failed to update order")
		}
	} else {
		newPaymentObj := models.NewPayment(order.Id, result.PaymentMethod, result.PaymentIntentId, "", models.Stripe)
		_, repoErr := ps.paymentRepo.CreatePayment(ctx, newPaymentObj)
		if repoErr != nil {
			return service.HandleRepoErr(repoErr, "Failed to create order")
		}
	}

	// Return the result containing the payment status
	return nil
}
