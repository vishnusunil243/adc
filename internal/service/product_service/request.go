package product_service

type CreateProductReq struct {
	Name   string   `json:"name"`
	Images []string `json:"images"`
	Price  float64  `json:"price"`
}

type UpdateProductReq struct {
	Id     string    `json:"id"`
	Name   *string   `json:"name"`
	Images *[]string `json:"images"`
	Price  *float64  `json:"price"`
}

type DeleteProductReq struct {
	Ids []string `json:"ids"`
}

type GetProductReq struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ListProductReq struct {
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Ids    []string `json:"ids"`
}
