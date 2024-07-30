"use client"
import Image from "next/image";

import axios from "@/utils/axios";

import { useParams, useRouter } from "next/navigation";

import { cn } from "@/lib/utils";
import { ActionTooltip } from "@/components/ActionTooltip";
import { useEffect, useRef } from "react";
import { useWebSocket } from "../providers/SocketProvider";

interface NavigationItemProps {
    id: string;
    imageUrl: string;
    name: string;
}

export const NavigationItem = ({
    id,
    imageUrl,
    name
}: NavigationItemProps) => {
    const { isConnected, joinedServer } = useWebSocket();
    const params = useParams();
    const router = useRouter();
    const isFirstTime = useRef(true);

    const getParticipants = async (serverId: string) => {
        try {
            console.log("Joining server : ", serverId);
            joinedServer(serverId);
        } catch (err) {
            console.error(err);
        }
    }

    useEffect(() => {
        if (isConnected && params?.serverId === id && isFirstTime.current) {
            getParticipants(id);
            isFirstTime.current = false;
        }
    }, [isConnected]);

    const onClick = () => {
        getParticipants(id);
        router.push(`/servers/${id}`);
    }

    return (  
        <ActionTooltip
            side="right"
            align="center"
            label={name}
        >
            <button
                onClick={onClick}
                className="group relative flex items-center"
            >
                <div className={cn(
                    "absolute left-0 bg-primary rounded-r-full transition-all w-[4px]",
                    params?.serverId !== id && "group-hover:h-[20px]",
                    params?.serverId === id ? "h-[36px]" : "h-[8px]"
                )} />
                <div className={cn(
                    "relative group flex mx-3 h-[48px] w-[48px] rounded-[24px] group-hover:rounded-[16px] transition-all overflow-hidden",
                    params?.serverId === id && "bg-primary/10 text-primary rounded-[16px]"
                )}>
                    <Image 
                        fill
                        src={imageUrl}
                        alt="Channel"
                    />
                </div>
            </button>
        </ActionTooltip>
    );
}