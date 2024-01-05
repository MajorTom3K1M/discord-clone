package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DirectMessage struct {
	ID             uuid.UUID    `gorm:"type:uuid;primary_key;"`
	Content        string       `gorm:"type:text"`
	FileURL        *string      `gorm:"type:text"`
	MemberID       uuid.UUID    // Foreign key for Member
	Member         Member       `gorm:"foreignKey:MemberID;constraint:OnDelete:CASCADE;"`
	ConversationID uuid.UUID    // Foreign key for Conversation
	Conversation   Conversation `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE;"`
	Deleted        bool         `gorm:"default:false"`
	CreatedAt      time.Time    `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time
}

func (directMessage *DirectMessage) BeforeCreate(tx *gorm.DB) (err error) {
	directMessage.ID = uuid.New()
	return
}
