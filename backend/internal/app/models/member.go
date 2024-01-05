package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberRole string

const (
	Admin     MemberRole = "ADMIN"
	Moderator MemberRole = "MODERATOR"
	Guest     MemberRole = "GUEST"
)

type Member struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primary_key;"`
	Role                   MemberRole `gorm:"type:varchar(100);default:'GUEST'"`
	ProfileID              uuid.UUID
	Profile                Profile `gorm:"foreignKey:ProfileID;references:ID;onDelete:CASCADE"`
	ServerID               uuid.UUID
	Server                 Server `gorm:"foreignKey:ServerID;references:ID;onDelete:CASCADE"`
	Messages               []Message
	DirectMessages         []DirectMessage
	ConversationsInitiated []Conversation `gorm:"foreignKey:MemberOneID"`
	ConversationsReceived  []Conversation `gorm:"foreignKey:MemberTwoID"`
	CreatedAt              time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt              time.Time
}

func (member *Member) BeforeCreate(tx *gorm.DB) (err error) {
	member.ID = uuid.New()
	return
}
