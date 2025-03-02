package models

import "main.go/common"

type Product struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Images []string `json:"images"`
	Price  string   `json:"price"`
	*common.AuditFields
}

func NewProduct(name string, images []string, price string) *Product {
	return &Product{
		Name:        name,
		Images:      images,
		Price:       price,
		AuditFields: common.NewAuditFields(),
	}
}
