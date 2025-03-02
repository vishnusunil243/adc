package common

import "time"

type AuditFields struct {
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt int64 `json:"deleted_at"`
	IsDeleted int   `json:"is_deleted" gorm:"default:0"`
}

func NewAuditFields() *AuditFields {
	return &AuditFields{
		CreatedAt: time.Now().UTC().UnixMilli(),
		UpdatedAt: time.Now().UTC().UnixMilli(),
		IsDeleted: 0,
	}
}

func GetFieldsForDelete() map[string]interface{} {
	return map[string]interface{}{
		"is_deleted": 0,
		"deleted_at": time.Now().UTC().UnixMilli(),
	}
}
