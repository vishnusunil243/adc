package cart_service

// **Request Structs** (These define inputs for service functions)
type AddToCartRequest struct {
	ProductId string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
	UserId    string `json:"user_id"`
}

type UpdateCartRequest struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Quantity  *int64 `json:"quantity"`
}

type DeleteCartRequest struct {
	Ids []string `json:"ids"`
}

type ListCartRequest struct {
	Ids    []string `json:"ids,omitempty"`
	UserId string   `json:"user_id,omitempty"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

type GetCartRequest struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
}
