"use client";

import { useWebRTC } from "@/components/providers/WebRTCProvider";
// import { useWebRTC } from "@/hooks/useWebRTC";
import { MediaPanel } from "@/components/chat/video/MediaPanel";
import { cn } from "@/lib/utils";
import { useAuth } from "@/components/providers/AuthProvider";
import { useEffect, useRef } from "react";
import { useParticipant } from "@/components/providers/ParticipantProvider";

interface StreamConfig {
    audio: boolean;
    video: boolean;
}

interface VideoConferenceProps {
    chatId: string;
    serverId: string;
    streamConfig: StreamConfig;
}

const getGridClassNames = (participantCount: number) => {
    return cn(
        "grid gap-4 w-full p-5",
        participantCount <= 2
            ? "grid-cols-2 min-w-sm max-w-[80vh]"
            : "md:grid-cols-2 h-full lg:h-auto lg:grid-cols-6"
    );
};

const getMediaPanelClassNames = (participantCount: number) => {
    return participantCount > 2
        ? "lg:last:nth-[3n-1]:col-end-[-2] lg:nth-last-child-[2]:nth-[3n+1]:col-end-4 lg:last:nth-[3n-2]:col-end-5"
        : "";
};

export const VideoConference = ({ chatId, serverId, streamConfig = { audio: true, video: true } }: VideoConferenceProps) => {
    const { authState } = useAuth();
    const { localStream, remoteStreams, isConnected, joinChannel } = useWebRTC();
    const { participant } = useParticipant();

    const isMounted = useRef(false);

    useEffect(() => {
        if (!isMounted.current && isConnected) {
            isMounted.current = true;
            joinChannel(chatId, serverId, streamConfig);
        }
    }, [isConnected]);

    const participantCount = remoteStreams.length + 1;
    const gridClassNames = getGridClassNames(participantCount);
    const mediaPanelClassNames = getMediaPanelClassNames(participantCount);

    return (
        <div className="lg:h-full flex items-center justify-center overflow-y-auto bg-[#000000]">
            <div className={gridClassNames}>
                <MediaPanel
                    media={localStream}
                    name={authState.profile?.name}
                    className={mediaPanelClassNames}
                    isRemote={false}
                />
                {remoteStreams.map((remoteStream) => (
                    <MediaPanel
                        key={remoteStream.id}
                        media={remoteStream}
                        name={participant[chatId]?.find((p) => p.streamId === remoteStream.id)?.username}
                        className={mediaPanelClassNames}
                        isRemote={true}
                    />
                ))}
            </div>
        </div>
    );
};