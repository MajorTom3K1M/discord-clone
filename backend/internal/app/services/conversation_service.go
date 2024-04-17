package services

import (
	"discord-backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationService struct {
	DB *gorm.DB
}

func NewConversationService(db *gorm.DB) *ConversationService {
	return &ConversationService{DB: db}
}

func (c *ConversationService) GetConversation(conversationID, profileID uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	err := c.DB.Preload("MemberOne.Profile").Preload("MemberTwo.Profile").
		Where("id = ? AND (member_one_id IN (SELECT id FROM members WHERE profile_id = ?) OR member_two_id IN (SELECT id FROM members WHERE profile_id = ?))",
			conversationID, profileID, profileID).
		First(&conversation).Error

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (c *ConversationService) FindConversation(memberOneID, memberTwoID uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	err := c.DB.Preload("MemberOne.Profile").Preload("MemberTwo.Profile").
		Where("member_one_id = ? AND member_two_id = ?", memberOneID, memberTwoID).
		Or("member_one_id = ? AND member_two_id = ?", memberTwoID, memberOneID).
		First(&conversation).Error

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (c *ConversationService) CreateConversation(memberOneID, memberTwoID uuid.UUID) (*models.Conversation, error) {
	conversation := models.Conversation{
		MemberOneID: memberOneID,
		MemberTwoID: memberTwoID,
	}

	if err := c.DB.Create(&conversation).Error; err != nil {
		return nil, err
	}

	var fetchedConversation models.Conversation
	if err := c.DB.Preload("MemberOne.Profile").Preload("MemberTwo.Profile").
		First(&fetchedConversation, conversation.ID).Error; err != nil {
		return nil, err
	}

	return &fetchedConversation, nil
}
