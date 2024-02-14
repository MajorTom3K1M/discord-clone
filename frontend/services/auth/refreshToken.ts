import axios from '@/utils/axios';

export const refreshToken = async () => {
    try {
        const response = await axios.get('/refresh', { withCredentials: true });
        return response.data.accessToken;
    } catch (error) {
        throw new Error('Failed to refresh token');
    }
};