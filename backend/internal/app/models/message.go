package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Content   string    `gorm:"type:text" json:"content"`
	FileURL   *string   `gorm:"type:text" json:"fileUrl"`
	MemberID  uuid.UUID `json:"memberID"`
	Member    Member    `gorm:"foreignKey:MemberID;references:ID;onDelete:CASCADE" json:"member"`
	ChannelID uuid.UUID `json:"channelID"`
	Channel   Channel   `gorm:"foreignKey:ChannelID;references:ID;onDelete:CASCADE" json:"channel"`
	Deleted   bool      `gorm:"default:false" json:"deleted"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (message *Message) BeforeCreate(tx *gorm.DB) (err error) {
	message.ID = uuid.New()
	return
}
