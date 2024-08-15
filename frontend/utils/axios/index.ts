import axios from 'axios';
import { setupInterceptors } from './interceptors';

const axiosInstance = axios.create({
    baseURL: process.env.NEXT_PUBLIC_BACKEND_URL,
    withCredentials: true
});

setupInterceptors(axiosInstance);

export const fetchDataWithCancellation = (url: string) => {
    const source = axios.CancelToken.source();
    
    const request = axiosInstance.get(url, {
        cancelToken: source.token,
        withCredentials: true
    });

    return {
        request,
        cancel: source.cancel
    };
};

export const isCancel = (error: string) => {
    return axios.isCancel(error);
};

export default axiosInstance;