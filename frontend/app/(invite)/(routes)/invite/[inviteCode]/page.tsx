import axios from "@/utils/axios";

import { currentProfile } from "@/lib/currentProfile";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { Server } from "@/types/models";

const GetServerByInviteCode = async (inviteCode: string): Promise<Server | null> => {
    try {
        const server = await axios.get<{
            message: string,
            server: Server
        }>(`/servers/invite-code/${inviteCode}`, {
            headers: {
                Cookie: cookies().toString()
            }
        })

        return server.data.server;
    } catch (error) {
        console.error(error);
        return null;
    }
};

const CreateServerMember = async (inviteCode: string): Promise<Server | null> => {
    try {
        const server = await axios.post<{
            message: string,
            server: Server
        }>(`/servers/invite-code/${inviteCode}/members`, {
            headers: {
                Cookie: cookies().toString()
            }
        });

        return server.data.server;
    } catch {
        return null;
    }
};

interface InviteCodePageProps {
    params: {
        inviteCode: string;
    };
}

const InviteCodePage = async ({
    params
}: InviteCodePageProps) => {
    const profile = await currentProfile();

    if (!profile) {
        return redirect("/sign-in");
    }

    if (!params.inviteCode) {
        return redirect("/");
    }

    const existingServer = await GetServerByInviteCode(params.inviteCode);

    if (existingServer) {
        return redirect(`/servers/${existingServer.id}`);
    }

    const server = await CreateServerMember(params.inviteCode);

    if (server) {
        return redirect(`/servers/${server.id}`);
    }

    return null;
}

export default InviteCodePage;