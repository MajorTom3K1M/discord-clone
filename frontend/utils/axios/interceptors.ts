import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { serverRefreshToken } from '@/services/auth/serverRefreshToken';
import { clientRefreshToken } from '@/services/auth/clientRefreshToken';
import { handleRefreshTokenError } from '@/services/auth/errorHandlers';

interface RetryQueueItem {
  resolve: (value: any) => void;
  reject: (error: any) => void;
  config: AxiosRequestConfig & { headers: { Cookie?: string } };
}

export const setupInterceptors = (axiosInstance: AxiosInstance) => {
  let isRefreshing = false;
  const refreshAndRetryQueue: RetryQueueItem[] = [];

  axiosInstance.interceptors.request.use(
    config => config,
    error => Promise.reject(error)
  );

  axiosInstance.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest: AxiosRequestConfig & { _retry: boolean, headers: { Cookie: string } } = error.config;

      if (!axios.isCancel(error) && error.response && error.response.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;

        if (!isRefreshing) {
          isRefreshing = true;
          try {
            let token: { accessToken?: string } | null;
            const isServer = typeof window === 'undefined';

            if (isServer) {
              token = await serverRefreshToken();
              if(token?.accessToken)
                originalRequest.headers['Cookie'] = `access_token=${token?.accessToken};`;
            } else {
              await clientRefreshToken();
            }

            isRefreshing = false;

            refreshAndRetryQueue.forEach(({ config, resolve, reject }) => {
              if (token?.accessToken) {
                const retryRequest: AxiosRequestConfig & { headers: { Cookie?: string } } = config;
                retryRequest.headers['Cookie'] = `access_token=${token?.accessToken};`;
              }
              axiosInstance
                .request(config)
                .then((response) => resolve(response))
                .catch((err) => reject(err));
            });

            refreshAndRetryQueue.length = 0;

            return axiosInstance(originalRequest);
          } catch (refreshError) {
            isRefreshing = false;

            refreshAndRetryQueue.forEach(({ reject }) => reject(refreshError));
            refreshAndRetryQueue.length = 0;

            return Promise.reject(error);
          }
        }
        
        return new Promise<void>((resolve, reject) => {
          refreshAndRetryQueue.push({
            config: originalRequest,
            resolve: (response) => {
              originalRequest._retry = false;
              resolve(response);
            },
            reject: (error) => {
              originalRequest._retry = false;
              reject(error);
            },
          });
        });
      }

      return Promise.reject(error);
    }
  );
};


// axiosInstance.interceptors.response.use(
//     reponse => reponse,
//     async error => {
//         const originalRequest = error.config;


//         if (!axios.isCancel(error) && error.response.status === 401 && !originalRequest._retry) {
//             originalRequest._retry = true;
//             originalRequest.sent = true;

//             try {
//                 const cookies = originalRequest.headers['Cookie'];
//                 await refreshToken(cookies);

//                 // Retry all requests in the queue with the new cookies
//                 refreshAndRetryQueue.forEach(({ config, resolve, reject }) => {
//                     axiosInstance
//                         .request(config)
//                         .then((response) => resolve(response))
//                         .catch((err) => reject(err));
//                 });

//                 // Clear the Queue
//                 refreshAndRetryQueue.length = 0;

//                 // Retry original request
//                 return axiosInstance(originalRequest);
//             } catch (error) {
//                 handleRefreshTokenError();
//                 return Promise.reject(error);
//             }
//         }

//         return Promise.reject(error);
//     }
// )