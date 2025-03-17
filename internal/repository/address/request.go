package address

type GetAddressRequest struct {
	Id string `json:"id"`
}

type ListAddressRequest struct {
	Limit  int
	Offset int
	Ids    []string
	UserId string
}
