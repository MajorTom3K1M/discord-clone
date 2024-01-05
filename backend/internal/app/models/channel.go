package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChannelType string

const (
	Text  ChannelType = "TEXT"
	Audio ChannelType = "AUDIO"
	Video ChannelType = "VIDEO"
)

type Channel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name      string
	Type      ChannelType `gorm:"type:varchar(100);default:'TEXT'"`
	ProfileID uuid.UUID
	Profile   Profile `gorm:"foreignKey:ProfileID;references:ID;onDelete:CASCADE"`
	ServerID  uuid.UUID
	Server    Server `gorm:"foreignKey:ServerID;references:ID;onDelete:CASCADE"`
	Messages  []Message
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time
}

func (channel *Channel) BeforeCreate(tx *gorm.DB) (err error) {
	channel.ID = uuid.New()
	return
}
