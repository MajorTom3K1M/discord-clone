package handlers

import (
	"discord-backend/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServerHandler struct {
	ServerService *services.ServerService
}

func NewServerHandler(serverService *services.ServerService) *ServerHandler {
	return &ServerHandler{ServerService: serverService}
}

func (s *ServerHandler) CreateServer(c *gin.Context) {
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

	var serverData struct {
		Name     string `json:"name"`
		ImageURL string `json:"imageUrl"`
	}
	if err := c.ShouldBindJSON(&serverData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body: " + err.Error()})
		return
	}

	server, err := s.ServerService.CreateServer(profileID, serverData.Name, serverData.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating server: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Server created successfully", "server": server})
}

func (s *ServerHandler) GetServers(c *gin.Context) {
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

	servers, err := s.ServerService.GetServers(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting servers: " + err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get servers successfully", "servers": servers})
}

func (s *ServerHandler) GetServer(c *gin.Context) {
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

	server, err := s.ServerService.GetServer(profileID, serverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting server: " + err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get server successfully", "server": server})
}

func (s *ServerHandler) GetServerDetails(c *gin.Context) {
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

	server, err := s.ServerService.GetServerDetails(profileID, serverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting server: " + err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get server successfully", "server": server})
}
