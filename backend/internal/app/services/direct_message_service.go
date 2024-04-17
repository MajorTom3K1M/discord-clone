package services

import (
	"discord-backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const DIRECT_MESSAGES_BATCH = 10

type DirectMessageService struct {
	DB *gorm.DB
}

func NewDirectMessageService(db *gorm.DB) *DirectMessageService {
	return &DirectMessageService{DB: db}
}

func (s *DirectMessageService) CreateDirectMessage(conversationID, memberID uuid.UUID, content, fileUrl string) (*models.DirectMessage, error) {
	directMessage := models.DirectMessage{
		Content:        content,
		FileURL:        &fileUrl,
		ConversationID: conversationID,
		MemberID:       memberID,
	}

	if err := s.DB.Create(&directMessage).Error; err != nil {
		return nil, err
	}

	var reponseMessage models.DirectMessage
	if err := s.DB.Preload("Member.Profile").Where("id = ?", directMessage.ID).
		First(&reponseMessage).Error; err != nil {
		return nil, err
	}

	return &reponseMessage, nil
}

func (s *DirectMessageService) GetDirectMessages(conversationID uuid.UUID, cursor string) ([]models.DirectMessage, string, error) {
	var directMessages []models.DirectMessage

	query := s.DB.Preload("Member.Profile").Where("conversation_id = ?", conversationID).
		Order("created_at DESC").Limit(DIRECT_MESSAGES_BATCH)

	if cursor != "" {
		cursorUUID, err := uuid.Parse(cursor)
		if err != nil {
			return nil, "", err
		}
		query = query.Where("id < ?", cursorUUID)
	}

	if err := query.Find(&directMessages).Error; err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(directMessages) == DIRECT_MESSAGES_BATCH {
		nextCursor = directMessages[DIRECT_MESSAGES_BATCH-1].ID.String()
	}

	return directMessages, nextCursor, nil
}

func (s *DirectMessageService) GetDirectMessage(conversationID, directMessageID uuid.UUID) (*models.DirectMessage, error) {
	var directMessage models.DirectMessage
	if err := s.DB.Preload("Member.Profile").Where("id = ? AND conversation_id = ? AND deleted = false", directMessageID, conversationID).
		First(&directMessage).Error; err != nil {
		return nil, err
	}

	return &directMessage, nil
}

func (s *DirectMessageService) UpdateDirectMessage(conversationID, directMessageID uuid.UUID, content string) (*models.DirectMessage, error) {
	if err := s.DB.Model(&models.DirectMessage{}).
		Where("id = ? AND conversation_id = ?", directMessageID, conversationID).
		Update("content", content).Error; err != nil {
		return nil, err
	}

	var directMessage models.DirectMessage
	if err := s.DB.Preload("Member.Profile").First(&directMessage, directMessageID).Error; err != nil {
		return nil, err
	}

	return &directMessage, nil
}

func (s *DirectMessageService) DeleteDirectMessage(conversationID, directMessageID uuid.UUID) (*models.DirectMessage, error) {
	if err := s.DB.Model(&models.DirectMessage{}).
		Where("id = ? AND conversation_id = ?", directMessageID, conversationID).
		Updates(models.DirectMessage{FileURL: nil, Content: "This message has been deleted.", Deleted: true}).
		Error; err != nil {
		return nil, err
	}

	var directMessage models.DirectMessage
	if err := s.DB.Preload("Member.Profile").First(&directMessage, directMessageID).Error; err != nil {
		return nil, err
	}

	return &directMessage, nil
}
