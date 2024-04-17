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
	ServerService        *services.ServerService
	ConversationService  *services.ConversationService
	ChannelService       *services.ChannelService
	MessageService       *services.MessageService
	DirectMessageService *services.DirectMessageService
}

func NewWebsocketHandler(
	serverService *services.ServerService,
	conversationService *services.ConversationService,
	channelService *services.ChannelService,
	messageService *services.MessageService,
	directMessageService *services.DirectMessageService,
) *WebsocketHandler {
	return &WebsocketHandler{
		ServerService:        serverService,
		ConversationService:  conversationService,
		ChannelService:       channelService,
		MessageService:       messageService,
		DirectMessageService: directMessageService,
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

		if err := message.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content in message is missing"})
			return
		}

		channelKey := fmt.Sprintf("chat:%s:messages", channelIDStr)
		msg := ws.Message{
			Type:    "message",
			Channel: channelKey,
			Content: message,
		}
		hub.BroadcastToChannel(msg)

		c.JSON(http.StatusOK, gin.H{"message": "Message created successfully", "data": message})
	}
}

func (h *WebsocketHandler) WebScoketEditMessageHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramMessageID := c.Param("messageId")
		messageID, err := uuid.Parse(paramMessageID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Server UUID format"})
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

		message, err := h.MessageService.GetMessage(channelID, messageID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}

		isMessageOwner := message.MemberID == member.ID
		// isAdmin := member.Role == models.Admin
		// isModerator := member.Role == models.Moderator
		// canModify := isMessageOwner || isAdmin || isModerator

		var input struct {
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if !isMessageOwner {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		message, err = h.MessageService.UpdateMessage(channelID, messageID, input.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating message: " + err.Error()})
			return
		}

		if err := message.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content in message is missing"})
			return
		}

		channelKey := fmt.Sprintf("chat:%s:messages:update", channelIDStr)
		msg := ws.Message{
			Type:    "message",
			Channel: channelKey,
			Content: message,
		}
		hub.BroadcastToChannel(msg)

		c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully", "data": message})
	}
}

func (h *WebsocketHandler) WebScoketDeleteMessageHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramMessageID := c.Param("messageId")
		messageID, err := uuid.Parse(paramMessageID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Server UUID format"})
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

		message, err := h.MessageService.GetMessage(channelID, messageID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}

		isMessageOwner := message.MemberID == member.ID
		isAdmin := member.Role == models.Admin
		isModerator := member.Role == models.Moderator
		canModify := isMessageOwner || isAdmin || isModerator

		if !canModify {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		message, err = h.MessageService.DeleteMessage(channelID, messageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating message: " + err.Error()})
			return
		}

		if err := message.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content in message is missing"})
			return
		}

		channelKey := fmt.Sprintf("chat:%s:messages:update", channelIDStr)
		msg := ws.Message{
			Type:    "message",
			Channel: channelKey,
			Content: message,
		}
		hub.BroadcastToChannel(msg)

		c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully", "data": message})
	}
}

func (h *WebsocketHandler) WebSocketDirectMessageHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Content string `json:"content"`
			FileURL string `json:"fileUrl"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		conversationIDStr := c.Query("conversationId")
		if conversationIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing conversationId"})
			return
		}

		conversationID, err := uuid.Parse(conversationIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversationId"})
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

		conversation, err := h.ConversationService.GetConversation(conversationID, profileID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting conversation: " + err.Error()})
			return
		}

		var member models.Member
		if conversation.MemberOne.ProfileID == profileID {
			member = conversation.MemberOne
		} else if conversation.MemberTwo.ProfileID == profileID {
			member = conversation.MemberTwo
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}

		directMessage, err := h.DirectMessageService.CreateDirectMessage(conversationID, member.ID, input.Content, input.FileURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create direct message"})
			return
		}

		if err := directMessage.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content in direct message is missing"})
			return
		}

		channelKey := fmt.Sprintf("chat:%s:messages", conversationIDStr)
		msg := ws.Message{
			Type:    "message",
			Channel: channelKey,
			Content: directMessage,
		}
		hub.BroadcastToChannel(msg)

		c.JSON(http.StatusOK, gin.H{"message": "Direct Message created successfully", "data": directMessage})
	}
}

func (h *WebsocketHandler) WebSocketEditDirectMessageHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramDirectMessageID := c.Param("directMessageId")
		directMessageID, err := uuid.Parse(paramDirectMessageID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Server UUID format"})
			return
		}

		conversationIDStr := c.Query("conversationId")
		if conversationIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing conversationId"})
			return
		}

		conversationID, err := uuid.Parse(conversationIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversationId"})
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

		conversation, err := h.ConversationService.GetConversation(conversationID, profileID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting conversation: " + err.Error()})
			return
		}

		var member models.Member
		if conversation.MemberOne.ProfileID == profileID {
			member = conversation.MemberOne
		} else if conversation.MemberTwo.ProfileID == profileID {
			member = conversation.MemberTwo
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}

		directMessage, err := h.DirectMessageService.GetDirectMessage(conversationID, directMessageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create direct message"})
			return
		}

		isMessageOwner := directMessage.MemberID == member.ID
		canModify := isMessageOwner

		var input struct {
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if !canModify {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		directMessage, err = h.DirectMessageService.UpdateDirectMessage(conversationID, directMessageID, input.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating direct message: " + err.Error()})
			return
		}

		if err := directMessage.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content in direct message is missing"})
			return
		}

		channelKey := fmt.Sprintf("chat:%s:messages:update", conversationIDStr)
		msg := ws.Message{
			Type:    "message",
			Channel: channelKey,
			Content: directMessage,
		}
		hub.BroadcastToChannel(msg)

		c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully", "data": directMessage})
	}
}

func (h *WebsocketHandler) WebSocketDeleteDirectMessageHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramDirectMessageID := c.Param("directMessageId")
		directMessageID, err := uuid.Parse(paramDirectMessageID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Server UUID format"})
			return
		}

		conversationIDStr := c.Query("conversationId")
		if conversationIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing conversationId"})
			return
		}

		conversationID, err := uuid.Parse(conversationIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversationId"})
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

		conversation, err := h.ConversationService.GetConversation(conversationID, profileID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting conversation: " + err.Error()})
			return
		}

		var member models.Member
		if conversation.MemberOne.ProfileID == profileID {
			member = conversation.MemberOne
		} else if conversation.MemberTwo.ProfileID == profileID {
			member = conversation.MemberTwo
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}

		directMessage, err := h.DirectMessageService.GetDirectMessage(conversationID, directMessageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create direct message"})
			return
		}

		isMessageOwner := directMessage.MemberID == member.ID
		canModify := isMessageOwner

		if !canModify {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		directMessage, err = h.DirectMessageService.DeleteDirectMessage(conversationID, directMessageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating direct message: " + err.Error()})
			return
		}

		if err := directMessage.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content in direct message is missing"})
			return
		}

		channelKey := fmt.Sprintf("chat:%s:messages:update", conversationIDStr)
		msg := ws.Message{
			Type:    "message",
			Channel: channelKey,
			Content: directMessage,
		}
		hub.BroadcastToChannel(msg)

		c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully", "data": directMessage})
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
