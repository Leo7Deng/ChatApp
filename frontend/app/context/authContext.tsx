"use client";

import React, { useState, useContext, useEffect } from "react";
import { useRouter } from "next/navigation";

interface AuthContextType {
    getAccessToken: () => Promise<string>;
}

const AuthContext = React.createContext<AuthContextType | undefined>(undefined);

// Hook to use the Auth Context
export function useAuth() {
    return useContext(AuthContext);
}

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [accessToken, setAccessToken] = useState("");
    const [expiryTime, setExpiryTime] = useState(0);
    const [shouldRedirect, setShouldRedirect] = useState(false);

    async function getAccessToken(): Promise<string> {
        if (Date.now() >= expiryTime) {
            const token = await refreshAccessToken();
            return token;
        }
        return accessToken;
    }

    useEffect(() => {
        if (typeof window === "undefined") return;
        getAccessToken();
    }, []);

    const router = useRouter();
    async function refreshAccessToken(): Promise<string> {
        try {
            const response = await fetch("http://localhost:8000/refresh", {
                method: "POST",
                credentials: "include",
            });
            const data = await response.json();
            if (!response.ok) {
                setShouldRedirect(true);
                console.log("Failed to refresh access token");
            }
            setAccessToken(data.access_token);
            setExpiryTime(Date.now() + 15 * 60 * 1000); // 15 minutes
            return data.access_token;
        } catch (error) {
            setShouldRedirect(true);
            console.log("Error refreshing access token:", error);
            return "";
        }
    }

    useEffect(() => {
        if (shouldRedirect) {
            router.push("/login");
        }
    }, [shouldRedirect, router]);

    return (
        <AuthContext.Provider value={{ getAccessToken }}>
            {children}
        </AuthContext.Provider>
    );
};