package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;"`
	MemberOneID    uuid.UUID
	MemberOne      Member `gorm:"foreignKey:MemberOneID;references:ID;onDelete:CASCADE"`
	MemberTwoID    uuid.UUID
	MemberTwo      Member `gorm:"foreignKey:MemberTwoID;references:ID;onDelete:CASCADE"`
	DirectMessages []DirectMessage
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time
}

func (conversation *Conversation) BeforeCreate(tx *gorm.DB) (err error) {
	conversation.ID = uuid.New()
	return
}
