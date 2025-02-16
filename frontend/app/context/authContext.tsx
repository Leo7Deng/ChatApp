import React, { useState, useContext, useEffect } from "react";

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

    async function getAccessToken(): Promise<string> {
        if (Date.now() >= expiryTime) {
            const token = await refreshAccessToken();
            return token;
        }
        return accessToken;
    }

    useEffect(() => {
        getAccessToken();
    }, []);

    async function refreshAccessToken(): Promise<string> {
        try {
            const response = await fetch("https://127.0.0.1:8000/refresh", {
                method: "POST",
                credentials: "include",
            });

            if (!response.ok) throw new Error("Failed to refresh token");
            const data = await response.json();

            setAccessToken(data.access_token);
            setExpiryTime(Date.now() + 15 * 60 * 1000); // 15 minutes
            return data.access_token;
        } catch (error) {
            console.error("Error refreshing access token:", error);
            return "";
        }
    }

    return (
        <AuthContext.Provider value={{ getAccessToken }}>
            {children}
        </AuthContext.Provider>
    );
};