import axios from '@/utils/axios';
import { cookies } from 'next/headers';

async function getServers() {
    try {
        const servers = await axios.get("/servers", {
            headers: {
                Cookie: cookies().toString()
            }
        })
 
        return servers.data;
    } catch (error) {
        return { servers: [] };
    }
}

export const NavigationSidebar = async () => {
    const servers = await getServers();

    return (
        <div>
            Navigation Sidebar
        </div>
    );
}
