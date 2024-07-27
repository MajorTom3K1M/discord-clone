"use client";

import React, { createContext, useContext, useEffect, useState, useRef } from 'react';
import { useAuth } from './AuthProvider';
import { Message, WebRTCMessage } from '@/types/models';
import { channel } from 'diagnostics_channel';

type WebSocketContextType = {
    send: (event: string, data: any) => void;
    joinedServer: (serverId: string) => void;
    leaveServer: (serverId: string, channelId: string) => void;
    sendWebRTCMessage: (type: string, channel: string, serverId: string, message: WebRTCMessage) => void,
    on: (event: string, handler: (data: any) => void) => void;
    off: (event: string) => void;
    socket: WebSocket | null;
    isConnected: boolean;
};

const WebSocketContext = createContext<WebSocketContextType>({
    send: (event: string, data: any) => { },
    joinedServer: (serverId: string) => { },
    leaveServer: (serverId: string, channelId: string) => { },
    sendWebRTCMessage: (type: string, channel: string, serverId: string, message: WebRTCMessage) => {},
    on: (event: string, handler: (data: any) => void) => { },
    off: (event: string) => { },
    socket: null,
    isConnected: false
});

export const useWebSocket = () => {
    const context = useContext(WebSocketContext);
    if (!context) {
        throw new Error('useWebSocket must be used within a WebSocketProvider');
    }
    return context;
};

export const WebSocketProvider = ({
    children
}: {
    children: React.ReactNode
}) => {
    const listeners = useRef<Map<string, Set<Function>>>(new Map());
    const subscribeQueue = useRef<Set<string>>(new Set());
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [isConnected, setIsConnected] = useState(false);
    
    const { authState } = useAuth();

    useEffect(() => {
        if (!window.WebSocket) {
            console.error('WebSocket is not supported by your browser.');
            return;
        }

        if (!authState.profile) {
            return;
        }

        const ws = new WebSocket("ws://localhost:8080/ws");

        ws.onopen = () => {
            console.log('WebSocket connection established.');
            setIsConnected(true);

            subscribeQueue.current.forEach(event => {
                ws.send(JSON.stringify({ type: "subscribe", channel: event }));
            });
            subscribeQueue.current.clear();
        };

        ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log("Received Message : ", message)
            if (message.channel) {
                const handlers = listeners.current.get(message.channel);
                if (handlers) {
                    handlers.forEach(handler => handler(message.content));
                }
            }
        };

        ws.onclose = () => {
            console.log('WebSocket disconnected.');
            setIsConnected(false);
        };

        ws.onerror = (error) => console.error('WebSocket error:', error);

        setSocket(ws);

        return () => {
            ws.close();
        };
    }, [authState.profile]);

    const send = (event: string, data: Message) => {
        if (socket?.readyState === WebSocket.OPEN) {
            console.log({ type: "message", channel: event, content: data })
            socket.send(JSON.stringify({ type: "message", channel: event, content: data }));
        }
    };

    const joinedServer = (serverId: string) => {
        if (socket?.readyState === WebSocket.OPEN) {
            console.log({ type: "joined", serverId: serverId });
            socket.send(JSON.stringify({ type: "joined", serverId }));
        }
    }

    const leaveServer = (serverId: string, channelId: string) => {
        if (socket?.readyState === WebSocket.OPEN) {
            console.log({ type: "leave", serverId: serverId, channelId });
            socket.send(JSON.stringify({ type: "leave", serverId, channel: channelId }));
        }
    }

    const on = (event: string, handler: (data: any) => void) => {
        if (!listeners.current.has(event)) {
            if (socket?.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({ type: "subscribe", channel: event }));
            } else {
                subscribeQueue.current.add(event);
            }
            listeners.current.set(event, new Set());
        }
        listeners.current.get(event)?.add(handler);
    };

    const off = (event: string) => {
        const handlers = listeners.current.get(event);
        if (handlers) {
            listeners.current.delete(event);
            if (socket?.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({ type: "unsubscribe", channel: event }));
            }
        }
    };

    const sendWebRTCMessage = (type: string, channel: string, serverId: string, message: WebRTCMessage) => {
        if (socket?.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ type: type, serverId: serverId, channel: channel, content: message }));
        }
    }

    return (
        <WebSocketContext.Provider value={{ 
            send, on, off, socket, isConnected,  sendWebRTCMessage, joinedServer, leaveServer
        }}>
            {children}
        </WebSocketContext.Provider>
    );
}