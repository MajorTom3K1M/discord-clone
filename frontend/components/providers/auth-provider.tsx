import React, { createContext, useState, useEffect } from 'react';
import axios from 'axios';

export const AuthContext = createContext({
    profile: null,
    loading: false,
    signin: (email: string, password: string) => { },
    signup: (name: string, imageUrl: string, email: string, password: string) => { },
    signout: () => { },
    refreshToken: () => { },
    fetchProfile: () => { }
});

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
    const [profile, setProfile] = useState(null);
    const [loading, setLoading] = useState(true);

    const signin = async (email: string, password: string) => {
        try {
            const response = await axios.post('/signin', { email, password });
            setProfile(response.data.profile);
            return response.data;
        } catch (error) {
            throw error;
        }
    };

    const signup = async (email: string, password: string) => {
        try {
            const response = await axios.post('/signup', { email, password });
            setProfile(response.data.profile);
            return response.data;
        } catch (error) {
            throw error;
        }
    };

    const signout = async () => {
        try {
            await axios.post('/signout');
            setProfile(null);
        } catch (error) {
            throw error;
        }
    };

    const refreshToken = async () => {
        try {
            const response = await axios.post("/refresh");
            return response.data.token;
        } catch (error) {
            throw error;
        }
    };

    const fetchProfile = async () => {
        try {
            const response = await axios.get("/me");
            return response.data;
        } catch (error) {
            throw error;
        }
    };

    useEffect(() => {
        const initAuth = async () => {
            const token = "";
            if (token) {
                try {
                    const profile = await fetchProfile();
                    setProfile(profile);
                } catch {
                    // If token validation fails, remove the token
                }
            }
            setLoading(false);
        }
        initAuth();
    }, []);

    return (
        <AuthContext.Provider value={{ profile, loading, signin, signout, signup, refreshToken, fetchProfile }}>
            {children}
        </AuthContext.Provider>
    )
};