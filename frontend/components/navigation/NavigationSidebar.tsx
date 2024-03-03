import axios from '@/utils/axios';
import { cookies } from 'next/headers';

import { Separator } from '@/components/ui/Separator';
import { ScrollArea } from '@/components/ui/ScrollArea';
import { ModeToggle } from '@/components/ModeToggle';

import { NavigationAction } from './NavigationAction';
import { NavigationItem } from './NavigationItem';

interface Server {
    id: string;
    name: string;
    imageUrl: string;
    inviteCode: string;
    profileID: string;
    created_at: string;
    updated_at: string;
}

async function getServers() {
    try {
        const servers = await axios.get("/servers", {
            headers: {
                Cookie: cookies().toString()
            }
        })

        return servers.data.servers;
    } catch (error) {
        return [];
    }
}

export const NavigationSidebar = async () => {
    const servers = await getServers();

    return (
        <div
            className="space-y-4 flex flex-col items-center h-full text-primary w-full dark:bg-[#1E1F22] bg-[#E3E5E8] py-3"
        >
            <NavigationAction />
            <Separator
                className="h-[2px] bg-zinc-300 dark:bg-zinc-700 rounded-md w-10 mx-auto"
            />
            <ScrollArea className="flex-1 w-full">
                {servers.map((server: Server) => (
                    <div key={server.id} className="mb-4">
                        <NavigationItem
                            id={server.id}
                            name={server.name}
                            imageUrl={server.imageUrl}
                        />
                    </div>
                ))}
            </ScrollArea>
            <div className="pb-3 mt-auto flex items-center flex-col gap-y-4">
                <ModeToggle />
                { /* TODO: ADD User Button Here! */ }
            </div>
        </div>
    );
}
