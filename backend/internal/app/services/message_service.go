package services

import (
	"discord-backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const MESSAGES_BATCH = 10

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

	var reponseMessage models.Message
	if err := s.DB.Preload("Member.Profile").Where("id = ?", message.ID).
		First(&reponseMessage).Error; err != nil {
		return nil, err
	}

	return &reponseMessage, nil
}

func (s *MessageService) GetMessages(channelID uuid.UUID, cursor string) ([]models.Message, string, error) {
	var messages []models.Message

	query := s.DB.Preload("Member.Profile").Where("channel_id = ?", channelID).
		Order("created_at DESC").Limit(MESSAGES_BATCH)

	if cursor != "" {
		cursorUUID, err := uuid.Parse(cursor)
		if err != nil {
			return nil, "", err
		}
		query = query.Where("id < ?", cursorUUID)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(messages) == MESSAGES_BATCH {
		nextCursor = messages[MESSAGES_BATCH-1].ID.String()
	}

	return messages, nextCursor, nil
}

func (s *MessageService) GetMessage(channelID, messageID uuid.UUID) (*models.Message, error) {
	var message models.Message
	if err := s.DB.Preload("Member.Profile").Where("id = ? AND channel_id = ? AND deleted = false", messageID, channelID).
		First(&message).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *MessageService) UpdateMessage(channelID, messageID uuid.UUID, content string) (*models.Message, error) {
	if err := s.DB.Model(&models.Message{}).Where("id = ? AND channel_id = ?", messageID, channelID).
		Update("content", content).Error; err != nil {
		return nil, err
	}

	var message models.Message
	if err := s.DB.Preload("Member.Profile").First(&message, messageID).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *MessageService) DeleteMessage(channelID, messageID uuid.UUID) (*models.Message, error) {
	if err := s.DB.Model(&models.Message{}).
		Where("id = ? AND channel_id = ?", messageID, channelID).
		Updates(models.Message{FileURL: nil, Content: "This message has been deleted.", Deleted: true}).
		Error; err != nil {
		return nil, err
	}

	var message models.Message
	if err := s.DB.Preload("Member.Profile").First(&message, messageID).Error; err != nil {
		return nil, err
	}

	return &message, nil
}
