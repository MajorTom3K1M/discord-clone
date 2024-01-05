package db

import (
	"discord-backend/internal/app/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Profile{},
		&models.Server{},
		&models.Channel{},
		&models.Member{},
		&models.Message{},
		&models.DirectMessage{},
		&models.Conversation{},
		&models.RefreshToken{},
	)
}
