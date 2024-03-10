import { ServerSidebar } from '@/components/server/ServerSidebar';
import axios from '@/utils/axios';

import { cookies } from 'next/headers';
import { redirect } from "next/navigation";

async function getServer(serverId: string) {
    try {
        const server = await axios.get(`/servers/${serverId}`, {
            headers: {
                Cookie: cookies().toString()
            }
        })

        return server.data.server;
    } catch (error) {
        return null;
    }
}

const ServerIdLayout = async ({
    children,
    params,
}: {
    children: React.ReactNode,
    params: { serverId: string }
}) => {
    const server = await getServer(params.serverId);

    if (!server) {
        return redirect("/");
    }

    return (
        <div className="h-full">
            <div
                className="hidden md:flex h-full w-60 z-20 flex-col fixed inset-y-0">
                <ServerSidebar serverId={params.serverId} />
            </div>
            <main className="h-full md:pl-60">
                {children}
            </main>
        </div>
    );
}

export default ServerIdLayout;