package websocket

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

type PeerConnectionState struct {
	peerConnection       *webrtc.PeerConnection
	client               *Client
	pendingCandidates    []webrtc.ICECandidateInit
	remoteDescriptionSet bool
	currentChannel       string
	currentServer        string
}

func NewPeerConnectionState(c *Client, serverId string, channel string) (*PeerConnectionState, error) {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return nil, err
	}

	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			return nil, err
		}
	}

	c.Hub.Lock()
	if _, ok := c.Hub.TrackChannels[channel]; !ok {
		c.Hub.TrackChannels[channel] = make(map[string]*webrtc.TrackLocalStaticRTP)
	}

	if _, ok := c.Hub.PeerChannels[serverId]; !ok {
		c.Hub.PeerChannels[serverId] = make(map[string]map[*PeerConnectionState]bool)
	}

	if _, ok := c.Hub.PeerChannels[serverId][channel]; !ok {
		c.Hub.PeerChannels[serverId][channel] = make(map[*PeerConnectionState]bool)
	}

	peerConnectionState := &PeerConnectionState{
		peerConnection:       peerConnection,
		client:               c,
		pendingCandidates:    make([]webrtc.ICECandidateInit, 0),
		remoteDescriptionSet: false,
		currentChannel:       channel,
		currentServer:        serverId,
	}

	// Add the new PeerConnectionState to PeerChannels
	c.Hub.PeerChannels[serverId][channel][peerConnectionState] = true
	c.Hub.Unlock()

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		fmt.Println("New ICE candidate found:", i.String())
		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Println(err)
			return
		}

		c.Hub.SendToClient(c, Message{
			Type:     "candidate",
			Channel:  channel,
			ServerID: serverId,
			Content: &ContentData{
				Data: string(candidateString),
			},
		})
	})

	peerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s \n", p.String())
		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Print(err)
			}
		case webrtc.PeerConnectionStateClosed:
			c.Hub.signalPeerConnections(serverId, channel)
		default:
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		fmt.Printf("Track received: %s \n", t.Kind().String())

		// Create a track to fan out our incoming video to all peers
		trackLocal := c.Hub.addTrack(serverId, channel, t)

		defer c.Hub.removeTrack(serverId, channel, trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := t.Read(buf)
			if err != nil {
				return
			}

			if _, err = trackLocal.Write(buf[:i]); err != nil {
				return
			}
		}
	})

	c.Hub.signalPeerConnections(serverId, channel)

	// return peerConnectionState, nil
	return peerConnectionState, nil
}

func (ps *PeerConnectionState) closePeerConnection() {
	log.Printf("Closing peer connection for channel %s", ps.currentChannel)

	if ps.peerConnection != nil {
		ps.peerConnection.OnICECandidate(nil)
		ps.peerConnection.OnConnectionStateChange(nil)
		ps.peerConnection.OnTrack(nil)
		ps.peerConnection.Close()
		// ps.peerConnection = nil

		ps.client.Hub.BroadcastToServer(Message{
			Type:     "participant",
			Channel:  ps.currentChannel,
			ServerID: ps.currentServer,
			Content: &ContentData{
				Data:     "left",
				Username: ps.client.Username,
				StreamID: ps.client.StreamID,
				ImageURL: ps.client.ImageURL,
				ClientID: ps.client.ID,
			},
		})
	}

	ps.client.Hub.Lock()
	defer ps.client.Hub.Unlock()
	for pcState := range ps.client.Hub.PeerChannels[ps.currentServer][ps.currentChannel] {
		if pcState.peerConnection == ps.peerConnection {
			delete(ps.client.Hub.PeerChannels[ps.currentServer][ps.currentChannel], pcState)
			break
		}
	}
	log.Printf("Peer connection closed for channel %s", ps.currentChannel)
}

// func (ps *PeerConnectionState) ChangeChannel(newServerId, newChannel string) error {
func (c *Client) ChangeChannel(newServerId, newChannel string) (*PeerConnectionState, error) {
	// log.Printf("Attempting to change channel from %s to %s", ps.currentChannel, newChannel)

	c.PeerConnectionState.closePeerConnection()

	// err := ps.initNewPeerConnection(newServerId, newChannel)
	peerConnectionState, err := NewPeerConnectionState(c, newServerId, newChannel)
	if err != nil {
		log.Printf("Error initializing new peer connection: %v", err)
		return nil, err
	}

	log.Printf("Successfully changed channel to %s", newChannel)
	return peerConnectionState, nil
}

func (ps *PeerConnectionState) initNewPeerConnection(serverId string, channel string) error {
	log.Printf("Initializing new peer connection for channel %s and server %s ", channel, serverId)

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Printf("Error creating peer connection: %v", err)
		return err
	}

	ps.peerConnection = peerConnection
	ps.currentChannel = channel
	ps.currentServer = serverId
	ps.remoteDescriptionSet = false
	ps.pendingCandidates = []webrtc.ICECandidateInit{}

	log.Printf("Adding transceivers for channel %s", channel)
	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			log.Printf("Error adding transceiver: %v", err)
			return err
		}
	}

	ps.client.Hub.Lock()
	if _, ok := ps.client.Hub.TrackChannels[channel]; !ok {
		ps.client.Hub.TrackChannels[channel] = make(map[string]*webrtc.TrackLocalStaticRTP)
	}

	if _, ok := ps.client.Hub.PeerChannels[serverId]; !ok {
		ps.client.Hub.PeerChannels[serverId] = make(map[string]map[*PeerConnectionState]bool)
	}

	if _, ok := ps.client.Hub.PeerChannels[serverId][channel]; !ok {
		ps.client.Hub.PeerChannels[serverId][channel] = make(map[*PeerConnectionState]bool)
	}

	ps.client.Hub.PeerChannels[serverId][channel][ps] = true
	ps.client.Hub.Unlock()

	log.Printf("Setting up event handlers for peer connection in channel %s", channel)
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		log.Printf("New ICE candidate found: %s", i.String())
		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Printf("Error marshaling ICE candidate: %v", err)
			return
		}

		ps.client.Hub.SendToClient(ps.client, Message{
			Type:     "candidate",
			Channel:  channel,
			ServerID: serverId,
			Content: &ContentData{
				Data: string(candidateString),
			},
		})
	})

	peerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		log.Printf("Peer Connection State has changed: %s", p.String())
		switch p {
		case webrtc.PeerConnectionStateFailed:
			log.Println("PeerConnectionStateFailed")
			if err := peerConnection.Close(); err != nil {
				log.Printf("Error closing peer connection: %v", err)
			}
		case webrtc.PeerConnectionStateClosed:
			log.Println("PeerConnectionStateClosed")
			ps.client.Hub.signalPeerConnections(serverId, channel)
		default:
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		log.Printf("Track received: %s", t.Kind().String())

		// Create a track to fan out our incoming video to all peers
		trackLocal := ps.client.Hub.addTrack(serverId, channel, t)

		defer ps.client.Hub.removeTrack(serverId, channel, trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := t.Read(buf)
			if err != nil {
				if err == io.EOF {
					log.Println("Track reading ended (EOF)")
					break
				}
				log.Printf("Error reading from track: %v", err)
				return
			}

			if _, err = trackLocal.Write(buf[:i]); err != nil {
				log.Printf("Error writing to track: %v", err)
				return
			}
		}
	})

	log.Printf("Signaling peer connections for channel %s", channel)
	ps.client.Hub.signalPeerConnections(serverId, channel)

	log.Printf("Peer connection initialized for channel %s", channel)
	return nil
}

func (ps *PeerConnectionState) SetRemoteDescription(desc webrtc.SessionDescription) error {
	if err := ps.peerConnection.SetRemoteDescription(desc); err != nil {
		return err
	}
	ps.remoteDescriptionSet = true

	// Process any pending ICE candidates
	for _, candidate := range ps.pendingCandidates {
		if err := ps.peerConnection.AddICECandidate(candidate); err != nil {
			log.Println("Failed to add ICE candidate:", err)
		}
	}
	ps.pendingCandidates = nil
	return nil
}

func (ps *PeerConnectionState) AddICECandidate(candidate webrtc.ICECandidateInit) error {
	if ps.remoteDescriptionSet {
		return ps.peerConnection.AddICECandidate(candidate)
	}
	// Queue the candidate if the remote description is not set
	ps.pendingCandidates = append(ps.pendingCandidates, candidate)
	return nil
}

// Add to list of tracks and fire renegotation for all PeerConnections
func (h *Hub) addTrack(serverId, channel string, t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	h.Lock()
	defer func() {
		h.Unlock()
		h.signalPeerConnections(serverId, channel)
	}()

	// Create a new TrackLocal with the same codec as our incoming
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}

	h.TrackChannels[channel][t.ID()] = trackLocal
	return trackLocal
}

// Remove from list of tracks and fire renegotation for all PeerConnections
func (h *Hub) removeTrack(serverId, channel string, t *webrtc.TrackLocalStaticRTP) {
	h.Lock()
	defer func() {
		h.Unlock()
		h.signalPeerConnections(serverId, channel)
	}()

	delete(h.TrackChannels[channel], t.ID())
}

func (h *Hub) signalPeerConnections(serverId, channel string) {
	h.Lock()
	defer func() {
		h.Unlock()
		h.dispatchKeyFrame(serverId, channel)
	}()

	attemptSync := func() (tryAgain bool) {
		if _, ok := h.PeerChannels[serverId]; !ok {
			return false
		}

		if _, ok := h.PeerChannels[serverId][channel]; !ok {
			return false
		}

		for pcState := range h.PeerChannels[serverId][channel] {
			if pcState.peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				delete(h.PeerChannels[serverId][channel], pcState)
				return true
			}

			existingSenders := map[string]bool{}

			for _, sender := range pcState.peerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				if _, ok := h.TrackChannels[channel][sender.Track().ID()]; !ok {
					if err := pcState.peerConnection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			for _, receiver := range pcState.peerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			for trackID := range h.TrackChannels[channel] {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := pcState.peerConnection.AddTrack(h.TrackChannels[channel][trackID]); err != nil {
						return true
					}
				}
			}

			offer, err := pcState.peerConnection.CreateOffer(nil)
			if err != nil {
				return true
			}

			if err = pcState.peerConnection.SetLocalDescription(offer); err != nil {
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				return true
			}

			currentClient := pcState.client
			currentClient.Hub.SendToClient(
				currentClient,
				Message{
					Type:     "offer",
					Channel:  channel,
					ServerID: serverId,
					Content: &ContentData{
						Data: string(offerString),
					},
				},
			)
		}
		return
	}

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			go func() {
				time.Sleep(time.Second * 3)
				h.signalPeerConnections(serverId, channel)
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}

func (h *Hub) dispatchKeyFrame(serverId, channel string) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.PeerChannels[serverId]; ok {
		if _, ok := h.PeerChannels[serverId][channel]; ok {
			for state := range h.PeerChannels[serverId][channel] {
				for _, receiver := range state.peerConnection.GetReceivers() {
					if receiver.Track() == nil {
						continue
					}

					_ = state.peerConnection.WriteRTCP([]rtcp.Packet{
						&rtcp.PictureLossIndication{
							MediaSSRC: uint32(receiver.Track().SSRC()),
						},
					})
				}
			}
		}
	}
}
