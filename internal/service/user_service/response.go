package user_service

import "main.go/internal/models"

type UserRes struct {
	Id       string          `json:"id"`
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Logo     string          `json:"logo"`
	UserType models.UserType `json:"user_type"`
}

type LoginRes struct {
	*UserRes
	Token string `json:"token"`
}

func NewUserRes(user *models.User) *UserRes {
	return &UserRes{
		Id:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Logo:     user.Logo,
		UserType: user.UserType,
	}
}

func NewLoginRes(user *models.User, token string) *LoginRes {
	return &LoginRes{
		UserRes: NewUserRes(user),
		Token:   token,
	}
}

type ListUsersRes []*UserRes

func NewListUsersRes(users []*models.User) ListUsersRes {
	var listUsersRes ListUsersRes
	for _, user := range users {
		listUsersRes = append(listUsersRes, NewUserRes(user))
	}
	return listUsersRes
}
