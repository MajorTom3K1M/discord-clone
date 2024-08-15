package websocket

import (
	"log"
	"sync"

	"github.com/pion/webrtc/v3"
	pionwebrtc "github.com/pion/webrtc/v3"
)

type Hub struct {
	Clients         map[*Client]bool
	BroadcastServer chan Message
	Broadcast       chan Message
	Register        chan *Client
	Unregister      chan *Client
	Channels        map[string]map[*Client]bool
	Servers         map[string]map[*Client]bool
	PeerChannels    map[string]map[string]map[*PeerConnectionState]bool
	TrackChannels   map[string]map[string]*pionwebrtc.TrackLocalStaticRTP
	// Channels      map[string]map[*Client]bool
	// PeerChannels map[string]map[*PeerConnectionState]bool
	sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		BroadcastServer: make(chan Message),
		Broadcast:       make(chan Message),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		Clients:         make(map[*Client]bool),
		Channels:        make(map[string]map[*Client]bool),
		Servers:         make(map[string]map[*Client]bool),
		PeerChannels:    make(map[string]map[string]map[*PeerConnectionState]bool),
		TrackChannels:   make(map[string]map[string]*webrtc.TrackLocalStaticRTP),
		// PeerChannels: make(map[string]map[*PeerConnectionState]bool),
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
				for server := range h.Servers {
					delete(h.Servers[server], client)
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
		case message := <-h.BroadcastServer:
			log.Println("Broadcasting to server : %s", message.ServerID)
			for client := range h.Servers[message.ServerID] {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
					for server := range h.Servers {
						delete(h.Servers[server], client)
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

// func (h *Hub) SendToClient(client *Client, msg Message) {
// 	if pcStates, ok := h.PeerChannels[msg.ServerID][msg.Channel]; ok {
// 		for state := range pcStates {
// 			if state.client == client {
// 				select {
// 				case state.client.Send <- msg:
// 				default:
// 					close(state.client.Send)
// 					delete(h.Clients, state.client)

// 					delete(h.PeerChannels[msg.ServerID][msg.Channel], state)
// 				}
// 				return
// 			}
// 		}
// 	}
// }

func (h *Hub) SendToClient(client *Client, msg Message) {
	client.Lock()
	defer client.Unlock()
	select {
	case client.Send <- msg:
	default:
		close(client.Send)
		delete(h.Clients, client)
		delete(h.PeerChannels[msg.ServerID][msg.Channel], client.PeerConnectionState)
	}
}

func (h *Hub) BroadcastToPeerChannel(msg Message) {
	if pcState, ok := h.PeerChannels[msg.ServerID][msg.Channel]; ok {
		for state := range pcState {
			select {
			case state.client.Send <- msg:
			default:
				close(state.client.Send)
				delete(h.Clients, state.client)
				for ch := range h.PeerChannels[msg.ServerID] {
					delete(h.PeerChannels[msg.ServerID][ch], state)
				}
			}
		}
	} else {
		log.Printf("No subscribers in channel: %s", msg.Channel)
	}
}

func (h *Hub) BroadcastToServer(msg Message) {
	h.RLock()
	defer h.RUnlock()
	if clients, ok := h.Servers[msg.ServerID]; ok {
		for client := range clients {
			select {
			case client.Send <- msg:
			default:
				close(client.Send)
				delete(h.Clients, client)
				for sv := range h.Servers {
					delete(h.Servers[sv], client)
				}
			}
		}
	} else {
		log.Printf("No subscribers in server: %s", msg.ServerID)
	}
}

func (h *Hub) BroadcastToPeerServer(msg Message) {
	if channels, ok := h.PeerChannels[msg.ServerID]; ok {
		for ch := range channels {
			for state := range h.PeerChannels[msg.ServerID][ch] {
				select {
				case state.client.Send <- msg:
				default:
					close(state.client.Send)
					delete(h.Clients, state.client)
					delete(h.PeerChannels[msg.ServerID][ch], state)
				}
			}
		}
	} else {
		log.Printf("No subscribers in server: %s", msg.ServerID)
	}
}

func (h *Hub) GetUsersFromPeerChannel(serverId string, channel string) []ContentData {
	h.RLock()
	defer h.RUnlock()

	var users []ContentData
	if peerConnStates, ok := h.PeerChannels[serverId][channel]; ok {
		for peerConnState := range peerConnStates {
			if peerConnState.client != nil {
				users = append(users, ContentData{
					Username: peerConnState.client.Username,
					StreamID: peerConnState.client.StreamID,
					ImageURL: peerConnState.client.ImageURL,
				})
			}
		}
	}
	return users
}

func (h *Hub) GetUsersFromPeerChannelsServer(serverId string) map[string]map[string]ContentData {
	h.RLock()
	defer h.RUnlock()

	var users map[string]map[string]ContentData = make(map[string]map[string]ContentData)
	for ch := range h.PeerChannels[serverId] {
		if peerConnStates, ok := h.PeerChannels[serverId][ch]; ok {
			for peerConnState := range peerConnStates {
				if peerConnState.client != nil {
					if _, ok := users[ch]; !ok {
						users[ch] = make(map[string]ContentData)
					}

					users[ch][peerConnState.client.StreamID] = ContentData{
						Data:     "joined",
						StreamID: peerConnState.client.StreamID,
						Username: peerConnState.client.Username,
						ImageURL: peerConnState.client.ImageURL,
						ClientID: peerConnState.client.ID,
					}
				}
			}
		}
	}

	return users
}
