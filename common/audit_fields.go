package common

import "time"

type AuditFields struct {
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt int64  `json:"deleted_at"`
	IsDeleted int    `json:"is_deleted" gorm:"default:0"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

func NewAuditFields() *AuditFields {
	return &AuditFields{
		CreatedAt: time.Now().UTC().UnixMilli(),
		UpdatedAt: time.Now().UTC().UnixMilli(),
		IsDeleted: 0,
	}
}

func NewAuditFieldsWithCreatedBy(userId string) *AuditFields {
	return &AuditFields{
		CreatedAt: time.Now().UTC().UnixMilli(),
		UpdatedAt: time.Now().UTC().UnixMilli(),
		CreatedBy: userId,
		UpdatedBy: userId,
		IsDeleted: 0,
	}
}

func (a *AuditFields) GetAuditFieldsUserIds() []string {
	userIds := []string{}
	userIds = append(userIds, a.CreatedBy, a.UpdatedBy)
	return userIds
}

func GetFieldsForDelete() map[string]interface{} {
	return map[string]interface{}{
		"is_deleted": 1,
		"deleted_at": time.Now().UTC().UnixMilli(),
	}
}
