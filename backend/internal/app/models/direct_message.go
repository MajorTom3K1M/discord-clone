package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DirectMessage struct {
	ID             uuid.UUID    `gorm:"type:uuid;primary_key;" json:"id"`
	Content        string       `gorm:"type:text" json:"content"`
	FileURL        *string      `gorm:"type:text" json:"fileUrl"`
	MemberID       uuid.UUID    `json:"memberID"`
	Member         Member       `gorm:"foreignKey:MemberID;references:ID;constraint:OnDelete:CASCADE;" json:"member"`
	ConversationID uuid.UUID    `json:"conversationId"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE;" json:"conversation"`
	Deleted        bool         `gorm:"default:false" json:"deleted"`
	CreatedAt      time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (directMessage *DirectMessage) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID need to sortable for pagination so I decide to using uuidv7
	directMessage.ID, err = uuid.NewV7()

	if err != nil {
		return err
	}

	return
}

func (directMessage *DirectMessage) Validate() error {
	if directMessage.Content == "" {
		return fmt.Errorf("Content cannot be empty")
	}
	return nil
}
