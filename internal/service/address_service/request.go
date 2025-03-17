package address_service

type GetAddressReq struct {
	Id string `json:"id"`
}

type ListAddressReq struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	UserId string `json:"user_id"`
}

type CreateAddressReq struct {
	UserId  string `json:"user_id"`
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	Pincode string `json:"pincode"`
	Street  string `json:"street"`
	Area    string `json:"area"`
}

type UpdateAddressReq struct {
	Id      string `json:"id"`
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	Pincode string `json:"pincode"`
	Street  string `json:"street"`
	Area    string `json:"area"`
}

type DeleteAddressReq struct {
	Ids []string `json:"ids"`
}
