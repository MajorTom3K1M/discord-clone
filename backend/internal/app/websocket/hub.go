package websocket

import "log"

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Channels   map[string]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Channels:   make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				for channel := range h.Channels {
					delete(h.Channels[channel], client)
				}
			}
		case message := <-h.Broadcast:
			subscribers := h.Channels[message.Channel]
			for client := range subscribers {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
					for channel := range h.Channels {
						delete(h.Channels[channel], client)
					}
				}
			}
		}
	}
}

func (h *Hub) BroadcastToChannel(msg Message) {
	if clients, ok := h.Channels[msg.Channel]; ok {
		for client := range clients {
			select {
			case client.Send <- msg:
			default:
				close(client.Send)
				delete(h.Clients, client)
				for ch := range h.Channels {
					delete(h.Channels[ch], client)
				}
			}
		}
	} else {
		log.Printf("No subscribers in channel: %s", msg.Channel)
	}
}
