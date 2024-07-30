"use client";

import React, { createContext, use, useContext, useEffect, useState } from 'react';
import { useWebSocket } from './SocketProvider';

type Participant = {
    data: string;
    username: string;
    streamId: string;
    imageURL: string;
    clientId: string;
};

interface ParticipantsMessage {
    [channelId: string]: {
        [streamId: string]: Participant;
    }
}

type ParticipantContextType = {
    participant: { [channelId: string]: Participant[] };
    setSingleParticipant: (channelId: string, participant: Participant) => void;
    setMultipleParticipants: (participants: { [channelId: string]: { [streamId: string]: Participant } }) => void;
};

const ParticipantContext = createContext<ParticipantContextType>({
    participant: {},
    setSingleParticipant: (channelId: string, participant: Participant) => { },
    setMultipleParticipants: (participants: { [channelId: string]: { [streamId: string]: Participant } }) => { },
});

export const useParticipant = () => {
    const context = useContext(ParticipantContext);
    if (!context) {
        throw new Error('useParticipant must be used within a ParticipantProvider');
    }
    return context;
};

export const ParticipantProvider = ({
    children
}: {
    children: React.ReactNode
}) => {
    const [participant, setParticipant] = useState<{ [channelId: string]: Participant[] }>({});

    const { socket, isConnected } = useWebSocket();

    useEffect(() => {
        const handleMessage = (event: MessageEvent) => {
            const message = JSON.parse(event.data) as { type: string, channel: string, content: any };
            if (message.type === 'participant') {
                const channelId = message.channel;
                setSingleParticipant(channelId, message.content);
            } else if (message.type === 'participants') {
                const participants = message.content as ParticipantsMessage;
                setMultipleParticipants(participants);
            }
        };

        if (isConnected) {
            socket!.addEventListener('message', handleMessage);
            return () => {
                socket!.removeEventListener('message', handleMessage);
            };
        }
    }, [isConnected, socket]);

    const setSingleParticipant = (channelId: string, participant: Participant) => {
        setParticipant((prevParticipant) => {
            let newParticipant = prevParticipant[channelId] ?? [];

            newParticipant = [...newParticipant.filter((p) => p.clientId !== participant.clientId), participant];

            return {
                ...prevParticipant,
                [channelId]: newParticipant
            };
        });
    }

    const setMultipleParticipants = (participants: { [channelId: string]: { [streamId: string]: Participant } }) => {
        const newParticipants = Object.keys(participants).reduce((acc, channelId) => {
            acc[channelId] = Object.values(participants[channelId]);
            return acc;
        }, {} as { [channelId: string]: Participant[] });
        setParticipant(newParticipants);
    };

    return (
        <ParticipantContext.Provider value={{ 
            participant: participant, setSingleParticipant, setMultipleParticipants 
        }}>
            {children}
        </ParticipantContext.Provider>
    );
};