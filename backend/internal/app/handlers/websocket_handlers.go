package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"discord-backend/internal/app/models"
	"discord-backend/internal/app/services"
	ws "discord-backend/internal/app/websocket"
)

type WebsocketHandler struct {
	ServerService  *services.ServerService
	ChannelService *services.ChannelService
	MessageService *services.MessageService
}

func NewWebsocketHandler(
	serverService *services.ServerService,
	channelService *services.ChannelService,
	messageService *services.MessageService,
) *WebsocketHandler {
	return &WebsocketHandler{
		ServerService:  serverService,
		ChannelService: channelService,
		MessageService: messageService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			http.Error(c.Writer, "Could not upgrade to WebSocket", http.StatusBadRequest)
			return
		}
		client := &ws.Client{Hub: hub, Conn: conn, Send: make(chan ws.Message), ID: c.Request.RemoteAddr}
		hub.Register <- client

		go client.ReadPump()
		go client.WritePump()
	}
}

func (h *WebsocketHandler) WebSocketMessageHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Content string `json:"content"`
			FileURL string `json:"fileUrl"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		serverIDStr := c.Query("serverId")
		channelIDStr := c.Query("channelId")
		if serverIDStr == "" || channelIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing serverId or channelId"})
			return
		}

		serverID, err := uuid.Parse(serverIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid serverId"})
			return
		}

		channelID, err := uuid.Parse(channelIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channelId"})
			return
		}

		profileIDInterface, exists := c.Get("profile_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "profile_id not found"})
			return
		}

		profileIDStr, ok := profileIDInterface.(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID format"})
			return
		}

		profileID, err := uuid.Parse(profileIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
			return
		}

		server, err := h.ServerService.GetServer(profileID, serverID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting server: " + err.Error()})
			return
		}

		_, err = h.ChannelService.GetChannel(channelID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting channel: " + err.Error()})
			return
		}

		member, err := FindMember(server.Members, profileID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}

		message, err := h.MessageService.CreateMessage(channelID, member.ID, input.Content, input.FileURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
			return
		}

		channelKey := fmt.Sprintf("chat:%s:messages", channelIDStr)
		fmt.Println(channelKey)
		msg := ws.Message{
			Type:    "message",
			Channel: channelKey,
			Content: ws.Content{
				Message: message.Content,
				FileUrl: *message.FileURL,
			},
		}
		hub.BroadcastToChannel(msg)

		c.JSON(http.StatusOK, gin.H{"message": "Message created successfully", "data": message})
	}
}

func FindMember(members []models.Member, profileID uuid.UUID) (*models.Member, error) {
	for _, member := range members {
		if member.ProfileID == profileID {
			return &member, nil
		}
	}
	return nil, errors.New("member not found")
}
