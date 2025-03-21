package models

import (
	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/utils"
)

type Product struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Images []string `gorm:"type:text[]" json:"images"`
	Price  float64  `json:"price"`
	*common.AuditFields
}

func NewProduct(name string, images []string, price float64) *Product {
	return &Product{
		Name:        name,
		Images:      images,
		Price:       price,
		AuditFields: common.NewAuditFields(),
	}
}

type ProductBaseResponse struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Images []string `json:"images"`
	Price  float64  `json:"price"`
}

func NewProductBaseResponse(product *Product) *ProductBaseResponse {
	if product == nil {
		return nil
	}
	return &ProductBaseResponse{
		Id:     product.Id,
		Name:   product.Name,
		Images: product.Images,
		Price:  product.Price,
	}
}

func (p *Product) UpdateName(name *string) {
	if name != nil {
		p.Name = *name
	}
}

func (p *Product) UpdateImages(images *[]string) {
	if images != nil {
		p.Images = *images
	}
}

func (p *Product) UpdatePrice(price *float64) {
	if price != nil {
		p.Price = *price
	}
}

type ListProductResponse []*Product

func (l *ListProductResponse) GetUserIds() []string {
	userIds := []string{}
	for _, product := range *l {
		userIds = append(userIds, product.GetAuditFieldsUserIds()...)
	}
	return userIds
}

func (l *ListProductResponse) ToMap() map[string]*Product {
	res := map[string]*Product{}
	for _, prod := range *l {
		res[prod.Id] = prod
	}
	return res
}

func (u *Product) BeforeCreate(tx *gorm.DB) error {
	if u.Id == "" {
		u.Id = utils.GenerateReadableID(16)
	}
	return nil
}
