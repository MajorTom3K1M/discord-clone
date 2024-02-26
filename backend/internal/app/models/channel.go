package models

import (
	"encoding/json"
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
	ID        uuid.UUID   `gorm:"type:uuid;primary_key;" json:"id"`
	Name      string      `json:"name"`
	Type      ChannelType `gorm:"type:varchar(100);default:'TEXT'" json:"type"`
	ProfileID uuid.UUID   `json:"profileID"`
	Profile   Profile     `gorm:"foreignKey:ProfileID;references:ID;onDelete:CASCADE" json:"profile"`
	ServerID  uuid.UUID   `json:"serverID"`
	Server    Server      `gorm:"foreignKey:ServerID;references:ID;onDelete:CASCADE" json:"server"`
	Messages  []Message   `json:"messages"`
	CreatedAt time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (channel *Channel) BeforeCreate(tx *gorm.DB) (err error) {
	channel.ID = uuid.New()
	return
}

func (channel *Channel) MarshalJSON() ([]byte, error) {
	type Alias Channel
	alias := (*Alias)(channel)

	temp := struct {
		*Alias
		Server  *Server  `json:"server,omitempty"`
		Profile *Profile `json:"profile,omitempty"`
	}{
		Alias: alias,
	}

	if channel.Server.ID != uuid.Nil {
		temp.Server = &channel.Server
	}

	if channel.Profile.ID != uuid.Nil {
		temp.Profile = &channel.Profile
	}

	return json.Marshal(temp)
}
