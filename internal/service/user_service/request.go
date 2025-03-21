package user_service

import "main.go/internal/models"

type SignupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ListUserReq struct {
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Ids    []string `json:"ids"`
}

type GetUserReq struct {
	Email    string          `json:"email"`
	Id       string          `json:"id"`
	UserType models.UserType `json:"user_type"`
}
