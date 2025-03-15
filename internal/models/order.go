package models

import "main.go/common"

type Order struct {
	Id     string  `json:"id"`
	UserId string  `json:"user_id"`
	Total  float64 `json:"total"`
	*common.AuditFields
}

type OrderProduct struct {
	Id        string  `json:"id"`
	ProductId string  `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
	*common.AuditFields
}
