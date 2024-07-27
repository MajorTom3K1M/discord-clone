import { useWebSocket } from "@/components/providers/SocketProvider";
import { channel } from "diagnostics_channel";
import { useEffect, useRef, useState } from "react";

type ParticipantSocketProps = {
    channelId?: string;
    serverId: string;
}

interface Participant {
    data: string;
    username: string;
    streamId: string;
    imageURL: string;
    clientId: string;
}

interface ParticipantsMessage {
    [channelId: string]: {
        [streamId: string]: Participant;
    }
}

export const useParticipantSocket = ({ channelId, serverId }: ParticipantSocketProps) => {
    const { socket, isConnected, leaveServer } = useWebSocket();
    const [participant, setParticipant] = useState<{ [channelId: string]: Participant[] }>({});

    useEffect(() => {
        console.log("participant", participant);
    }, [participant])

    useEffect(() => {
        const handleMessage = (event: MessageEvent) => {
            const message = JSON.parse(event.data) as { type: string, channel: string, content: any };
            if (message.type === 'participant') {
                console.log("message : ", message);
                const channelId = message.channel;
                const participant: Participant = message.content;
                setParticipant((prevParticipant) => {
                    let newParticipant = prevParticipant[channelId] ?? [];
                    const isParticipantExist = (p: Participant) => String(p.clientId) === String(participant.clientId);

                    if (newParticipant?.some(isParticipantExist)) {
                        if (participant.data === 'left') {
                            return {
                                ...prevParticipant,
                                [channelId]: newParticipant.filter((p) => p.clientId !== participant.clientId)
                            };
                        }

                    }

                    newParticipant = [...newParticipant.filter((p) => p.clientId !== participant.clientId), participant];

                    return {
                        ...prevParticipant,
                        [channelId]: newParticipant
                    };
                });
            } else if (message.type === 'participants') {
                const participantsMessage: ParticipantsMessage = message.content;
                const participannts = Object.entries(participantsMessage).reduce((acc, [channelId, participants]) => {
                    acc[channelId] = Object.values(participants);
                    return acc;
                }, {} as { [channelId: string]: Participant[] });
                setParticipant(participannts);
            }
        }

        if (isConnected) {
            socket!.addEventListener('message', handleMessage);

            return () => {
                socket!.removeEventListener('message', handleMessage);
            };
        }
    }, [isConnected, socket, channelId, serverId]);

    return { participant, isConnected, leaveServer };
}