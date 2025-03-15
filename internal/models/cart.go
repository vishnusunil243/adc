package models

import "main.go/common"

type Cart struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
	UserId    string `json:"user_id"`
	*common.AuditFields
}

type ListCart []*Cart

func (l *ListCart) GetUserIds() []string {
	userIds := []string{}
	for _, cart := range *l {
		userIds = append(userIds, cart.GetAuditFieldsUserIds()...)
	}
	return userIds
}

func (l *ListCart) GetProductIds() []string {
	productIds := []string{}
	for _, cart := range *l {
		productIds = append(productIds, cart.ProductId)
	}
	return productIds
}

func NewCart(productId string, quantity int64, createdby string) *Cart {
	return &Cart{
		ProductId:   productId,
		Quantity:    quantity,
		UserId:      createdby,
		AuditFields: common.NewAuditFieldsWithCreatedBy(createdby),
	}
}

func (c *Cart) UpdateQuantity(qty *int64) {
	if qty != nil {
		c.Quantity = *qty
	}
}
