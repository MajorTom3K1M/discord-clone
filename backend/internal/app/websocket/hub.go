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
	ClientMessage   chan ClientMessage
	Register        chan *Client
	RegisterServer  chan ClientMessage
	Unregister      chan *Client
	UnregisterPeer  chan *PeerConnectionState
	Channels        map[string]map[*Client]bool
	Servers         map[string]map[*Client]bool
	PeerChannels    map[string]map[string]map[*PeerConnectionState]bool
	TrackChannels   map[string]map[string]*pionwebrtc.TrackLocalStaticRTP
	sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		BroadcastServer: make(chan Message),
		Broadcast:       make(chan Message),
		ClientMessage:   make(chan ClientMessage),
		Register:        make(chan *Client),
		RegisterServer:  make(chan ClientMessage),
		Unregister:      make(chan *Client),
		UnregisterPeer:  make(chan *PeerConnectionState),
		Clients:         make(map[*Client]bool),
		Channels:        make(map[string]map[*Client]bool),
		Servers:         make(map[string]map[*Client]bool),
		PeerChannels:    make(map[string]map[string]map[*PeerConnectionState]bool),
		TrackChannels:   make(map[string]map[string]*webrtc.TrackLocalStaticRTP),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				log.Printf("Closing Client : %s", client.ID)
				delete(h.Clients, client)
				close(client.Send)
				h.cleanupClient(client)
			}
		case peer := <-h.UnregisterPeer:
			if _, ok := h.PeerChannels[peer.currentServer][peer.currentChannel][peer]; ok {
				delete(h.PeerChannels[peer.currentServer][peer.currentChannel], peer)
			}
		case clientMessage := <-h.RegisterServer:
			if _, ok := h.Servers[clientMessage.ServerID]; !ok {
				h.Servers[clientMessage.ServerID] = make(map[*Client]bool)
			}

			h.Servers[clientMessage.ServerID][clientMessage.Client] = true
		case message := <-h.Broadcast:
			subscribers := h.Channels[message.Channel]
			for client := range subscribers {
				select {
				case client.Send <- message:
				default:
					log.Printf("Closing Broadcasting to channel : %s", message.Channel)
					close(client.Send)
					delete(h.Clients, client)
					h.cleanupClient(client)
				}
			}
		case message := <-h.BroadcastServer:
			log.Println("Broadcasting to server : %s", message.ServerID)
			for client := range h.Servers[message.ServerID] {
				client.Add(1)
				go func(client *Client) {
					defer client.Done()
					select {
					case client.Send <- message:
					default:
						log.Println("Closing Broadcasting to server : %s", message.ServerID)
						client.Wait()
						close(client.Send)
						delete(h.Clients, client)
						h.cleanupClient(client)
					}
				}(client)
				// select {
				// case client.Send <- message:
				// default:
				// 	log.Println("Closing Broadcasting to server : %s", message.ServerID)
				// 	close(client.Send)
				// 	delete(h.Clients, client)
				// 	h.cleanupClient(client)
				// }
			}
		case clientMessage := <-h.ClientMessage:
			client := clientMessage.Client
			client.Add(1)
			go func(client *Client) {
				defer client.Done()
				select {
				case client.Send <- clientMessage.Message:
				default:
					log.Printf("Closing Client : %s", client.ID)
					client.Wait()
					close(client.Send)
					delete(h.Clients, client)
					h.cleanupClient(client)
				}
			}(client)
			// select {
			// case client.Send <- clientMessage.Message:
			// default:
			// 	log.Printf("Closing Client : %s", client.ID)
			// 	close(client.Send)
			// 	delete(h.Clients, client)
			// 	h.cleanupClient(client)
			// }
		}
	}
}

func (h *Hub) cleanupClient(client *Client) {
	for channel := range h.Channels {
		delete(h.Channels[channel], client)
	}
	for server := range h.Servers {
		delete(h.Servers[server], client)
	}
}

func (h *Hub) BroadcastToChannel(msg Message) {
	if clients, ok := h.Channels[msg.Channel]; ok {
		for client := range clients {
			select {
			case client.Send <- msg:
			default:
				log.Println("BroadcastToChannel cause CLOSE")
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
