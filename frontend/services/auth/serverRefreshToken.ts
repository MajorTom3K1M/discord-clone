"use server"
import axios from "@/utils/axios";
import { cookies } from 'next/headers';

export const serverRefreshToken = async () => {
    try {
        const response = await axios.get("/server/refresh", {
            withCredentials: true,
            headers: {
                Cookie: cookies().toString()
            }
        });
    
        const accessToken = response.data.access_token
        return accessToken;
    } catch (err) {
        throw err;
    }
};

const parseCookies = (setCookieArray: string[]) => {
    const cookies: any = {};
    
    setCookieArray.forEach(cookieString => {
      const [firstPart] = cookieString.split('; ');
      const [key, value] = firstPart.split('=');
      cookies[key] = value;
    });
    
    return cookies;
  }