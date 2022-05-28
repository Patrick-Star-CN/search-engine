package database

import (
	"gorm.io/gorm"
	"search-engine/app/models"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.DocRaw{},
		&models.DocID{},
		&models.WordMap{},
		&models.Collection{},
		&models.User{})
}
