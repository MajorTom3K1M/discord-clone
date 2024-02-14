import axios, { AxiosInstance } from 'axios';
import { refreshToken } from '@/services/auth/refreshToken';
import { handleRefreshTokenError } from '@/services/auth/errorHandlers';
import { useRouter } from 'next/navigation';

export const setupInterceptors = (axiosInstance: AxiosInstance) => {
    axiosInstance.interceptors.request.use(
        config => config,
        error => Promise.reject(error)
    );

    axiosInstance.interceptors.response.use(
        reponse => reponse,
        async error => {
            const originalRequest = error.config;

      
            if (!axios.isCancel(error) && error.response.status === 401 && !originalRequest._retry) {
                originalRequest._retry = true;

                try {
                    await refreshToken();
                    return axiosInstance(originalRequest);
                } catch (error) {
                    handleRefreshTokenError();
                    return Promise.reject(error);
                }
            }

            return Promise.reject(error);
        }
    )
};