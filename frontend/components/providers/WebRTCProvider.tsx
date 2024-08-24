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

type WebRTCContextType = {
    localStream: MediaStream | null;
    remoteStreams: MediaStream[];
    isConnected: boolean;
    joinChannel: (channel: string, serverId: string) => void;
    closeChannel: () => void;
};

const WebRTCContext = createContext<WebRTCContextType>({
    localStream: null,
    remoteStreams: [],
    isConnected: false,
    joinChannel: () => {},
    closeChannel: () => {}
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
    const [isReady, setIsReady] = useState<boolean>(false);
    const [channelConfig, setChannelConfig] = useState<ChannelConfig | null>(null);

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

    useEffect(() => {
        if (isReady && isConnected && channelConfig) {
            console.log("Send initializeCall");
            sendWebRTCMessage('initializeCall', channelConfig.channel, channelConfig.serverId, { streamId: localStream?.id });
        }
    }, [isReady, isConnected, channelConfig])


    const createPeerConnection = async (channel: string, serverId: string) => {
        pcRef.current = new RTCPeerConnection(configuration);

        pcRef.current.ontrack = (event) => {
            console.log("ontrack ", event);
            if (event.track.kind === 'video') {
                setRemoteStreams(prevStreams => [...prevStreams, event.streams[0]]);
            }

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

    const joinChannel = async (channel: string, serverId: string) => {
        console.log("Joining channel", channel, serverId, isConnected);
        if (!serverId) {
            throw new Error("Server ID must be specified.")
        }
        
        if (!channel) {
            throw new Error("Channel must be specified.")
        }

        setChannelConfig({ channel, serverId });
        
        await createPeerConnection(channel, serverId);

        const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
        setLocalStream(stream);

        stream.getTracks().forEach(track => {
            pcRef.current!.addTrack(track, stream);
            setIsReady(true);
        });
    };

    const closeChannel = async () => {
        if (pcRef.current) pcRef.current.close();
        setLocalStream(null);
        setRemoteStreams([]);
        setIsReady(false);
        setChannelConfig(null);
        sendWebRTCMessage('leave', channelConfig!.channel, channelConfig!.serverId, {});
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

        if (channelConfig?.channel &&  channelConfig?.serverId) {
            sendWebRTCMessage('answer', channelConfig.channel, channelConfig.serverId, { answer });
        }
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