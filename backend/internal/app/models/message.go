package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Content   string    `gorm:"type:text"`
	FileURL   *string   `gorm:"type:text"`
	MemberID  uuid.UUID
	Member    Member `gorm:"foreignKey:MemberID;references:ID;onDelete:CASCADE"`
	ChannelID uuid.UUID
	Channel   Channel   `gorm:"foreignKey:ChannelID;references:ID;onDelete:CASCADE"`
	Deleted   bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time
}

func (message *Message) BeforeCreate(tx *gorm.DB) (err error) {
	message.ID = uuid.New()
	return
}
