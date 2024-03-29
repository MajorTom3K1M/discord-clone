package handlers

import (
	"discord-backend/internal/app/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationHandler struct {
	ConversationService *services.ConversationService
}

func NewConversationHandler(conversationService *services.ConversationService) *ConversationHandler {
	return &ConversationHandler{ConversationService: conversationService}
}

func (s *ConversationHandler) GetOrCreateConversation(c *gin.Context) {
	var req struct {
		MemberOneID uuid.UUID `json:"memberOneId"`
		MemberTwoID uuid.UUID `json:"memberTwoId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conversation, err := s.ConversationService.FindConversation(req.MemberOneID, req.MemberTwoID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Get conversation successfully", "conversation": conversation})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		conversation, err := s.ConversationService.CreateConversation(req.MemberOneID, req.MemberTwoID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Create conversation successfully", "conversation": conversation})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func (s *ConversationHandler) FindConversation(c *gin.Context) {
	memberOneIdStr := c.Param("memberOneId")
	memberTwoIdStr := c.Param("memberTwoId")

	memberOneId, err := uuid.Parse(memberOneIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid memberOneId"})
		return
	}

	memberTwoId, err := uuid.Parse(memberTwoIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid memberTwoId"})
		return
	}

	conversation, err := s.ConversationService.FindConversation(memberOneId, memberTwoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get conversation successfully", "conversation": conversation})
}
