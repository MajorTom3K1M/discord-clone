import axios from "@/utils/axios";

import { currentProfile } from "@/lib/currentProfile";
import { Channel, Member } from "@/types/models";
import { redirect } from "next/navigation";
import { cookies } from "next/headers";

import { ChatHeader } from "@/components/chat/ChatHeader";

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
        </div>
    );
}

export default ChannelIdPage;