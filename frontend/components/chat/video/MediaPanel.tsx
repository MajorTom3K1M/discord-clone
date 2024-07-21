"use client";

import { cn } from '@/lib/utils';

interface MediaPanelProps {
    media?: MediaStream | null;
    className?: string;
}

export const MediaPanel = ({ media, className }: MediaPanelProps) => {
    return (
        <div className={cn("relative w-full col-span-2", className)}>
            <div className="pb-[56.25%]"></div>
            <div className="absolute top-0 left-0 w-full h-full shadow-sm rounded-lg flex items-center justify-center">
                <video
                    className='rounded-lg h-full w-full dark:bg-[#1E1F22] bg-[#E3E5E8]'
                    autoPlay
                    muted
                    ref={(video) => {
                        if (video && media) {
                            video.srcObject = media;
                        }
                    }}
                />
            </div>
        </div>
    );
};