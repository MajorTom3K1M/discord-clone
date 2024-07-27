"use client";

import { cn } from "@/lib/utils";
import { UserAvatar } from "../UserAvatar";
import { Channel } from "@/types/models";
import { useWebRTC } from "@/hooks/useWebRTC";
import { useWebSocket } from "../providers/SocketProvider";
import { useEffect, useState } from "react";
import { useParticipantSocket } from "@/hooks/useParticipantSocket";

interface ServerChannelMemberProps {
    // src: string;
    // name: string;
    // className?: string;
    channel: Channel;
    serverId: string;
    className?: string;
}

interface Users {
    imageURL: string;
    streamId: string;
    username: string;
}

export const ServerChannelMember = ({
    channel,
    serverId,
    className
}: ServerChannelMemberProps) => {
    const { participant, leaveServer } = useParticipantSocket({ channelId: channel.id, serverId });

    useEffect(() => {
        return () => {
            // leaveServer(serverId, channel.id);
            console.log("leave server : ", channel.name);
        }
    }, []);

    return (
        <div key={channel.id}>
            {
                participant[channel.id]?.map((user) => (
                    <div key={user.streamId} className={cn(
                        `text-zinc-500 dark:text-zinc-400 hover:text-zinc-600 
                        dark:hover:text-zinc-300 transition text-[14px]`,
                        className
                    )}>
                        <div className='relative flex items-center cursor-pointer flex-1 
                            hover:bg-zinc-700/10 dark:hover:bg-zinc-700/50 transition 
                            rounded-md pt-1 pb-1'>
                            <UserAvatar
                                className="h-6 w-6 md:h-6 md:w-6 ml-2 mr-2"
                                src={user.imageURL}
                            />
                            {user.username}
                        </div>
                    </div>
                ))
            }
        </div>
    );
};