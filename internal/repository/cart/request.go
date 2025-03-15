package cart

type ListCartRequest struct {
	Limit  int
	Offset int
	Ids    []string
	UserId string
}

type GetCartRequest struct {
	Id        string
	UserId    string
	ProductId string
}
