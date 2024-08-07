"use client";

import { useInView } from 'react-intersection-observer';
import { Loader2, ServerCrash } from "lucide-react";
import { Fragment, useRef, ElementRef, useEffect } from "react";
import { format } from 'date-fns';

import { Member, Message } from "@/types/models";
import { ChatWelcome } from "@/components/chat/ChatWelcome";
import { ChatItem } from "@/components/chat/ChatItem";
import { useChatQuery } from "@/hooks/useChatQuery";
import { useChatSocket } from "@/hooks/useChatSocket";
import { useChatScroll } from "@/hooks/useChatScroll";

const DATE_FORMAT = "d MMM yyyy, HH:mm"

interface ChatMessagesProps {
    name: string;
    member: Member;
    chatId: string;
    apiUrl: string;
    socketUrl: string;
    socketQuery: Record<string, string>;
    paramKey: "channelId" | "conversationId";
    paramValue: string;
    type: "channel" | "conversation";
}

export const ChatMessages = ({
    name,
    member,
    chatId,
    apiUrl,
    socketUrl,
    socketQuery,
    paramKey,
    paramValue,
    type
}: ChatMessagesProps) => {
    const queryKey = `chat:${chatId}`;
    const addKey = `chat:${chatId}:messages`;
    const updateKey = `chat:${chatId}:messages:update`;

    const chatRef = useRef<ElementRef<"div">>(null);
    const bottomRef = useRef<ElementRef<"div">>(null);

    const {
        data,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
        status
    } = useChatQuery({
        queryKey,
        apiUrl,
        paramKey,
        paramValue
    });
    const { ref: loadMoreRef, inView } = useInView({
        rootMargin: '100px',
        threshold: 0
    });
    useChatSocket({ queryKey, addKey, updateKey });
    useChatScroll({
        chatRef,
        bottomRef,
        inView,
        loadMore: fetchNextPage,
        shouldLoadMore: !isFetchingNextPage && !!hasNextPage,
        count: data?.pages?.[0]?.items?.length ?? 0
    })

    if (status === "pending") {
        return (
            <div className="flex flex-col flex-1 justify-center items-center">
                <Loader2 className="h-7 w-7 text-zinc-500 animate-spin my-4" />
                <p className="text-xs text-zinc-500 dark:text-zinc-400">
                    Loading messages...
                </p>
            </div>
        )
    }

    if (status === "error") {
        return (
            <div className="flex flex-col flex-1 justify-center items-center">
                <ServerCrash className="h-7 w-7 text-zinc-500 my-4" />
                <p className="text-xs text-zinc-500 dark:text-zinc-400">
                    Something went wrong!
                </p>
            </div>
        )
    }

    return (
        <div className="flex-1 flex flex-col-reverse py-4 overflow-y-auto">
            <div ref={bottomRef} />
            <div className="flex flex-col-reverse mt-auto">
                {data?.pages?.map((group, i) => (
                    <Fragment key={i}>
                        {group.items.map((message: Message) => (
                            <ChatItem
                                key={message.id}
                                id={message.id}
                                currentMember={member}
                                member={message.member}
                                content={message.content}
                                fileUrl={message.fileUrl}
                                deleted={message.deleted}
                                timestamp={format(new Date(message.created_at), DATE_FORMAT)}
                                isUpdated={message.updated_at !== message.created_at}
                                socketUrl={socketUrl}
                                socketQuery={socketQuery}
                            />
                        ))}
                    </Fragment>
                ))}
            </div>
            {hasNextPage && (
                <div className="flex justify-center">
                    {isFetchingNextPage ? (
                        <Loader2 className="h-6 w-6 text-zinc-500 animate-spin my-4" />
                    ) : (
                        <button
                            ref={loadMoreRef}
                            onClick={() => fetchNextPage()}
                            className="text-zinc-500 hover:text-zinc-600 dark:text-zinc-400 text-xs
                            my-4 dark:hover:text-zinc-300 transition"
                        >
                            Load previous messages
                        </button>
                    )}
                </div>
            )}
            {!hasNextPage && (
                <ChatWelcome
                    type={type}
                    name={name}
                />
            )}
            {!hasNextPage && <div className="flex-1" />}
            <div ref={chatRef} />
        </div>
    )
};