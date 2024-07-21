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
}

func NewPeerConnectionState(c *Client, channel string) (*PeerConnectionState, error) {
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
	c.Hub.PeerChannels[channel] = append(c.Hub.PeerChannels[channel], PeerConnectionState{peerConnection, c, make([]webrtc.ICECandidateInit, 0), false, channel})
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

		c.Lock()
		c.Hub.SendToClient(c, Message{
			Type:    "candidate",
			Channel: channel,
			Content: &ContentData{
				Data: string(candidateString),
			},
		})

		c.Unlock()
	})

	peerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s \n", p.String())
		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Print(err)
			}
		case webrtc.PeerConnectionStateClosed:
			c.Hub.signalPeerConnections(channel)
		default:
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		fmt.Printf("Track received: %s \n", t.Kind().String())

		// Create a track to fan out our incoming video to all peers
		trackLocal := c.Hub.addTrack(channel, t)

		defer c.Hub.removeTrack(channel, trackLocal)

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

	c.Hub.signalPeerConnections(channel)

	return &PeerConnectionState{
		peerConnection: peerConnection,
		client:         c,
	}, nil
}

func (ps *PeerConnectionState) closePeerConnection() {
	log.Printf("Closing peer connection for channel %s", ps.currentChannel)

	if ps.peerConnection != nil {
		ps.peerConnection.OnICECandidate(nil)
		ps.peerConnection.OnConnectionStateChange(nil)
		ps.peerConnection.OnTrack(nil)
		ps.peerConnection.Close()
		ps.peerConnection = nil
	}

	ps.client.Hub.Lock()
	defer ps.client.Hub.Unlock()
	for i, pcState := range ps.client.Hub.PeerChannels[ps.currentChannel] {
		if pcState.peerConnection == ps.peerConnection {
			ps.client.Hub.PeerChannels[ps.currentChannel] = append(ps.client.Hub.PeerChannels[ps.currentChannel][:i], ps.client.Hub.PeerChannels[ps.currentChannel][i+1:]...)
			break
		}
	}
	log.Printf("Peer connection closed for channel %s", ps.currentChannel)
}

func (ps *PeerConnectionState) ChangeChannel(newChannel string) error {
	log.Printf("Attempting to change channel from %s to %s", ps.currentChannel, newChannel)

	ps.closePeerConnection()

	err := ps.initNewPeerConnection(newChannel)
	if err != nil {
		log.Printf("Error initializing new peer connection: %v", err)
		return err
	}

	log.Printf("Successfully changed channel to %s", newChannel)
	return nil
}

func (ps *PeerConnectionState) initNewPeerConnection(channel string) error {
	log.Printf("Initializing new peer connection for channel %s", channel)

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
	ps.client.Hub.PeerChannels[channel] = append(ps.client.Hub.PeerChannels[channel], *ps)
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

		ps.client.Lock()
		defer ps.client.Unlock()

		ps.client.Hub.SendToClient(ps.client, Message{
			Type:    "candidate",
			Channel: channel,
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
			ps.client.Hub.signalPeerConnections(channel)
		default:
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		log.Printf("Track received: %s", t.Kind().String())

		// Create a track to fan out our incoming video to all peers
		trackLocal := ps.client.Hub.addTrack(channel, t)

		defer ps.client.Hub.removeTrack(channel, trackLocal)

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
	ps.client.Hub.signalPeerConnections(channel)

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
func (h *Hub) addTrack(channel string, t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	h.Lock()
	defer func() {
		h.Unlock()
		h.signalPeerConnections(channel)
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
func (h *Hub) removeTrack(channel string, t *webrtc.TrackLocalStaticRTP) {
	h.Lock()
	defer func() {
		h.Unlock()
		h.signalPeerConnections(channel)
	}()

	delete(h.TrackChannels[channel], t.ID())
}

func (h *Hub) signalPeerConnections(channel string) {
	h.Lock()
	defer func() {
		h.Unlock()
		h.dispatchKeyFrame(channel)
	}()

	attemptSync := func() (tryAgain bool) {
		for i := range h.PeerChannels[channel] {
			if h.PeerChannels[channel][i].peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				h.PeerChannels[channel] = append(h.PeerChannels[channel][:i], h.PeerChannels[channel][i+1:]...)
				return true
			}

			existingSenders := map[string]bool{}

			for _, sender := range h.PeerChannels[channel][i].peerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				if _, ok := h.TrackChannels[channel][sender.Track().ID()]; !ok {
					if err := h.PeerChannels[channel][i].peerConnection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			for _, receiver := range h.PeerChannels[channel][i].peerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			for trackID := range h.TrackChannels[channel] {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := h.PeerChannels[channel][i].peerConnection.AddTrack(h.TrackChannels[channel][trackID]); err != nil {
						return true
					}
				}
			}

			offer, err := h.PeerChannels[channel][i].peerConnection.CreateOffer(nil)
			if err != nil {
				return true
			}

			if err = h.PeerChannels[channel][i].peerConnection.SetLocalDescription(offer); err != nil {
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				return true
			}

			currentClient := h.PeerChannels[channel][i].client
			currentClient.Hub.SendToClient(
				currentClient,
				Message{
					Type:    "offer",
					Channel: channel,
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
				h.signalPeerConnections(channel)
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}

func (h *Hub) dispatchKeyFrame(channel string) {
	h.Lock()
	defer h.Unlock()

	for i := range h.PeerChannels[channel] {
		for _, receiver := range h.PeerChannels[channel][i].peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			_ = h.PeerChannels[channel][i].peerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}
