package service

import (
	"main.go/common"
	"main.go/internal/models"
)

type AuditFieldRes struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
}

func NewAuditFieldRes(user *models.User) *AuditFieldRes {
	if user == nil {
		return nil
	}
	return &AuditFieldRes{
		UserId: user.Id,
		Name:   user.Name,
		Logo:   user.Logo,
	}
}

type AuditFieldResponse struct {
	CreatedBy *AuditFieldRes `json:"created_by"`
	UpdatedBy *AuditFieldRes `json:"updated_by"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
}

func NewAuditFieldResponse(auditField *common.AuditFields, userdata map[string]*models.User) *AuditFieldResponse {
	return &AuditFieldResponse{
		CreatedBy: NewAuditFieldRes(userdata[auditField.CreatedBy]),
		UpdatedBy: NewAuditFieldRes(userdata[auditField.UpdatedBy]),
		CreatedAt: auditField.CreatedAt,
		UpdatedAt: auditField.UpdatedAt,
	}
}
