package models

import (
	"encoding/json"
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
	ID                     uuid.UUID       `gorm:"type:uuid;primary_key;" json:"id"`
	Role                   MemberRole      `gorm:"type:varchar(100);default:'GUEST'" json:"role"`
	ProfileID              uuid.UUID       `json:"profileID"`
	Profile                Profile         `gorm:"foreignKey:ProfileID;references:ID;onDelete:CASCADE" json:"profile"`
	ServerID               uuid.UUID       `json:"serverID"`
	Server                 Server          `gorm:"foreignKey:ServerID;references:ID;onDelete:CASCADE" json:"server"`
	Messages               []Message       `json:"messages"`
	DirectMessages         []DirectMessage `json:"directMessages"`
	ConversationsInitiated []Conversation  `gorm:"foreignKey:MemberOneID" json:"conversationsInitiated"`
	ConversationsReceived  []Conversation  `gorm:"foreignKey:MemberTwoID" json:"conversationsReceived"`
	CreatedAt              time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
}

func (member *Member) BeforeCreate(tx *gorm.DB) (err error) {
	member.ID = uuid.New()
	return
}

func (member *Member) MarshalJSON() ([]byte, error) {
	type Alias Member
	alias := (*Alias)(member)

	temp := struct {
		*Alias
		Server  *Server  `json:"server,omitempty"`
		Profile *Profile `json:"profile,omitempty"`
	}{
		Alias: alias,
	}

	if member.Server.ID != uuid.Nil {
		temp.Server = &member.Server
	}

	if member.Profile.ID != uuid.Nil {
		temp.Profile = &member.Profile
	}

	return json.Marshal(temp)
}
