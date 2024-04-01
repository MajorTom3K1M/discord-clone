import axios from "@/utils/axios";

import { Conversation } from "@/types/models";
import { cookies } from 'next/headers';

export const getOrCreateConversation = async (memberOneId: string, memberTwoId: string) => {
    try {
        const response = await axios.post<{
            message: string;
            conversation: Conversation
        }>(`/conversations/ensure`, {
            memberOneId,
            memberTwoId
        }, {
            headers: {
                Cookie: cookies().toString()
            }
        });

        return response.data.conversation;
    } catch (err) {
        return null;
    }
};