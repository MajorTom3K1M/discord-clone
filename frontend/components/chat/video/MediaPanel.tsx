"use client";

import { Badge } from '@/components/ui/Badge';
import { cn } from '@/lib/utils';

interface MediaPanelProps {
    media?: MediaStream | null;
    name?: string;
    className?: string;
}

export const MediaPanel = ({ media, className, name }: MediaPanelProps) => {
    const setVideoRef = (video: HTMLVideoElement | null) => {
        if (video && media) {
            video.srcObject = media;
        }
    };

    return (
        <div className={cn("relative w-full col-span-2", className)}>
            <div className="pb-[56.25%]"></div>
            <div className="absolute top-0 left-0 w-full h-full shadow-sm rounded-lg flex items-center justify-center">
                <video
                    className="rounded-lg h-full w-full dark:bg-[#1E1F22] bg-[#E3E5E8]"
                    autoPlay
                    muted
                    ref={setVideoRef}
                />
                <div className="absolute flex bottom-0 justify-between w-full p-3 h-[55px]">
                    <Badge variant="outline" className="rounded-lg bg-black text-sm text-white border-none">
                        {name}
                    </Badge>
                </div>
            </div>
        </div>
    );
};