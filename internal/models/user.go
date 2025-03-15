package models

import (
	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/utils"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Logo     string `json:"logo"`
	*common.AuditFields
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:        name,
		Email:       email,
		Password:    password,
		AuditFields: common.NewAuditFields(),
	}
}

type ListResponse []*User

func (l *ListResponse) ToMap() map[string]*User {
	res := map[string]*User{}
	for _, user := range *l {
		res[user.Id] = user
	}
	return res
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Id == "" {
		u.Id = utils.GenerateReadableID(16)
	}
	return nil
}
