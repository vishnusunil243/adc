package product

type ListProductRequest struct {
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Ids    []string `json:"ids"`
}

type GetProductRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
