package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name      string    `json:"name"`
	ImageURL  string    `gorm:"type:text" json:"imageUrl"`
	Email     string    `gorm:"type:text" json:"email"`
	Password  string    `json:"-"`
	Servers   []Server  `gorm:"foreignKey:ProfileID" json:"servers"`
	Members   []Member  `gorm:"foreignKey:ProfileID" json:"members"`
	Channels  []Channel `gorm:"foreignKey:ProfileID" json:"channels"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
