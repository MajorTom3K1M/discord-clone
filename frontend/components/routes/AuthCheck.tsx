"use client"

import { useAuth } from "@/components/providers/AuthProvider";
import React, { useState, useLayoutEffect } from "react";
import { useRouter, usePathname } from "next/navigation";

export const AuthCheck = ({ children }: { children: React.ReactNode }) => {
    const { authState, isPublicPage } = useAuth();
    const [isChecking, setIsChecking] = useState(true);

    const router = useRouter();
    const pathname = usePathname();

    useLayoutEffect(() => {
        console.log({ profile: authState.profile });
        if (!authState.profile && !isPublicPage) {
            // router.push(`/sign-in?redirect=${encodeURIComponent(pathname)}`);
        } else {
            setIsChecking(false);
        }
    }, [router, authState.profile, isPublicPage, pathname]);

    return isChecking ? <>Loading...</> : children;
};
