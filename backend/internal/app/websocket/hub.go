package websocket

import (
	"log"
	"sync"

	"github.com/pion/webrtc/v3"
	pionwebrtc "github.com/pion/webrtc/v3"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Channels   map[string]map[*Client]bool
	// PeerChannels map[string]map[string]map[*PeerConnectionState]bool
	PeerChannels map[string]map[*PeerConnectionState]bool
	// PeerChannelsServer map[string]map[string]map[*PeerConnectionState]bool
	TrackChannels map[string]map[string]*pionwebrtc.TrackLocalStaticRTP
	sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:    make(chan Message),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		Clients:      make(map[*Client]bool),
		Channels:     make(map[string]map[*Client]bool),
		PeerChannels: make(map[string]map[*PeerConnectionState]bool),
		// PeerChannels:  make(map[string]map[string]map[*PeerConnectionState]bool),
		TrackChannels: make(map[string]map[string]*webrtc.TrackLocalStaticRTP),
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

// func (h *Hub) SendToClient(client *Client, msg Message) {
// 	if pcStates, ok := h.PeerChannels[msg.ServerID][msg.Channel]; ok {
// 		for state := range pcStates {
// 			if state.client == client {
// 				select {
// 				case state.client.Send <- msg:
// 				default:
// 					close(state.client.Send)
// 					delete(h.Clients, state.client)

// 					delete(h.PeerChannels[msg.Channel][msg.ServerID], state)
// 				}
// 				return
// 			}
// 		}
// 	}
// }

// func (h *Hub) BroadcastToPeerChannel(msg Message) {
// 	if pcState, ok := h.PeerChannels[msg.ServerID][msg.Channel]; ok {
// 		for state := range pcState {
// 			select {
// 			case state.client.Send <- msg:
// 			default:
// 				close(state.client.Send)
// 				delete(h.Clients, state.client)
// 				for ch := range h.PeerChannels[msg.ServerID] {
// 					delete(h.PeerChannels[msg.ServerID][ch], state)
// 				}
// 			}
// 		}
// 	} else {
// 		log.Printf("No subscribers in channel: %s", msg.Channel)
// 	}
// }

// func (h *Hub) GetUsersFromPeerChannel(serverId string, channel string) []ContentData {
// 	h.RLock()
// 	defer h.RUnlock()

// 	var users []ContentData
// 	if peerConnStates, ok := h.PeerChannels[serverId][channel]; ok {
// 		for peerConnState := range peerConnStates {
// 			if peerConnState.client != nil {
// 				users = append(users, ContentData{
// 					Username: peerConnState.client.Username,
// 					StreamID: peerConnState.client.StreamID,
// 					ImageURL: peerConnState.client.ImageURL,
// 				})
// 			}
// 		}
// 	}
// 	return users
// }

func (h *Hub) BroadcastToPeerChannel(msg Message) {
	if pcState, ok := h.PeerChannels[msg.Channel]; ok {
		for state := range pcState {
			select {
			case state.client.Send <- msg:
			default:
				close(state.client.Send)
				delete(h.Clients, state.client)
				for ch := range h.PeerChannels {
					delete(h.PeerChannels[ch], state)
				}
			}
		}
	} else {
		log.Printf("No subscribers in channel: %s", msg.Channel)
	}
}

func (h *Hub) SendToClient(client *Client, msg Message) {
	if pcStates, ok := h.PeerChannels[msg.Channel]; ok {
		for state := range pcStates {
			if state.client == client {
				select {
				case state.client.Send <- msg:
				default:
					close(state.client.Send)
					delete(h.Clients, state.client)
					for ch := range h.PeerChannels {
						delete(h.PeerChannels[ch], state)
					}
				}
				return
			}
		}
	}
}

func (h *Hub) GetUsersFromPeerChannel(channel string) []ContentData {
	h.RLock()
	defer h.RUnlock()

	var users []ContentData
	if peerConnStates, ok := h.PeerChannels[channel]; ok {
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
