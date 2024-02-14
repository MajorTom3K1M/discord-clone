package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name      string
	ImageURL  string `gorm:"type:text"`
	Email     string `gorm:"type:text"`
	Password  string
	Servers   []Server  `gorm:"foreignKey:ProfileID"`
	Members   []Member  `gorm:"foreignKey:ProfileID"`
	Channels  []Channel `gorm:"foreignKey:ProfileID"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time
}

func (profile *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	profile.ID = uuid.New()
	return
}

type ProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"imageUrl"`
	Email     string    `json:"email"`
	Servers   []Server  `json:"servers"`
	Members   []Member  `json:"members"`
	Channels  []Channel `json:"channels"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
