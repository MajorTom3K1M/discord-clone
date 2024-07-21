"use client";

import { useWebRTC } from "@/hooks/useWebRTC";
import { MediaPanel } from "@/components/chat/video/MediaPanel";
import { Card } from "@/components/ui/Card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/Avatar";
import { useState } from "react";
import { cn } from "@/lib/utils";

interface VideoConferenceProps {
    chatId: string;
}

export const VideoConference = ({
    chatId
}: VideoConferenceProps) => {
    const participant: any[] = [1]
    const { joinChannel, localStream, remoteStreams, peerConnection, isConnected } = useWebRTC({ channel: chatId });
    const participantCount = remoteStreams.length + 1;

    return (
        <div className="lg:h-full flex items-center justify-center overflow-y-auto bg-[#000000]">
            <div className={cn(
                "grid gap-4 w-full p-5",
                participantCount <= 2 ? "grid-cols-2 min-w-sm max-w-[80vh]" : "md:grid-cols-2 h-full lg:h-auto lg:grid-cols-6"
            )}>
                <MediaPanel
                    media={localStream}
                    className={
                        participantCount > 2 ?
                            "lg:last:nth-[3n-1]:col-end-[-2] lg:nth-last-child-[2]:nth-[3n+1]:col-end-4 lg:last:nth-[3n-2]:col-end-5" :
                            ""
                    } />
                {remoteStreams.map((media, index) => (
                    <MediaPanel
                        key={index} media={media}
                        className={
                            participantCount > 2 ?
                                "lg:last:nth-[3n-1]:col-end-[-2] lg:nth-last-child-[2]:nth-[3n+1]:col-end-4 lg:last:nth-[3n-2]:col-end-5" :
                                ""
                        } />
                ))}
            </div>
        </div>

    )
};