package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Server struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name       string
	ImageURL   string `gorm:"type:text"`
	InviteCode string `gorm:"unique;"`
	ProfileID  uuid.UUID
	Profile    Profile   `gorm:"foreignKey:ProfileID;references:ID;onDelete:CASCADE"`
	Members    []Member  `gorm:"foreignKey:ServerID"`
	Channels   []Channel `gorm:"foreignKey:ServerID"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time
}

func (server *Server) BeforeCreate(tx *gorm.DB) (err error) {
	server.ID = uuid.New()
	return
}
