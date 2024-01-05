package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	ProfileID uuid.UUID `gorm:"unique;"`
	Profile   Profile   `gorm:"foreignKey:ProfileID;references:ID;onDelete:CASCADE"`
	Token     string    `gorm:"size:500"`
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (refreshToken *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	refreshToken.ID = uuid.New()
	return
}
