import axios from "@/utils/axios";

import { currentProfile } from "@/lib/currentProfile";
import { Channel, ChannelType, Member } from "@/types/models";
import { redirect } from "next/navigation";
import { cookies } from "next/headers";

import { ChatHeader } from "@/components/chat/ChatHeader";
import { ChatInput } from "@/components/chat/ChatInput";
import { ChatMessages } from "@/components/chat/ChatMessages";
import { VideoConference } from "@/components/chat/video/VideoConference";

interface ChannelIdPageProps {
    params: {
        serverId: string;
        channelId: string;
    }
}

const getChannel = async (channelId: string) => {
    try {
        const response = await axios.get<{
            message: string;
            channel: Channel;
        }>(`/channels/${channelId}`, {
            headers: {
                Cookie: cookies().toString()
            }
        });

        return response.data.channel;
    } catch (err) {
        return null
    }
};

const getMember = async (serverId: string) => {
    try {
        const response = await axios.get<{
            message: string;
            member: Member;
        }>(`/servers/${serverId}/members`, {
            headers: {
                Cookie: cookies().toString()
            }
        });

        return response.data.member;
    } catch (err) {
        return null
    }
};

const ChannelIdPage = async ({
    params
}: ChannelIdPageProps) => {
    const profile = await currentProfile();

    if (!profile) {
        return redirect("/sign-in")
    }

    const channel = await getChannel(params.channelId);

    const member = await getMember(params.serverId);

    if (!channel || !member) {
        return redirect("/");
    }

    return (
        <div className="bg-white dark:bg-[#313338] flex flex-col h-full">
            <ChatHeader
                name={channel.name}
                serverId={channel.serverID}
                type="channel"
            />
            {channel.type === ChannelType.TEXT ? (
                <>
                    <ChatMessages
                        member={member}
                        name={channel.name}
                        chatId={channel.id}
                        type="channel"
                        apiUrl="/messages"
                        socketUrl="/ws/messages"
                        socketQuery={{
                            channelId: channel.id,
                            serverId: channel.serverID
                        }}
                        paramKey="channelId"
                        paramValue={channel.id}
                    />
                    <ChatInput
                        name={channel.name}
                        type="channel"
                        apiUrl="/ws/messages"
                        query={{
                            channelId: channel.id,
                            serverId: channel.serverID
                        }}
                    />
                </>
            ) : null}

            {channel.type === ChannelType.AUDIO ? (
                <VideoConference chatId={channel.id} serverId={channel.serverID} streamConfig={{video: false, audio: true}} />
            ) : null}

            {channel.type === ChannelType.VIDEO ? (
                <VideoConference chatId={channel.id} serverId={channel.serverID} streamConfig={{video: true, audio: true}} />
            ) : null}
        </div>
    );
}

export default ChannelIdPage;