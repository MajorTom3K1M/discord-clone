package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key;" json:"id"`
	MemberOneID    uuid.UUID       `gorm:"index:,unique" json:"memberOneID"`
	MemberOne      Member          `gorm:"foreignKey:MemberOneID;references:ID;onDelete:CASCADE" json:"memberOne"`
	MemberTwoID    uuid.UUID       `gorm:"index:,unique;index" json:"memberTwoID"`
	MemberTwo      Member          `gorm:"foreignKey:MemberTwoID;references:ID;onDelete:CASCADE" json:"memberTwo"`
	DirectMessages []DirectMessage `json:"directMessages"`
	CreatedAt      time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func (conversation *Conversation) BeforeCreate(tx *gorm.DB) (err error) {
	conversation.ID = uuid.New()
	return
}
