package order_service

import "main.go/internal/models"

type AddOrderRequest struct {
	Products []*OrderProduct `json:"products"`
}

func (a *AddOrderRequest) GetProductIds() []string {
	productIds := []string{}
	for _, prod := range a.Products {
		productIds = append(productIds, prod.Id)
	}
	return productIds
}

type OrderProduct struct {
	Id       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

type ListOrderRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type GetOrderRequest struct {
	Id string `json:"id"`
}

type UpdateOrderRequest struct {
	Id     string        `json:"id"`
	Status models.Status `json:"status"`
}
