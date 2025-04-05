package models

import (
	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/utils"
)

type PaymentType string

type PaymentStatus string

var (
	Stripe PaymentType = "stripe"

	PaymentCompleted  PaymentStatus = "completed"
	PaymentCancelled  PaymentStatus = "cancelled"
	PaymentPending    PaymentStatus = "pending"
	PaymentProcessing PaymentStatus = "processing"
)

type Payment struct {
	Id              string        `json:"id"`
	OrderId         string        `json:"order_id"`
	Type            PaymentType   `json:"type"`
	PaymentMethod   string        `json:"payment_method"` //id of the payment method
	PaymentStatus   PaymentStatus `json:"payment_status"`
	PaymentIntentId string        `json:"payment_intent_id"`
	*common.AuditFields
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.Id == "" {
		p.Id = utils.GenerateReadableID(16)
	}
	return nil
}

// Constructor to create a new payment
func NewPayment(orderId, paymentMethod, paymentIntentId, createdBy string, paymentType PaymentType) *Payment {
	return &Payment{
		OrderId:         orderId,
		Type:            paymentType,
		PaymentMethod:   paymentMethod,
		PaymentStatus:   PaymentPending,
		PaymentIntentId: paymentIntentId,
		AuditFields:     common.NewAuditFieldsWithCreatedBy(createdBy),
	}
}

// Function to update a payment's details
func (p *Payment) UpdatePaymentDetails(paymentMethod, paymentIntentId string, paymentStatus PaymentStatus) {
	if paymentMethod != "" {
		p.PaymentMethod = paymentMethod
	}
	if paymentStatus != "" {
		p.PaymentStatus = paymentStatus
	}
	if paymentIntentId != "" {
		p.PaymentIntentId = paymentIntentId
	}
}

// Function to update payment status (specific to status changes)
func (p *Payment) UpdatePaymentStatus(newStatus PaymentStatus) {
	if newStatus != "" {
		p.PaymentStatus = newStatus
	}
}

type ListPayment []*Payment

// Function to get a list of payment IDs from a list of payments
func (l *ListPayment) GetPaymentIds() []string {
	paymentIds := []string{}
	for _, payment := range *l {
		paymentIds = append(paymentIds, payment.Id)
	}
	return paymentIds
}

// Function to get a list of order IDs from a list of payments
func (l *ListPayment) GetOrderIds() []string {
	orderIds := []string{}
	for _, payment := range *l {
		orderIds = append(orderIds, payment.OrderId)
	}
	return orderIds
}
