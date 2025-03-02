package models

import "main.go/common"

type OAuth2Token struct {
	Id     string `json:"id"`
	Token  string `json:"token"`
	UserId string `json:"user_id"`
	*common.AuditFields
}

func NewOauth2Token(token, userId string) *OAuth2Token {
	return &OAuth2Token{
		Token:       token,
		UserId:      userId,
		AuditFields: common.NewAuditFields(),
	}
}
