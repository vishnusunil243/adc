package models

import (
	"gorm.io/gorm"
	"main.go/common"
	"main.go/common/utils"
)

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

func (u *OAuth2Token) BeforeCreate(tx *gorm.DB) error {
	if u.Id == "" {
		u.Id = utils.GenerateReadableID(16)
	}
	return nil
}
