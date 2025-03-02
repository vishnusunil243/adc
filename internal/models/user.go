package models

import "main.go/common"

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
