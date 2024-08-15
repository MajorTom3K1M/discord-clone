import axios from '@/utils/axios';

import { ChannelType, MemberRole, Server } from '@/types/models';
import { decodeJwtPayload } from '@/lib/utils';

import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';

import { Hash, Mic, ShieldAlert, ShieldCheck, Video, Settings } from 'lucide-react';

import { currentProfile } from '@/lib/currentProfile';
import { ScrollArea } from '@/components/ui/ScrollArea';
import { Separator } from '@/components/ui/Separator';

import { ServerHeader } from './ServerHeader';
import { ServerSearch } from './ServerSearch';
import { ServerSection } from './ServerSection';
import { ServerChannel } from './ServerChannel';
import { ServerMember } from './ServerMember';
import { UserAvatar } from '../UserAvatar';
import { ServerChannelMember } from './ServerChannelMember';
import { useWebRTC } from '@/hooks/useWebRTC';
import { ServerFooter } from './ServerFooter';

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
                        <div className="space-y-[2px]">
                            {textChannels.map((channel) => (
                                <ServerChannel
                                    key={channel.id}
                                    channel={channel}
                                    role={role}
                                    server={server}
                                />
                            ))}
                        </div>
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
                        <div className="space-y-[2px]">
                            {audioChannels.map((channel) => (
                                <div className='flex flex-col' key={channel.id}>
                                    <ServerChannel
                                        key={channel.id}
                                        channel={channel}
                                        role={role}
                                        server={server}
                                    />
                                    {/* <ServerChannelMember 
                                        className='pl-8 pb-2'
                                        channel={channel}
                                        serverId={serverId}
                                        // src='https://cdn.discordapp.com/avatars/293021833350086656/8b125db71d831bd62fb80818bc574777.webp'
                                        // name={'MajorTom'}
                                    /> */}
                                    {/* <div className='pl-8 pb-2 text-zinc-500 dark:text-zinc-400 
                                        hover:text-zinc-600 dark:hover:text-zinc-300 transition
                                        text-[14px]'
                                    >
                                        <div className='relative flex items-center cursor-pointer flex-1 
                                            hover:bg-zinc-700/10 dark:hover:bg-zinc-700/50 transition rounded-md
                                            pt-1 pb-1'>
                                            <UserAvatar
                                                className='h-6 w-6 md:h-6 md:w-6 ml-2 mr-2'
                                                src='https://cdn.discordapp.com/avatars/293021833350086656/8b125db71d831bd62fb80818bc574777.webp'
                                            />
                                            MajorTom
                                        </div>
                                    </div> */}
                                </div>
                            ))}
                        </div>
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
                        <div className="space-y-[2px]">
                            {videoChannels.map((channel) => (
                                <div className='flex flex-col' key={channel.id}>
                                    <ServerChannel
                                        key={channel.id}
                                        channel={channel}
                                        role={role}
                                        server={server}
                                    />
                                    <ServerChannelMember
                                        className='pl-8 pb-2'
                                        channel={channel}
                                        serverId={serverId}
                                    />
                                </div>
                            ))}
                        </div>
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
            <ServerFooter
                role={role}
            />
            {/* <div className="flex items-center p-2 dark:bg-[#222327] bg-[#E3E5E8] drop-shadow-md justify-between">
                <div className='flex items-center'>
                    <UserAvatar
                        src={"https://images.pexels.com/photos/220453/pexels-photo-220453.jpeg?auto=compress&cs=tinysrgb&dpr=2&w=80"}
                        className="h-8 w-8 md:h-8 md:w-8 ml-1 mr-1"
                    />
                    <div className="text-xs">
                        <h3 className='font-bold'>Matt Cooper</h3>
                        <p className='font-thin'>Moderator</p>
                    </div>
                </div>
                <Settings className='mr-2 h-5 w-5' />
            </div> */}
        </div>
    )
}