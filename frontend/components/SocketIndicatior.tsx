"use client";

import { useWebSocket } from '@/components/providers/SocketProvider';
import { Badge } from '@/components/ui/Badge';

export const SocketIndicator = () => {
    const { isConnected } = useWebSocket();

    if (!isConnected) {
        return (
            <Badge 
                variant="outline" 
                className="bg-yellow-600 text-white border-none"
            >
                Fallback: Polling every 1s
            </Badge>
        );
    }

    return (
        <Badge 
            variant="outline" 
            className="bg-emerald-600 text-white border-none"
        >
            Live: Real-time updates
        </Badge>
    );
}