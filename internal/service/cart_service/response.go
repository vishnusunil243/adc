package cart_service

import (
	"main.go/internal/models"
	"main.go/internal/service"
)

// **Response Structs** (Ensures separation from database models)
type CartRes struct {
	Id        string                      `json:"id"`
	ProductId string                      `json:"product_id"`
	Quantity  int64                       `json:"quantity"`
	Product   *models.ProductBaseResponse `json:"product"`
	*service.AuditFieldResponse
}

// NewCartRes converts a Cart model into a response struct
func NewCartRes(cart *models.Cart, userData models.ListResponse, product *models.Product) *CartRes {
	userMap := userData.ToMap()
	return &CartRes{
		Id:                 cart.Id,
		ProductId:          cart.ProductId,
		Quantity:           cart.Quantity,
		Product:            models.NewProductBaseResponse(product),
		AuditFieldResponse: service.NewAuditFieldResponse(cart.AuditFields, userMap),
	}
}

// **List Response Wrapper**
type ListCartRes []*CartRes

func NewListCartRes(carts []*models.Cart, userData models.ListResponse, productList models.ListProductResponse) []*CartRes {
	res := []*CartRes{}
	prodMap := productList.ToMap()
	for _, cart := range carts {
		product := prodMap[cart.ProductId]
		res = append(res, NewCartRes(cart, userData, product))
	}
	return res
}
