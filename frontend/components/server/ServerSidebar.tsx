import axios from '@/utils/axios';

import { ChannelType, Server } from '@/types/models';
import { decodeJwtPayload } from '@/lib/utils';

import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';

import { ServerHeader } from './ServerHeader';
import { useRouter } from 'next/navigation';
import { currentProfile } from '@/lib/currentProfile';

interface ServerSidebarProps {
    serverId: string;
}

async function getServerDetails(serverId: string) {
    try {
        const server = await axios.get<{
            message: string,
            server: Server
        }>(`/servers/${serverId}/details`, {
            headers: {
                Cookie: cookies().toString()
            }
        })

        return server.data.server;
    } catch (error) {
        // console.error(error);
        return null;
    }
}

export const ServerSidebar = async ({
    serverId
}: ServerSidebarProps) => {
    const profile = await currentProfile();

    if(!profile) {
        return redirect("/");
    }

    const server = await getServerDetails(serverId);

    const textChannels = server?.channels?.filter((channel) => channel.type === ChannelType.TEXT);
    const audioChannels = server?.channels?.filter((channel) => channel.type === ChannelType.AUDIO);
    const videoChannels = server?.channels?.filter((channel) => channel.type === ChannelType.VIDEO);
    const members = server?.members?.filter((member) => member.profileID !== profile.id);

    if (!server) {       
        return redirect("/");
    }

    const role = server.members?.find((member) => member.profileID === profile.id)?.role;

    return (
        <div className='flex flex-col h-full text-primary w-full dark:bg-[#2B2D31] bg-[#F2F3F5]'>
            {/* Server Sidebar Components */}
            <ServerHeader 
                server={server}
                role={role}
            />
        </div>
    )
}