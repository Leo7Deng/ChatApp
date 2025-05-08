import { useEffect, useRef } from "react";
import { Circle, Message, WebSocketMessage } from "../types";
import { useAuth } from "../context/authContext";



export function useWebSocketDashboard(setCircles: React.Dispatch<React.SetStateAction<Circle[]>>, setAllMessages: React.Dispatch<React.SetStateAction<{ [key: string]: Message[] }>> ) {
    const ws = useRef<WebSocket | null>(null);
    const lastSentMessageTime = useRef<string>("");
    const authContext = useAuth();
    if (!authContext) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    const { getAccessToken } = authContext;

    // Connect to WebSocket
    useEffect(() => {
        let wsCurrent: WebSocket;

        (async () => {
            const token = await getAccessToken();
            wsCurrent = new WebSocket("ws://localhost:8000/ws", [token]);
            wsCurrent.onopen = () => console.log("ws opened");
            wsCurrent.onclose = () => console.log("ws closed");
            ws.current = wsCurrent;
        })();

        return () => {
            if (wsCurrent) {
                wsCurrent.close();
            }
        };
    }, [getAccessToken]);

    // Handle incoming messages
    useEffect(() => {
        if (!ws.current) return;
        ws.current.onmessage = (e) => {
            try {
                var parsedData = JSON.parse(e.data);
                switch (parsedData.origin) {
                    case "server": {
                        switch (parsedData.type) {
                            case "message": {
                                switch (parsedData.action) {
                                    case "create": {
                                        const message: Message = parsedData.message;
                                        const circleID = message.circle_id;

                                        const isDuplicate = lastSentMessageTime.current === message.created_at;
                                        if (isDuplicate) {
                                            const sentTime = new Date(message.created_at).getTime();
                                            const receivedTime = Date.now();
                                            console.log(`Message round-trip time: ${receivedTime - sentTime}ms`);
                                            break;
                                        }

                                        setAllMessages(prev => {
                                            const current = prev[circleID] ?? [];
                                            return {
                                                ...prev,
                                                [circleID]: [...current, message],
                                            };
                                        });
                                        break;
                                    }
                                    case "delete": {
                                        parsedData = parsedData.circle;
                                        const circleIDToDelete = parsedData.circle_id;
                                        console.log("Removing circle with ID:", circleIDToDelete);
                                        setCircles((prevCircles) => prevCircles.filter((circle) => circle.id !== circleIDToDelete));
                                        setAllMessages(prev => ({
                                            ...prev,
                                            [circleIDToDelete]: [],
                                        }));
                                        break;
                                    }
                                    default: {
                                        console.error("Unknown message action:", parsedData.action);
                                    }
                                }
                                break;
                            }
                            case "circle": {
                                switch (parsedData.action) {
                                    case "create": {
                                        parsedData = parsedData.circle;
                                        console.log("Adding circle:", parsedData.id);
                                        const newCircle : Circle = {
                                            id: parsedData.id,
                                            name: parsedData.name,
                                            created_at: parsedData.created_at,
                                        };
                                        setCircles((prevCircles) => [...prevCircles, newCircle]);
                                        break;
                                    }
                                    case "delete": {
                                        parsedData = parsedData.circle;
                                        const circleID = parsedData.id;
                                        console.log("Removing circle with ID:", circleID);
                                        setCircles((prevCircles) => prevCircles.filter((circle) => circle.id !== circleID));
                                        setAllMessages(prev => ({
                                            ...prev,
                                            [circleID]: [],
                                        }));
                                        break;
                                    }
                                    default: {
                                        console.error("Unknown message action:", parsedData.action);
                                    }
                                }
                                break;
                            }
                            default: {
                                console.error("Unknown message type:", parsedData.type);
                                return;
                            }
                        }
                        break;
                    }
                    case "client": {
                        break;
                    }
                    default: {
                        console.error("Unknown origin:", parsedData.origin);
                        return;
                    }
                }
            } catch (error) {
                console.log(error);
                console.error("Error parsing WebSocket message:", error);
            }
            // print a list of all messages
            // console.log("allMessages:", JSON.stringify(allMessages, null, 2));
        };
    }, [ws.current]);

    // Function to send messages
    const sendMessage = (messagePayload: WebSocketMessage) => {
        if (messagePayload.message) {
            lastSentMessageTime.current = messagePayload.message.created_at;
        } else {
            console.error("Message payload is undefined");
        }

        if (ws.current?.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify(messagePayload));
        } else {
            console.error("WebSocket is not open");
        }
    };

    return { sendMessage };
}
