package services

import (
	"discord-backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageService struct {
	DB *gorm.DB
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{DB: db}
}

func (s *MessageService) CreateMessage(channelID, memberID uuid.UUID, content, fileUrl string) (*models.Message, error) {
	message := models.Message{
		Content:   content,
		FileURL:   &fileUrl,
		ChannelID: channelID,
		MemberID:  memberID,
	}

	if err := s.DB.Create(&message).Error; err != nil {
		return nil, err
	}

	return &message, nil
}
