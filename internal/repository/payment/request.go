package payment

type ListPaymentRequest struct {
	Ids    []string
	Limit  int
	Offset int
}

type GetPaymentRequest struct {
	Id      string `json:"string"`
	OrderId string `json:"order_id"`
}
