"use client"
import { useWebSocket } from "@/components/providers/SocketProvider";
import { useQueryClient } from "@tanstack/react-query";
import { useEffect, useRef, useState } from 'react';

interface Message {
    type: string;
    channel: string;
    content: {
        data?: string;
        username?: string;
        streamId?: string;
        imageURL?: string;
    } | {
        data?: string;
        username?: string;
        streamId?: string;
        imageURL?: string;
    }[];
}

interface ChannelConfig {
    channel: string;
    serverId: string;
}


export const useWebRTC = ({ channel, serverId }: ChannelConfig) => {
    const { socket, isConnected, sendWebRTCMessage } = useWebSocket();
    const [localStream, setLocalStream] = useState<MediaStream | null>(null);
    const [remoteStreams, setRemoteStreams] = useState<MediaStream[]>([]);
    const [participant, setParticipant] = useState<Map<string, { data?: string, username: string, streamId: string }>>(new Map());
    const [isReady, setIsReady] = useState<boolean>(false);

    const pcRef = useRef<RTCPeerConnection | null>(null);

    const configuration: RTCConfiguration = {
        iceServers: [
            { urls: 'stun:stun.l.google.com:19302' }
        ]
    };

    useEffect(() => {
        return () => {
            if (pcRef.current) pcRef.current.close();
        };
    }, []);

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
                case 'participant':
                    handleParticipantMessage(message);
                    break;
                default:
                    console.log('Unknown message type:', message.type);
            }
        };

        if (isConnected) {
            joinChannel();

            socket!.addEventListener('message', handleMessage);

            return () => {
                socket!.removeEventListener('message', handleMessage);
            };
        }
    }, [isConnected]);

    useEffect(() => {
        if (isReady && isConnected) {
            console.log("Send initializeCall");
            sendWebRTCMessage('initializeCall', channel, { streamId: localStream?.id, serverId: serverId });
        }
    }, [isReady, isConnected])

    const createPeerConnection = async () => {
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
                sendWebRTCMessage('candidate', channel, { candidate: event.candidate });
            }
        };
    };

    const joinChannel = async () => {
        if (!channel) {
            throw new Error("Channel must be specified.")
        }

        await createPeerConnection();

        const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
        setLocalStream(stream);

        stream.getTracks().forEach(track => {
            pcRef.current!.addTrack(track, stream)
            setIsReady(true);
        });
    };

    const handleParticipantMessage = async (message: Message) => {
        if (!Array.isArray(message.content))
            return console.log('Failed to parse participant');

        const participantMap = message.content.reduce((participant, cur) => {
            if (!participant.has(cur.streamId)) participant.set(cur.streamId, cur);
            return participant;
        }, new Map())
        setParticipant(participantMap)
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

        sendWebRTCMessage('answer', channel, { answer });
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

    return { joinChannel, localStream, remoteStreams, participant, isConnected, peerConnection: pcRef }
};