import axios from "@/utils/axios";

import { currentProfile } from "@/lib/currentProfile";
import { redirect } from "next/navigation";
import { cookies } from 'next/headers';
import { Member } from "@/types/models";
import { getOrCreateConversation } from "@/lib/conversation";
import { ChatHeader } from "@/components/chat/ChatHeader";
import { ChatMessages } from "@/components/chat/ChatMessages";
import { ChatInput } from "@/components/chat/ChatInput";

const getMember = async (serverId: string) => {
    try {
        const response = await axios.get<{
            message: string;
            member: Member;
        }>(`/members/servers/${serverId}`, {
            headers: {
                Cookie: cookies().toString()
            }
        })

        return response.data.member;
    } catch (err) {
        return null;
    }
}

interface MemberIdPageProps {
    params: {
        memberId: string;
        serverId: string;
    }
}

const MemberIdPage = async ({
    params
}: MemberIdPageProps) => {
    const profile = await currentProfile();

    if (!profile) {
        return redirect("/sign-in");
    }

    const currentMember = await getMember(params.serverId);

    if (!currentMember) {
        return redirect("/")
    }

    const conversation = await getOrCreateConversation(currentMember.id, params.memberId);

    if (!conversation) {
        return redirect(`/servers/${params.serverId}`)
    }

    const { memberOne, memberTwo } = conversation;

    const otherMember = memberOne.profileID === profile.id ? memberTwo : memberOne;
 
    return (
        <div className="bg-white dark:bg-[#313338] flex flex-col h-full">
            <ChatHeader 
                imageUrl={otherMember?.profile?.imageUrl}
                name={otherMember?.profile?.name ?? ""}
                serverId={params.serverId}
                type="conversation"
            />
            <ChatMessages 
                member={currentMember}
                name={otherMember.profile.name}
                chatId={conversation.id}
                type="conversation"
                apiUrl="/direct-messages"
                paramKey="conversationId"
                paramValue={conversation.id}
                socketUrl="/ws/direct-messages"
                socketQuery={{
                    conversationId: conversation.id,
                }}
            />
            <ChatInput 
                name={otherMember.profile.name}
                type="conversation"
                apiUrl="/ws/direct-messages"
                query={{
                    conversationId: conversation.id
                }}
            />
        </div>
    );
}

export default MemberIdPage;