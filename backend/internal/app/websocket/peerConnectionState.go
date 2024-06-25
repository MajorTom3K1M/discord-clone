package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

type PeerConnectionState struct {
	peerConnection *webrtc.PeerConnection
	client         *Client
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
	c.Hub.PeerChannels[channel] = append(c.Hub.PeerChannels[channel], PeerConnectionState{peerConnection, c})
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
		if writeErr := c.Conn.WriteJSON(&Message{
			Type:    "candidate",
			Channel: channel,
			Content: &ContentData{
				Data: string(candidateString),
			},
		}); writeErr != nil {
			log.Println(writeErr)
		}
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

			if err = h.PeerChannels[channel][i].client.WriteJSON(&Message{
				Type:    "offer",
				Channel: channel,
				Content: &ContentData{
					Data: string(offerString),
				},
			}); err != nil {
				return true
			}
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
