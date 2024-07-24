"use client";

import { useWebRTC } from "@/hooks/useWebRTC";
import { MediaPanel } from "@/components/chat/video/MediaPanel";
import { cn } from "@/lib/utils";
import { useAuth } from "@/components/providers/AuthProvider";

interface VideoConferenceProps {
    chatId: string;
    serverId: string;
}

const useParticipantName = (participant: Map<string, { username: string }>, streamId: string) => {
    return participant.get(streamId)?.username ?? "";
};

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

export const VideoConference = ({ chatId, serverId }: VideoConferenceProps) => {
    const { authState } = useAuth();
    const { localStream, remoteStreams, participant } = useWebRTC({ channel: chatId, serverId: serverId });

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
                />
                {remoteStreams.map((remoteStream) => (
                    <MediaPanel
                        key={remoteStream.id}
                        media={remoteStream}
                        name={useParticipantName(participant, remoteStream.id)}
                        className={mediaPanelClassNames}
                    />
                ))}
            </div>
        </div>
    );
};