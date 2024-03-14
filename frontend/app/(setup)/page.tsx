import { InitialModal } from '@/components/modals/InitialModal';
import { Server } from '@/types/models';

import { cookies } from 'next/headers';

import axios from '@/utils/axios';
import { redirect } from 'next/navigation';

const GetServer = async (): Promise<Server | null> => {
    try {
        const server = await axios.get<{
            message: string,
            server: Server
        }>("/servers/by-profile", {
            headers: {
                Cookie: cookies().toString()
            }
        });
        return server.data.server
    } catch (err) {
        return null
    }
}

const SetupPage = async () => {
    const server = await GetServer();

    if (server) {
        return redirect(`/servers/${server.id}`);
    }

    return <InitialModal />
}

export default SetupPage;