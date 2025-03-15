package product_service

import (
	"main.go/internal/models"
	"main.go/internal/service"
)

type ProductRes struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Images []string `json:"images"`
	Price  float64  `json:"price"`
	*service.AuditFieldResponse
}

func NewProductRes(product *models.Product, userData models.ListResponse) *ProductRes {
	userMap := userData.ToMap()
	return &ProductRes{
		Id:                 product.Id,
		Name:               product.Name,
		Images:             product.Images,
		Price:              product.Price,
		AuditFieldResponse: service.NewAuditFieldResponse(product.AuditFields, userMap),
	}
}

type ListProductRes []*ProductRes

func NewListProductRes(products []*models.Product, userData models.ListResponse) []*ProductRes {
	res := []*ProductRes{}
	for _, product := range products {
		res = append(res, NewProductRes(product, userData))
	}
	return res
}
