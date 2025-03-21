package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(&Address{}, &User{}, &Cart{}, &Product{}, &Order{}, &OrderProduct{}, &OAuth2Token{})
	return nil
}
