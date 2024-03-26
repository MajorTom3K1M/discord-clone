import axios from "@/utils/axios";

import { currentProfile } from "@/lib/currentProfile";
import { Server } from "@/types/models";
import { cookies } from "next/headers";
import { redirect, useSearchParams } from 'next/navigation';

interface ServerIdPageProps {
    params: {
        serverId: string;
    }
}

const getServerGeneralChannel = async (serverId: string) => {
    try {
        const server = await axios.get<{
            message: string,
            server: Server
        }>(`/servers/${serverId}/channels/default`, {
            headers: {
                Cookie: cookies().toString()
            }
        });
        return server.data.server
    } catch (err) {
        return null
    }
};

const ServerIdPage = async ({
    params
}: ServerIdPageProps) => {
    const profile = await currentProfile();

    if (!profile) {
        return redirect("/sign-in");
    }

    const server = await getServerGeneralChannel(params.serverId);

    const initialChannel = server?.channels?.[0];

    if (!initialChannel || initialChannel?.name !== "general") {
        return null;
    }

    return redirect(`/servers/${params.serverId}/channels/${initialChannel?.id}`)
}

export default ServerIdPage;