package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Server struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name       string    `json:"name"`
	ImageURL   string    `gorm:"type:text" json:"imageUrl"`
	InviteCode string    `gorm:"unique;" json:"inviteCode"`
	ProfileID  uuid.UUID `json:"profileID"`
	Profile    Profile   `gorm:"foreignKey:ProfileID;references:ID;onDelete:CASCADE" json:"profile,omitempty"`
	Members    []Member  `gorm:"foreignKey:ServerID;onDelete:CASCADE" json:"members,omitempty"`
	Channels   []Channel `gorm:"foreignKey:ServerID;onDelete:CASCADE" json:"channels,omitempty"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (server *Server) BeforeCreate(tx *gorm.DB) (err error) {
	server.ID = uuid.New()
	return
}

// type ServerResponse struct {
// 	ID         uuid.UUID         `json:"id"`
// 	Name       string            `json:"name"`
// 	ImageURL   string            `json:"imageUrl"`
// 	InviteCode string            `json:"inviteCode"`
// 	ProfileID  uuid.UUID         `json:"profileID"`
// 	Profile    Profile           `json:"profile"`
// 	Members    []MemberResponse  `json:"members"`
// 	Channels   []ChannelResponse `json:"channels"`
// 	CreatedAt  time.Time         `json:"created_at"`
// 	UpdatedAt  time.Time         `json:"updated_at"`
// }

func (server *Server) MarshalJSON() ([]byte, error) {
	type Alias Server
	alias := (*Alias)(server)

	temp := struct {
		*Alias
		Profile *Profile `json:"profile,omitempty"`
	}{
		Alias: alias,
	}

	if server.Profile.ID != uuid.Nil {
		temp.Profile = &server.Profile
	}

	return json.Marshal(temp)
}
