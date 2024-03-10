import axios from '@/utils/axios';

import { cookies } from 'next/headers';

import { cache } from 'react';

import { Profile } from '@/types/models';

export const currentProfile = cache(async (): Promise<Profile | null> => {
    try {
        const profile = await axios.get<{
            profile: Profile
        }>("/profile/auth/me", {
            headers: {
                Cookie: cookies().toString()
            }
        });
    
        return profile.data.profile;
    } catch (err) {
        return null;
    }   
})