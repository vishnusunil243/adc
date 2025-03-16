package order_service

import "main.go/internal/models"

type OrderResponse struct {
	Id       string                        `json:"id"`
	Total    float64                       `json:"total"`
	Products []*models.ProductBaseResponse `json:"products"`
}

func NewOrderResponse(order *models.Order, productMap map[string]*models.Product, orderProducts []*models.OrderProduct) *OrderResponse {
	products := []*models.ProductBaseResponse{}
	for _, pr := range orderProducts {
		prod := productMap[pr.Id]
		if prod != nil {
			product := models.NewProductBaseResponse(prod)
			product.Price = pr.Price
			products = append(products, product)
		}
	}
	return &OrderResponse{
		Id:       order.Id,
		Total:    order.Total,
		Products: products,
	}
}

func NewOrderListResponse(orders []*models.Order, productMap map[string]*models.Product, orderProductMap map[string][]*models.OrderProduct) []*OrderResponse {
	res := []*OrderResponse{}
	for _, order := range orders {
		orderProducts := orderProductMap[order.Id]
		if len(orderProducts) > 0 {
			res = append(res, NewOrderResponse(order, productMap, orderProducts))
		}
	}
	return res
}
