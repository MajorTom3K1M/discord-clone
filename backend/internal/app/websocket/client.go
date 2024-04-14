package websocket

import (
	"discord-backend/internal/app/models"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan Message
	ID   string
}

type Content struct {
	Message string `json:"message"`
	FileUrl string `json:"fileUrl"`
}
type Message struct {
	Type    string         `json:"type"`              // e.g., "chat", "subscribe", "unsubscribe"
	Channel string         `json:"channel,omitempty"` // Channel or conversation ID
	Content models.Message `json:"content,omitempty"` // Actual message content
}

const (
	maxMessageSize = 512
	pongWait       = 60 * time.Second
	writeWait      = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch msg.Type {
		case "subscribe":
			if msg.Channel != "" {
				if _, ok := c.Hub.Channels[msg.Channel]; !ok {
					c.Hub.Channels[msg.Channel] = make(map[*Client]bool)
				}
				c.Hub.Channels[msg.Channel][c] = true
				log.Printf("Client %s subscribed to channel %s", c.ID, msg.Channel)
			}
		case "message":
			c.Hub.BroadcastToChannel(msg)
		default:
			log.Printf("Unknown message type received: %v", msg.Type)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			jsonMessage, err := json.Marshal(message)
			if err != nil {
				return
			}
			w.Write(jsonMessage)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				nextMsg, _ := json.Marshal(<-c.Send)
				w.Write(nextMsg)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// var newline = []byte{'\n'}

// func (c *Client) ReadPump(manager *Manager) {
// 	defer func() {
// 		manager.unregister <- c
// 		c.Conn.Close()
// 	}()

// 	c.Conn.SetReadLimit(maxMessageSize)              // Define maxMessageSize as needed
// 	c.Conn.SetReadDeadline(time.Now().Add(pongWait)) // pongWait is the duration to wait for a pong message
// 	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

// 	for {
// 		_, message, err := c.Conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}

// 		// Deserialize the message to a Message struct
// 		var msg Message
// 		err = json.Unmarshal(message, &msg)
// 		if err != nil {
// 			log.Println("error parsing message:", err)
// 			continue
// 		}

// 		// Handle the message based on its type
// 		switch msg.Type {
// 		case "chat":
// 			manager.BroadcastToChannel(msg.Channel, []byte(msg.Content))
// 		case "subscribe":
// 			manager.SubscribeClientToChannel(c, msg.Channel)
// 		case "unsubscribe":
// 			manager.UnsubscribeClientFromChannel(c, msg.Channel)
// 		}
// 	}
// }

// func (c *Client) WritePump() {
// 	ticker := time.NewTicker(pongWait)
// 	defer func() {
// 		ticker.Stop()
// 		c.Conn.Close()
// 	}()

// 	for {
// 		select {
// 		case message, ok := <-c.Send:
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if !ok {
// 				// The Send channel was closed
// 				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			w, err := c.Conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			w.Write(message)

// 			n := len(c.Send)
// 			for i := 0; i < n; i++ {
// 				w.Write(newline) // Assume newline is defined
// 				w.Write(<-c.Send)
// 			}

// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
