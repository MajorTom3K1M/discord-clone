"use client";
import { createContext, use, useContext, useEffect, useRef, useState } from "react";
import { useWebSocket } from "./SocketProvider";

interface Message {
    type: string;
    channel: string;
    serverId: string;
    content: {
        data?: string;
        username?: string;
        streamId?: string;
        imageURL?: string;
    }
};

interface ChannelConfig {
    channel: string;
    serverId: string;
};

interface StreamConfig {
    audio: boolean;
    video: boolean;
}

type WebRTCContextType = {
    localStream: MediaStream | null;
    remoteStreams: MediaStream[];
    isConnected: boolean;
    joinChannel: (channel: string, serverId: string, config: StreamConfig) => void;
    closeChannel: () => void;
};

const WebRTCContext = createContext<WebRTCContextType>({
    localStream: null,
    remoteStreams: [],
    isConnected: false,
    joinChannel: () => { },
    closeChannel: () => { }
});

export const useWebRTC = () => {
    const context = useContext(WebRTCContext);
    if (!context) {
        throw new Error('useWebRTC must be used within a WebRTCProvider');
    }
    return context;
};

export const WebRTCProvider = ({
    children
}: {
    children: React.ReactNode
}) => {
    const { socket, isConnected, sendWebRTCMessage } = useWebSocket();
    const [localStream, setLocalStream] = useState<MediaStream | null>(null);
    const [remoteStreams, setRemoteStreams] = useState<MediaStream[]>([]);

    const pcRef = useRef<RTCPeerConnection | null>(null);

    const configuration: RTCConfiguration = {
        iceServers: [
            { urls: 'stun:stun.l.google.com:19302' }
        ]
    };

    useEffect(() => {
        const handleMessage = (event: MessageEvent) => {
            const message: Message = JSON.parse(event.data);

            if (!message) {
                return console.log('Failed to parse message');
            }

            switch (message.type) {
                case 'offer':
                    handleOfferMessage(message);
                    break;
                case 'candidate':
                    handleCandidateMessage(message);
                    break;
                default:
                    console.log('Unknown message type:', message.type);
            }
        };

        if (isConnected) {
            socket!.addEventListener('message', handleMessage);

            return () => {
                socket!.removeEventListener('message', handleMessage);
            };
        }
    }, [isConnected]);

    const createPeerConnection = async (channel: string, serverId: string) => {
        pcRef.current = new RTCPeerConnection(configuration);

        pcRef.current.ontrack = (event) => {
            console.log("ontrack ", event);
            
            // Check if the stream associated with the track is already in the remoteStreams
            setRemoteStreams(prevStreams => {
                if (prevStreams.find(stream => stream.id === event.streams[0].id)) {
                    return prevStreams;
                }
                return [...prevStreams, event.streams[0]];
            });

            event.streams[0].onremovetrack = () => {
                setRemoteStreams(prevStreams => prevStreams.filter(stream => stream.id !== event.streams[0].id));
            };
        };

        pcRef.current.onicecandidate = (event) => {
            if (event.candidate) {
                console.log('ICE candidate:', event.candidate);
                sendWebRTCMessage('candidate', channel, serverId, { candidate: event.candidate });
            }
        };
    };

    const joinChannel = async (channel: string, serverId: string, config: StreamConfig = { video: true, audio: true }) => {
        console.log("Joining channel", channel, serverId, isConnected);
        if (!serverId) {
            throw new Error("Server ID must be specified.")
        }

        if (!channel) {
            throw new Error("Channel must be specified.")
        }

        await createPeerConnection(channel, serverId);

        const stream = await navigator.mediaDevices.getUserMedia(config);
        setLocalStream(stream);

        const addTrackPromises = stream.getTracks().map(track => {
            return pcRef.current!.addTrack(track, stream);
        });

        await Promise.all(addTrackPromises);

        // Before initializing the call, we need to remove the stream from the remoteStreams
        setRemoteStreams([]);

        console.log("Initialized call");
        sendWebRTCMessage('initializeCall', channel, serverId, { streamId: stream?.id });
    };

    const closeChannel = () => {
        if (pcRef.current) pcRef.current.close();
        setLocalStream(null);
        setRemoteStreams([]);
        sendWebRTCMessage('leave', "", "" , {});
    };

    const handleOfferMessage = async (message: Message) => {
        if (Array.isArray(message.content))
            return console.log('Failed to parse offer');

        const offer = message.content.data;
        if (!offer) {
            return console.log('Failed to parse offer');
        }

        await pcRef.current!.setRemoteDescription(JSON.parse(offer));

        const answer = await pcRef.current!.createAnswer();
        await pcRef.current!.setLocalDescription(answer);

        sendWebRTCMessage('answer', message.channel, message.serverId, { answer });
    };

    const handleCandidateMessage = (message: Message) => {
        if (Array.isArray(message.content))
            return console.log('Failed to parse candidate');

        const candidate = message.content.data;
        if (!candidate) {
            return console.log('Failed to parse candidate');
        }
        pcRef.current!.addIceCandidate(JSON.parse(candidate));
    };

    return (
        <WebRTCContext.Provider value={{ localStream, remoteStreams, isConnected, joinChannel, closeChannel }}>
            {children}
        </WebRTCContext.Provider>
    );
};