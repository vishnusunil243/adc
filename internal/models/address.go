package models

import (
	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/utils"
)

type Address struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	Pincode string `json:"pincode"`
	Street  string `json:"street"`
	Area    string `json:"area"`
	*common.AuditFields
}

func (u *Address) BeforeCreate(tx *gorm.DB) error {
	if u.Id == "" {
		u.Id = utils.GenerateReadableID(16)
	}
	return nil
}

type ListAddress []*Address

// Function to get a list of user IDs from a list of addresses
func (l *ListAddress) GetUserIds() []string {
	userIds := []string{}
	for _, address := range *l {
		userIds = append(userIds, address.GetAuditFieldsUserIds()...)
	}
	return userIds
}

// Function to get a list of all address IDs from a list of addresses
func (l *ListAddress) GetAddressIds() []string {
	addressIds := []string{}
	for _, address := range *l {
		addressIds = append(addressIds, address.Id)
	}
	return addressIds
}

// Constructor to create a new address
func NewAddress(userId, country, state, city, pincode, street, area, createdby string) *Address {
	return &Address{
		UserId:      userId,
		Country:     country,
		State:       state,
		City:        city,
		Pincode:     pincode,
		Street:      street,
		Area:        area,
		AuditFields: common.NewAuditFieldsWithCreatedBy(createdby),
	}
}

// Function to update an address' details
func (a *Address) UpdateDetails(country, state, city, pincode, street, area string) {
	if country != "" {
		a.Country = country
	}
	if state != "" {
		a.State = state
	}
	if city != "" {
		a.City = city
	}
	if pincode != "" {
		a.Pincode = pincode
	}
	if street != "" {
		a.Street = street
	}
	if area != "" {
		a.Area = area
	}
}
