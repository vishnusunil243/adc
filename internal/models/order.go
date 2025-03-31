package models

import (
	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/utils"
)

type Status string

const (
	Placed    Status = "placed"
	Shipped   Status = "shipped"
	Delivered Status = "delivered"
	Pending   Status = "pending"
	Cancelled Status = "cancelled"
	Refunded  Status = "refunded"
)

type Order struct {
	Id               string  `json:"id"`
	UserId           string  `json:"user_id"`
	Total            float64 `json:"total"`
	Status           Status  `json:"status"`
	DeliveredAt      int64   `json:"bigint"`
	ShippedAt        int64   `json:"shipped_at"`
	PlacedAt         int64   `json:"placed_at"`
	AddressId        string  `json:"address_id"`
	PaymentSessionId string  `json:"payment_session_id"`
	*common.AuditFields
}

func (o *Order) UpdateStatus(status Status) {
	if status != "" {
		o.Status = status
	}
}

func NewOrder(userId string, total float64, addressId string) *Order {
	return &Order{
		UserId:      userId,
		Total:       total,
		AddressId:   addressId,
		AuditFields: common.NewAuditFields(),
	}
}

type OrderBaseResponse struct {
	Id     string  `json:"id"`
	UserId string  `json:"user_id"`
	Total  float64 `json:"total"`
}

func NewOrderBaseResponse(order *Order) *OrderBaseResponse {
	if order == nil {
		return nil
	}
	return &OrderBaseResponse{
		Id:     order.Id,
		UserId: order.UserId,
		Total:  order.Total,
	}
}

type ListOrderResponse []*Order

func (l *ListOrderResponse) GetIds() []string {
	ids := []string{}
	for _, ord := range *l {
		ids = append(ids, ord.Id)
	}
	return ids
}

func (l *ListOrderResponse) GetUserIds() []string {
	userIds := []string{}
	for _, order := range *l {
		userIds = append(userIds, order.GetAuditFieldsUserIds()...)
	}
	return userIds
}

func (l *ListOrderResponse) ToMap() map[string]*Order {
	res := map[string]*Order{}
	for _, ord := range *l {
		res[ord.Id] = ord
	}
	return res
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.Id == "" {
		o.Id = utils.GenerateReadableID(16)
	}
	return nil
}

type OrderProduct struct {
	Id        string  `json:"id"`
	ProductId string  `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
	OrderId   string  `json:"order_id"`
	*common.AuditFields
}

func NewOrderProduct(productId string, quantity int64, price float64) *OrderProduct {
	return &OrderProduct{
		ProductId:   productId,
		Quantity:    quantity,
		Price:       price,
		AuditFields: common.NewAuditFields(),
	}
}

func NewOrderProductWithOrderId(productId string, quantity int64, price float64, orderId string) *OrderProduct {
	return &OrderProduct{
		ProductId:   productId,
		Quantity:    quantity,
		Price:       price,
		OrderId:     orderId,
		AuditFields: common.NewAuditFields(),
	}
}

type OrderProductBaseResponse struct {
	Id        string  `json:"id"`
	ProductId string  `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
}

func NewOrderProductBaseResponse(orderProduct *OrderProduct) *OrderProductBaseResponse {
	if orderProduct == nil {
		return nil
	}
	return &OrderProductBaseResponse{
		Id:        orderProduct.Id,
		ProductId: orderProduct.ProductId,
		Quantity:  orderProduct.Quantity,
		Price:     orderProduct.Price,
	}
}

type ListOrderProductResponse []*OrderProduct

func (l *ListOrderProductResponse) ToOrderMap() map[string][]*OrderProduct {
	res := map[string][]*OrderProduct{}
	for _, op := range *l {
		res[op.OrderId] = append(res[op.OrderId], op)
	}
	return res
}

func (l *ListOrderProductResponse) ListProductIds() []string {
	productIds := []string{}
	for _, prod := range *l {
		productIds = append(productIds, prod.ProductId)
	}
	return productIds
}

func (op *OrderProduct) BeforeCreate(tx *gorm.DB) error {
	if op.Id == "" {
		op.Id = utils.GenerateReadableID(16)
	}
	return nil
}
