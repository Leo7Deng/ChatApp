"use client";

import React, { useState, useContext, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";

interface AuthContextType {
    getAccessToken: () => Promise<string>;
}

const AuthContext = React.createContext<AuthContextType | undefined>(undefined);

export function useAuth() {
    return useContext(AuthContext);
}

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [accessToken, setAccessToken] = useState("");
    const [expiryTime, setExpiryTime] = useState(0);
    const [shouldRedirect, setShouldRedirect] = useState(false);

    async function getAccessToken(): Promise<string> {
        console.log(Date.now(), expiryTime);
        if (!accessToken || Date.now() >= expiryTime) {
            const token = await refreshAccessToken();
            return token;
        }
        return accessToken;
    }

    const router = useRouter();
    let pendingPromise: Promise<string> | null = null;

    async function refreshAccessToken(): Promise<string> {
        if (pendingPromise) {
            return pendingPromise;
        }

        pendingPromise = new Promise(async (resolve, reject) => {
            try {
                const response = await fetch("http://localhost:8000/refresh", {
                    method: "POST",
                    credentials: "include",
                });
                const data = await response.json();
                if (!response.ok) {
                    setShouldRedirect(true);
                    return "";
                }
                setAccessToken(data.access_token);
                setExpiryTime(Date.now() + 15 * 60 * 1000); // 15 minutes
                resolve(data.access_token);
            }
            catch (error) {
                setShouldRedirect(true);
                console.log("Error refreshing access token:", error);
                reject(error);
            }
            finally {
                pendingPromise = null;
            }
        }
        );
        return pendingPromise;
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