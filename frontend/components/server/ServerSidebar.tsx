import axios from '@/utils/axios';

import { ChannelType, MemberRole, Server } from '@/types/models';
import { decodeJwtPayload } from '@/lib/utils';

import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';

import { Hash, Mic, ShieldAlert, ShieldCheck, Video } from 'lucide-react';

import { currentProfile } from '@/lib/currentProfile';
import { ScrollArea } from '@/components/ui/ScrollArea';
import { Separator } from '@/components/ui/Separator';

import { ServerHeader } from './ServerHeader';
import { ServerSearch } from './ServerSearch';
import { ServerSection } from './ServerSection';
import { ServerChannel } from './ServerChannel';
import { ServerMember } from './ServerMember';

interface ServerSidebarProps {
    serverId: string;
}

const iconMap = {
    [ChannelType.TEXT]: <Hash className='mr-2 h-4 w-4' />,
    [ChannelType.AUDIO]: <Mic className='mr-2 h-4 w-4' />,
    [ChannelType.VIDEO]: <Video className='mr-2 h-4 w-4' />
};

const roleIconMap = {
    [MemberRole.GUEST]: null,
    [MemberRole.MODERATOR]: <ShieldCheck className='h-4 w-4 mr-2 text-indigo-500' />,
    [MemberRole.ADMIN]: <ShieldAlert className='h-4 w-4 mr-2 text-rose-500' />,
};

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

    if (!profile) {
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
            <ScrollArea className='flex-1 px-3'>
                <div className='mt-2'>
                    <ServerSearch
                        data={[
                            {
                                label: "Text Channels",
                                type: "channel",
                                data: textChannels?.map((channel) => ({
                                    id: channel.id,
                                    name: channel.name,
                                    icon: iconMap[channel.type],
                                }))
                            },
                            {
                                label: "Voice Channels",
                                type: "channel",
                                data: audioChannels?.map((channel) => ({
                                    id: channel.id,
                                    name: channel.name,
                                    icon: iconMap[channel.type],
                                }))
                            },
                            {
                                label: "Video Channels",
                                type: "channel",
                                data: videoChannels?.map((channel) => ({
                                    id: channel.id,
                                    name: channel.name,
                                    icon: iconMap[channel.type],
                                }))
                            },
                            {
                                label: "Members",
                                type: "member",
                                data: members?.map((member) => ({
                                    id: member.id,
                                    name: member.profile?.name,
                                    icon: roleIconMap[member.role],
                                }))
                            },
                        ]}
                    />
                </div>
                <Separator className='bg-zinc-200 dark:bg-zinc-700 rounded-md my-2' />
                {!!textChannels?.length && (
                    <div className='mb-2'>
                        <ServerSection
                            sectionType="channels"
                            channelType={ChannelType.TEXT}
                            role={role}
                            label="Text Channels"
                        />
                        {textChannels.map((channel) => (
                            <ServerChannel
                                key={channel.id}
                                channel={channel}
                                role={role}
                                server={server}
                            />
                        ))}
                    </div>
                )}
                {!!audioChannels?.length && (
                    <div className='mb-2'>
                        <ServerSection
                            sectionType="channels"
                            channelType={ChannelType.AUDIO}
                            role={role}
                            label="Voice Channels"
                        />
                        {audioChannels.map((channel) => (
                            <ServerChannel
                                key={channel.id}
                                channel={channel}
                                role={role}
                                server={server}
                            />
                        ))}
                    </div>
                )}
                {!!videoChannels?.length && (
                    <div className='mb-2'>
                        <ServerSection
                            sectionType="channels"
                            channelType={ChannelType.VIDEO}
                            role={role}
                            label="Video Channels"
                        />
                        {videoChannels.map((channel) => (
                            <ServerChannel
                                key={channel.id}
                                channel={channel}
                                role={role}
                                server={server}
                            />
                        ))}
                    </div>
                )}
                {!!members?.length && (
                    <div className="mb-2">
                        <ServerSection
                            sectionType="members"
                            role={role}
                            label="Members"
                            server={server}
                        />
                        <div className="space-y-[2px]">
                            {members.map((member) => (
                                <ServerMember
                                    key={member.id}
                                    member={member}
                                    server={server}
                                />
                            ))}
                        </div>
                    </div>
                )}
            </ScrollArea>
        </div>
    )
}