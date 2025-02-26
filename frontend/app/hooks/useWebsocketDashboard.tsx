import { useEffect, useRef } from "react";
import { Circle, Message, WebSocketMessage } from "../types";
import { useAuth } from "../context/authContext";



export function useWebSocketDashboard(setCircles: React.Dispatch<React.SetStateAction<Circle[]>>, allMessages: { [key: string]: Message[] }, setAllMessages: React.Dispatch<React.SetStateAction<{ [key: string]: Message[] }>> ) {
    const ws = useRef<WebSocket | null>(null);
    const lastSentMessageTime = useRef<string>("");
    const authContext = useAuth();
    if (!authContext) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    const { getAccessToken } = authContext;

    // Connect to WebSocket
    useEffect(() => {
        const connectWebSocket = async () => {
            const token = await getAccessToken();
            ws.current = new WebSocket("ws://localhost:8000/ws", [token]);
            ws.current.onopen = () => console.log("ws opened");
            ws.current.onclose = () => console.log("ws closed");
            const wsCurrent = ws.current;
            return () => {
                wsCurrent.close();
            };
        };
        connectWebSocket();
    }, []);

    // Handle incoming messages
    useEffect(() => {
        if (!ws.current) return;
        ws.current.onmessage = (e) => {
            try {
                var parsedData = JSON.parse(e.data);
                switch (parsedData.origin) {
                    case "server":
                        switch (parsedData.type) {
                            case "message":
                                switch (parsedData.action) {
                                    case "create":
                                        const circleID = parsedData.circle_id;
                                        parsedData = parsedData.message;
                                        const isDuplicate = lastSentMessageTime.current === parsedData.created_at
                                        console.log("isDuplicate:", isDuplicate);
                                        if (isDuplicate) {
                                            const sentTime = new Date(parsedData.created_at).getTime();
                                            const receivedTime = new Date().getTime();
                                            const timeDifference = receivedTime - sentTime;
                                            console.log(`Message round-trip time: ${timeDifference}ms`);
                                            break;
                                        }
                                        setAllMessages({
                                            ...allMessages,
                                            [parsedData.circle_id]: [
                                                ...allMessages[parsedData.circle_id] || [],
                                                parsedData,
                                            ],
                                        });
                                        break;
                                    case "delete":
                                        parsedData = parsedData.circle;
                                        const circleIDToDelete = parsedData.circle_id;
                                        console.log("Removing circle with ID:", circleID);
                                        setCircles((prevCircles) => prevCircles.filter((circle) => circle.id !== circleIDToDelete));
                                        setAllMessages({
                                            ...allMessages,
                                            [circleIDToDelete]: [],
                                        });
                                        break;
                                    default:
                                        console.error("Unknown message action:", parsedData.action);
                                }
                                break;
                            case "circle":
                                switch (parsedData.action) {
                                    case "create":
                                        parsedData = parsedData.circle;
                                        console.log("Adding circle:", parsedData.id);
                                        const newCircle : Circle = {
                                            id: parsedData.id,
                                            name: parsedData.name,
                                            created_at: parsedData.created_at,
                                        };
                                        setCircles((prevCircles) => [...prevCircles, newCircle]);
                                        break;
                                    case "delete":
                                        parsedData = parsedData.circle;
                                        const circleID = parsedData.id;
                                        console.log("Removing circle with ID:", circleID);
                                        setCircles((prevCircles) => prevCircles.filter((circle) => circle.id !== circleID));
                                        setAllMessages({
                                            ...allMessages,
                                            [circleID]: [],
                                        });
                                        break;
                                    default:
                                        console.error("Unknown message action:", parsedData.action);
                                }
                                break;
                            default:
                                console.error("Unknown message type:", parsedData.type);
                                return;
                        }
                        break;
                    case "client":
                        break;
                    default:
                        console.error("Unknown origin:", parsedData.origin);
                        return;
                }
            } catch (error) {
                console.log(error);
                console.error("Error parsing WebSocket message:", error);
            }
        };
    }, [ws.current]);

    // Function to send messages
    const sendMessage = (messagePayload: WebSocketMessage) => {
        const currentTime = new Date().toISOString();
        lastSentMessageTime.current = currentTime;

        if (ws.current?.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify(messagePayload));
        } else {
            console.error("WebSocket is not open");
        }
    };

    return { sendMessage };
}
