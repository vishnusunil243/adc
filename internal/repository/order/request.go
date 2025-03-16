package order

type ListOrderRequest struct {
	Limit  int
	Offset int
	UserId string
	Ids    []string
}

type GetOrderRequest struct {
	Id string
}

type GetOrderProductRequest struct {
	Id string
}

type ListOrderProductRequest struct {
	Ids      []string `json:"ids"`
	Limit    int
	Offset   int
	OrderId  string
	OrderIds []string
}
