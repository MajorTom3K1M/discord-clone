package handlers

import (
	"discord-backend/internal/app/models"
	"discord-backend/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChannelHandler struct {
	ChannelService *services.ChannelService
}

func NewChannelHandler(channelService *services.ChannelService) *ChannelHandler {
	return &ChannelHandler{ChannelService: channelService}
}

func (h *ChannelHandler) CreateChannel(c *gin.Context) {
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

	var channelData struct {
		Name        string             `json:"name"`
		ChannelType models.ChannelType `json:"type"`
	}
	if err := c.ShouldBindJSON(&channelData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if channelData.Name == "general" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name cannot be 'general'"})
		return
	}

	var channelType models.ChannelType
	switch channelData.ChannelType {
	case models.Text:
		channelType = models.Text
	case models.Audio:
		channelType = models.Audio
	case models.Video:
		channelType = models.Video
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channelType value"})
		return
	}

	server, err := h.ChannelService.CreateChannel(serverID, profileID, channelData.Name, channelType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Channel created successfully", "server": server})
}

func (h *ChannelHandler) DeleteChannel(c *gin.Context) {
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

	paramChannelID := c.Param("channelId")
	channelID, err := uuid.Parse(paramChannelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Channel UUID format"})
		return
	}

	server, err := h.ChannelService.DeleteChannel(serverID, profileID, channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Channel delete successfully", "server": server})
}
