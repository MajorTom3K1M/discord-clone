"use client"
import { Signal, Phone } from 'lucide-react';
import { useWebRTC } from '@/components/providers/WebRTCProvider';
import { redirect, useRouter } from 'next/navigation';

interface ServerActionProps {
    serverId: string;
}

const ServerAction = ({
    serverId
}: ServerActionProps) => {
    const { closeChannel, localStream } = useWebRTC();
    const router = useRouter();

    const handleDisconnect = () => {
        closeChannel();
        if (serverId) router.push(`/servers/${serverId}`);
    }

    return localStream ? (
        <div className="flex w-full items-center p-2 dark:bg-[#222327] bg-[#E3E5E8] justify-between border-b-[1px] dark:border-zinc-700 border-zinc-200">
            <div className="flex w-full items-center gap-2">
                <Signal size={16} className="text-emerald-600" />
                <div>
                    <p className="text-sm font-semibold text-emerald-600">Connected</p>
                </div>
                <div className="flex-grow"></div>
                <button 
                    className='ml-2 mr-2 items-center cursor-pointer' 
                    onClick={handleDisconnect}
                >
                    <Phone size={16} className="text-white" />
                </button>
            </div>
        </div>
    ) : null;
}

export default ServerAction;