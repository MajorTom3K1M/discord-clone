package handlers

import (
	"discord-backend/internal/app/models"
	"discord-backend/internal/app/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	ProfileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{ProfileService: profileService}
}

func (p *ProfileHandler) GetProfile(c *gin.Context) {
	paramProfileID := c.Param("id")
	profileID, err := uuid.Parse(paramProfileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	profile, err := p.ProfileService.GetProfileByID(profileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func (p *ProfileHandler) UpdateProfile(c *gin.Context) {
	var updatedData models.Profile
	paramProfileID := c.Param("id")
	profileID, err := uuid.Parse(paramProfileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := p.ProfileService.UpdateProfile(profileID, updatedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (p *ProfileHandler) GetMyProfile(c *gin.Context) {
	profileIDInterface, exists := c.Get("profile_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "profile_id not found"})
		return
	}

	profileIDStr, ok := profileIDInterface.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile_id format"})
		return
	}

	profileID, err := uuid.Parse(profileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile_id format"})
		return
	}

	fmt.Println(profileID)
	profile, err := p.ProfileService.GetProfileByID(profileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
