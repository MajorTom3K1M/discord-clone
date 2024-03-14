"use client"
import React, { createContext, useState, useContext, useEffect } from 'react';
import { usePathname, useSearchParams, useRouter, redirect } from 'next/navigation'; // Corrected import
import axios, { fetchDataWithCancellation, isCancel } from '@/utils/axios';

const publicPages = ['/sign-in', '/sign-up'];

interface Profile {
    id: string;
    name: string;
    imageUrl: string;
    email: string;
    servers: any[] | null;
    members: any[] | null;
    channels: any[] | null;
    created_at: string;
    updated_at: string;
}

interface AuthState {
    profile: Profile | null;
}

interface AuthContextType {
    authState: AuthState;
    isPublicPage: boolean;
    signin: (email: string, password: string) => Promise<void>;
    signup: (name: string, imageUrl: string, email: string, password: string) => Promise<void>;
    signout: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType>({
    authState: { profile: null },
    isPublicPage: false,
    signin: async (email: string, password: string) => { },
    signup: async (name: string, imageUrl: string, email: string, password: string) => { },
    signout: async () => { },
});

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
    const pathname = usePathname();
    const searchParams = useSearchParams();
    const router = useRouter();
    const [authState, setAuthState] = useState({
        profile: null
    });
    const isPublicPage = publicPages.includes(pathname);

    const signin = async (email: string, password: string) => {
        const response = await axios.post('/signin', { email, password });
        if (response.status === 200) {
            setAuthState({ profile: response.data.profile });
            const redirectTo = searchParams.get("redirect");
            router.push(redirectTo ?? "/");
        }
        return response.data;
    };

    const signup = async (email: string, password: string) => {
        const response = await axios.post('/signup', { email, password });
        setAuthState(response.data.profile);
        return response.data;
    };

    const signout = async () => {
        await axios.post('/signout');
        setAuthState({ profile: null });
    };

    useEffect(() => {
        const { request, cancel } = fetchDataWithCancellation("/profile/auth/me");

        async function checkAuth() {
            request.then(response => {
                setAuthState({ profile: response.data.profile });
            }).catch(error => {
                if(isCancel(error)) {
                    console.log('Request canceled:', error.message);
                } else if (!isPublicPage) {
                    if (pathname !== "/") {
                        router.push(`/sign-in?redirect=${encodeURIComponent(pathname)}`);
                    } else {
                        router.push("/sign-in");
                    }
                    router.refresh();
                } else {
                    console.error('An error occurred:', error.message);
                }
            })
        }

        checkAuth()
        return () => cancel('Component unmounted: Operation canceled by the user.');
    }, []);

    return (
        <AuthContext.Provider value={{ authState, signin, signout, signup, isPublicPage }}>
            {authState.profile === null && !isPublicPage ? <>Loading...</> : children}
        </AuthContext.Provider>
    )
};

export const useAuth = () => useContext(AuthContext);