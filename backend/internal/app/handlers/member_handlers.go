package handlers

import (
	"discord-backend/internal/app/models"
	"discord-backend/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberHandler struct {
	MemberService *services.MemberService
}

func NewMemberHandler(memberService *services.MemberService) *MemberHandler {
	return &MemberHandler{MemberService: memberService}
}

func (m *MemberHandler) UpdateMemberRole(c *gin.Context) {
	profileIDInterface, exists := c.Get("profile_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "profile_id not found"})
		return
	}

	profileIDString, ok := profileIDInterface.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID format"})
		return
	}

	profileID, err := uuid.Parse(profileIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
		return
	}

	paramServerID := c.Param("serverId")
	serverID, err := uuid.Parse(paramServerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Server UUID format"})
		return
	}

	paramMemberID := c.Param("memberId")
	memberID, err := uuid.Parse(paramMemberID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Member UUID format"})
		return
	}

	var payload struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var role models.MemberRole
	switch payload.Role {
	case string(models.Admin):
		role = models.Admin
	case string(models.Moderator):
		role = models.Moderator
	case string(models.Guest):
		role = models.Guest
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role value"})
		return
	}

	server, err := m.MemberService.UpdateMemberRole(serverID, profileID, memberID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member role updated successfully", "server": server})
}
