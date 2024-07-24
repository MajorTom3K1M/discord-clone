package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type Client struct {
	Hub                 *Hub
	Conn                *websocket.Conn
	Send                chan Message
	ID                  string
	Username            string
	ProfileID           uuid.UUID
	PeerConnectionState *PeerConnectionState
	StreamID            string
	ImageURL            string
	sync.Mutex
}

type ContentInterface interface{}

// type Content struct {
// 	Message string `json:"message"`
// 	FileUrl string `json:"fileUrl"`
// }

type ContentData struct {
	Data     string `json:"data"`
	Username string `json:"username,omitempty"`
	StreamID string `json:"streamId,omitempty"`
	ImageURL string `json:"imageURL,omitempty"`
}

type Message struct {
	Type     string           `json:"type"`
	Channel  string           `json:"channel,omitempty"`
	ServerID string           `json:"channel,omitempty"`
	Content  ContentInterface `json:"content,omitempty"`
}

type WebRTCMessage struct {
	Offer     webrtc.SessionDescription `json:"offer,omitempty"`
	Answer    webrtc.SessionDescription `json:"answer,omitempty"`
	Candidate webrtc.ICECandidateInit   `json:"candidate,omitempty"`
	StreamID  string                    `json:"streamId,omitempty"`
}

const (
	maxMessageSize = 10240
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
		if c.PeerConnectionState != nil {
			c.PeerConnectionState.closePeerConnection()
		}
	}()

	var err error

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			}
			log.Printf("Read error: %v", err)
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
		case "unsubscribe":
			if msg.Channel != "" {
				if clients, ok := c.Hub.Channels[msg.Channel]; ok {
					delete(clients, c)
					if len(clients) == 0 {
						delete(c.Hub.Channels, msg.Channel)
					}
					log.Printf("Client %s unsubscribed from channel %s", c.ID, msg.Channel)
				}
			}
		case "message":
			c.Hub.BroadcastToChannel(msg)
		case "initializeCall":
			if c.PeerConnectionState == nil {
				log.Printf("Client %s initializeCall", msg.Channel)
				webrtcMsg := msg.Content.(WebRTCMessage)
				c.StreamID = webrtcMsg.StreamID
				c.PeerConnectionState, err = NewPeerConnectionState(c, msg.Channel)
				if err != nil {
					log.Println("Error creating PeerConnection:", err)
					continue
				}
				// c.PeerConnectionState.streamID = webrtcMsg.StreamID
			} else if c.PeerConnectionState.currentChannel != msg.Channel {
				log.Printf("Client %s changing channel from %s to %s", c.ID, c.PeerConnectionState.currentChannel, msg.Channel)
				webrtcMsg := msg.Content.(WebRTCMessage)
				c.StreamID = webrtcMsg.StreamID
				err = c.PeerConnectionState.ChangeChannel(msg.Channel)
				if err != nil {
					log.Println("Error changing channel:", err)
					continue
				}
				// webrtcMsg := msg.Content.(WebRTCMessage)
				// c.PeerConnectionState.streamID = webrtcMsg.StreamID
			}
		case "answer":
			if c.PeerConnectionState != nil {
				webrtcMsg := msg.Content.(WebRTCMessage)
				if err := c.PeerConnectionState.SetRemoteDescription(webrtcMsg.Answer); err != nil {
					log.Println("Failed to set remote description:", err)
				}
			}
		case "candidate":
			if c.PeerConnectionState != nil {
				webrtcMsg := msg.Content.(WebRTCMessage)
				if err := c.PeerConnectionState.AddICECandidate(webrtcMsg.Candidate); err != nil {
					log.Println("Failed to add ICE candidate:", err)
				}

				c.Hub.BroadcastToPeerChannel(Message{
					Type:    "participant",
					Channel: msg.Channel,
					Content: c.Hub.GetUsersFromPeerChannel(msg.Channel),
				})
			}
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

func (c *Client) WriteJSON(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.Conn.WriteJSON(v)
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		Content json.RawMessage `json:"content,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Determine the type of content based on the message type
	switch aux.Type {
	case "candidate", "answer", "initializeCall":
		var webRTCMessage WebRTCMessage
		if err := json.Unmarshal(aux.Content, &webRTCMessage); err != nil {
			return err
		}
		m.Content = webRTCMessage
	default:
		m.Content = aux.Content
	}

	return nil
}
