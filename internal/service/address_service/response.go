package address_service

import (
	"main.go/internal/models"
	"main.go/internal/service"
)

type AddressRes struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	Pincode string `json:"pincode"`
	Street  string `json:"street"`
	Area    string `json:"area"`
	*service.AuditFieldResponse
}

func NewAddressRes(address *models.Address, userData models.ListResponse) *AddressRes {
	userMap := userData.ToMap()
	return &AddressRes{
		Id:                 address.Id,
		UserId:             address.UserId,
		Country:            address.Country,
		State:              address.State,
		City:               address.City,
		Pincode:            address.Pincode,
		Street:             address.Street,
		Area:               address.Area,
		AuditFieldResponse: service.NewAuditFieldResponse(address.AuditFields, userMap),
	}
}

type ListAddressRes []*AddressRes

func NewListAddressRes(addresses []*models.Address, userData models.ListResponse) []*AddressRes {
	res := []*AddressRes{}
	for _, address := range addresses {
		res = append(res, NewAddressRes(address, userData))
	}
	return res
}
