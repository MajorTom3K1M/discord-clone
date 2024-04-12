package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	ws "discord-backend/internal/app/websocket"
)

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
