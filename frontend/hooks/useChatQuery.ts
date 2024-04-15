import qs from 'query-string';
import axios from '@/utils/axios';
import { useInfiniteQuery } from '@tanstack/react-query';

import { useWebSocket } from '@/components/providers/SocketProvider';
import { Message } from '@/types/models';

interface ChatQueryProps {
    queryKey: string;
    apiUrl: string;
    paramKey: "channelId" | "conversationId";
    paramValue: string;
}

export const useChatQuery = ({
    queryKey,
    apiUrl,
    paramKey,
    paramValue
}: ChatQueryProps) => {
    const { isConnected } = useWebSocket();

    const fetchMessages = async ({ pageParam = undefined }: { pageParam: string | undefined }) => {
        const url = qs.stringifyUrl({
            url: apiUrl,
            query: {
                cursor: pageParam,
                [paramKey]: paramValue
            }
        }, { skipNull: true })

        const res = await axios.get(url);

        return res.data;
    };

    const {
        data,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
        status
    } = useInfiniteQuery({
        queryKey: [queryKey],
        queryFn: ({ pageParam }) => fetchMessages({pageParam}),
        initialPageParam: "",
        getNextPageParam: (lastPage) => {
            if (lastPage?.nextCursor !== "") {
                return lastPage?.nextCursor
            }
        }, 
        refetchInterval: isConnected ? false : 1000
    });

    return { data, fetchNextPage, hasNextPage, isFetchingNextPage, status };
}