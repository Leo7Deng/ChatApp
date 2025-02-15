"use client"

import "./Dashboard.css"
import React, { useEffect, useRef, useState } from "react";
import CreateCircleModal from "./CreateCircleModal";
import InviteModal from "./InviteModal";

interface Circle {
    id: string;
    name: string;
    created_at: string;
}

interface WebSocketMessage {
    origin: "server" | "client";
    type: "message" | "circle";
    action: "create" | "delete";
    message?: Message;
    circle?: Circle;
}

interface Message {
    circle_id: string;
    content: string;
    created_at: string;
    author_id: string;
}

interface Circle {
    id: string;
    name: string;
    createdAt?: string;
}

export default function Dashboard() {
    // Loading spinner
    const [isFetching, setIsFetching] = useState(true);

    // Store messages
    const [allMessages, setAllMessages] = useState<{ [key: string]: Message[] }>({});
    const lastSentMessageTime = useRef("");

    // Store access token
    const [accessToken, setAccessToken] = useState("");
    const [expiryTime, setExpiryTime] = useState(0);

    // Refresh access token
    async function refreshAccessToken() {
        if (new Date().getTime() < expiryTime || accessToken !== "") {
            return accessToken;
        }
        const response = await fetch('https://127.0.0.1:8000/refresh', {
            method: 'POST',
            credentials: 'include',
        });
        const data = await response.json();
        if (!response.ok) {
            console.log("Error in refreshing access token");
        }
        else {
            console.log("Data:", data);
            setAccessToken(data.access_token);
            setExpiryTime(new Date().getTime());
            return data.access_token
        }
    }

    // Get circles from the server
    const [circles, setCircles] = useState<Circle[]>([]);
    useEffect(() => {
        async function fetchUserData() {
            const token = await refreshAccessToken();
            const headers = {
                "Authorization": `Bearer ${token}`,
            };
            fetch('https://127.0.0.1:8000/api/circles', {
                method: 'GET',
                credentials: 'include',
                headers: headers,
            })
                .then(async (response) => {
                    const data = await response.json();
                    if (!response.ok) {
                        if (data == "refresh token not found") {
                            window.location.href = "./login";
                        }
                        console.log("Error in getting circles");
                    }
                    else {
                        console.log("Data:", data);
                        const mappedCircles = data.map((circle: any) => ({
                            id: circle.id,
                            name: circle.name,
                            created_at: circle.created_at,
                        }));
                        setCircles(mappedCircles);
                        setIsFetching(false);
                    }
                })
                .catch(error => {
                    console.log(error);
                });
        }
        fetchUserData();
    }, []);

    // Create a new circle
    const [openModal, setOpenModal] = useState(false);
    function createCircle() {
        setOpenModal(true);
    }
    interface HandleCloseEvent extends React.MouseEvent<HTMLDivElement> {
        target: EventTarget & HTMLDivElement;
    }
    const handleClose = (event: HandleCloseEvent) => {
        if (event.target.classList.contains('modal-container')) {
            setOpenModal(false);
        }
    };

    // Invite users to circle
    const [openInviteModal, setOpenInviteModal] = useState(false);
    function handleOpenInviteModal() {
        setOpenInviteModal(true);
    }
    const handleInviteClose = (event: HandleCloseEvent) => {
        if (event.target.classList.contains('modal-container')) {
            setOpenInviteModal(false);
        }
    };

    // Delete circle
    const handleDelete = async () => {
        const token = await refreshAccessToken();
        const headers = {
            "Authorization": `Bearer ${token}`,
        };
        fetch(`https://127.0.0.1:8000/api/circles/delete/${selectedCircleID}`, {
            method: 'DELETE',
            headers: headers,
            credentials: 'include',
        })
            .then(async (response) => {
                const data = await response.json();
                if (!response.ok) {
                    if (data == "user is not admin of circle") {
                        alert("You are not an admin of the circle");
                    }
                    console.log("data: ", data);
                }
                else {
                    console.log("Data:", data);
                    setCircles((prevCircles) => prevCircles.filter((circle) => circle.id !== selectedCircleID));
                }
            })
            .catch(error => {
                console.log(error);
            });
    }

    // Connect to WebSocket
    const ws = useRef<WebSocket>();
    useEffect(() => {
        const connectWebSocket = async () => {
            const token = await refreshAccessToken();
            ws.current = new WebSocket("wss://127.0.0.1:8000/ws", [token]);
            ws.current.onopen = () => console.log("ws opened");
            ws.current.onclose = () => console.log("ws closed");
            const wsCurrent = ws.current;
            return () => {
                wsCurrent.close();
            };
        };
        connectWebSocket();
    }, []);

    // Request for userID and username
    const [userID, setUserID] = useState("");
    const [username, setUsername] = useState("");
    useEffect(() => {
        async function fetchUserData() {
            const token = await refreshAccessToken();
            const headers = {
                "Authorization": `Bearer ${token}`,
            };
            fetch('https://127.0.0.1:8000/api/user', {
                method: 'GET',
                headers: headers,
                credentials: 'include'
            })
                .then(async (response) => {
                    const data = await response.json();
                    if (!response.ok) {
                        console.log("Error in getting user data");
                    }
                    else {
                        console.log("Data:", data);
                        setUserID(data.user_id);
                        setUsername(data.username);
                    }
                })
                .catch(error => {
                    console.log(error);
                });
        }
        fetchUserData();
    }, []);


    // Send message to WebSocket
    interface HandleEnterEvent extends React.KeyboardEvent<HTMLInputElement> { }
    const handleEnter = (event: HandleEnterEvent) => {
        if (!userID) {
            console.log("User ID not found");
            window.location.href = "./login";
            return;
        }
        if (event.currentTarget.value === '') {
            return;
        }
        if (event.key === 'Enter') {
            console.log("Enter key pressed");
            var currentTime = new Date().toISOString();
            lastSentMessageTime.current = currentTime;
            const message: Message = {
                circle_id: selectedCircleID,
                content: event.currentTarget.value,
                created_at: currentTime,
                author_id: String(userID),
            }
            const messagePayload: WebSocketMessage = {
                origin: "client",
                type: "message",
                action: "create",
                message: message,
            };
            if (!messagePayload.message) {
                console.error("Message payload is empty");
                return;
            }
            console.log("Sent message:", messagePayload);
            event.currentTarget.value = '';
            if (!ws.current) {
                console.error("WebSocket not connected");
                return;
            }
            ws.current.send(JSON.stringify(messagePayload));
            setAllMessages({
                ...allMessages,
                [selectedCircleID]: [
                    ...allMessages[selectedCircleID] || [],
                    message,
                ],
            });
        }
    };

    // Receive message from WebSocket
    useEffect(() => {
        if (!ws.current) return;
        ws.current.onmessage = (e) => {
            try {
                var parsedData = JSON.parse(e.data);
                console.log("Parsed data:", parsedData);
                if (parsedData.origin == "client") {
                    return;
                }
                if (parsedData.type == "message") {
                    console.log("Received message:", parsedData.message);
                    if (parsedData.action === "create") {
                        parsedData = parsedData.message;
                        console.log("Adding message:", parsedData.content);
                        const isDuplicate = lastSentMessageTime.current === parsedData.created_at
                        if (isDuplicate) {
                            const sentTime = new Date(parsedData.created_at).getTime();
                            const receivedTime = new Date().getTime();
                            const timeDifference = receivedTime - sentTime;
                            console.log(`Message round-trip time: ${timeDifference}ms`);
                            return;
                        }
                        setAllMessages({
                            ...allMessages,
                            [parsedData.circle_id]: [
                                ...allMessages[parsedData.circle_id] || [],
                                parsedData,
                            ],
                        });
                    }
                    else if (parsedData.action == "delete") {
                        parsedData = parsedData.message;
                        const messageCreatedAt = parsedData.created_at;
                        console.log("Removing message with created_at:", messageCreatedAt);
                        setAllMessages({
                            ...allMessages,
                            [parsedData.circle_id]: allMessages[parsedData.circle_id].filter((message) => message.created_at !== messageCreatedAt),
                        });
                    }
                }
                else if (parsedData.type == "circle") {
                    if (parsedData.action == "create") {
                        parsedData = parsedData.circle;
                        console.log("Adding circle:", parsedData.id);
                        const newCircle = {
                            id: parsedData.id,
                            name: parsedData.name,
                            created_at: parsedData.created_at,
                        };
                        setCircles((prevCircles) => [...prevCircles, newCircle]);
                    }
                    else if (parsedData.action == "delete") {
                        parsedData = parsedData.circle;
                        const circleID = parsedData.id;
                        console.log("Removing circle with ID:", circleID);
                        setCircles((prevCircles) => prevCircles.filter((circle) => circle.id !== circleID));
                        setAllMessages({
                            ...allMessages,
                            [circleID]: [],
                        });
                    }
                }
                else {
                    console.log("Unknown message type:", parsedData.type);
                }
            } catch (error) {
                console.error("Error parsing WebSocket message:", error);
            }
        };
    }, [ws.current]);

    const [selectedCircleID, setSelectedCircleID] = useState("");

    // Scroll when messages are updated
    useEffect(() => {
        const chatContainer = document.querySelector(".chat-placeholder") as HTMLElement;
        if (chatContainer) {
            const textBoxHeight = window.innerHeight * 0.94;
            chatContainer.scrollTop = chatContainer.scrollHeight - textBoxHeight;
        }
    }, [allMessages[selectedCircleID]]);

    return (
        <>
            {isFetching ?
                <div role="status" className="spinner-container">
                    <svg aria-hidden="true" className="w-8 h-8 text-gray-200 animate-spin spinner" viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" fill="currentColor" />
                        <path d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" fill="currentFill" />
                    </svg>
                </div>
                :
                <div className="dashboard">
                    <div className="flex h-screen w-16 flex-col justify-between border-e bg-white sidebar">
                        <div>
                            <div className="inline-flex size-16 items-center justify-center">
                                <span className="grid size-10 place-content-center rounded-lg bg-gray-100 text-xs text-gray-600">
                                    L
                                </span>
                            </div>

                            <div className="border-t" style={{ borderColor: 'hsl(0, 0%, 87%)' }}>
                                <div className="px-2">
                                    <div className="py-2">
                                        <a
                                            className="t group relative flex justify-center rounded px-2 py-1.5 text-gray-500 hover:bg-gray-50 hover:text-gray-700"
                                        >
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                className="size-5 opacity-75"
                                                fill="none"
                                                viewBox="0 0 24 24"
                                                stroke="currentColor"
                                                strokeWidth="2"
                                            >
                                                <path
                                                    strokeLinecap="round"
                                                    strokeLinejoin="round"
                                                    d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                                                />
                                                <path
                                                    strokeLinecap="round"
                                                    strokeLinejoin="round"
                                                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                                                />
                                            </svg>

                                            <span
                                                className="invisible absolute start-full top-1/2 ms-4 -translate-y-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible whitespace-nowrap"
                                            >
                                                General
                                            </span>
                                        </a>
                                    </div>

                                    <ul className="space-y-1 border-t pt-4" style={{ borderColor: 'hsl(0, 0%, 87%)' }}>
                                        {circles.map((circle) => (
                                            <li key={circle.id} className="flex justify-center">
                                                <button onClick={() => setSelectedCircleID(circle.id)} className="w-full">
                                                    <a
                                                        className={`group relative flex justify-center rounded px-2 py-1.5 ${selectedCircleID === circle.id
                                                            ? "bg-blue-100 text-gray-700"
                                                            : "text-gray-500 hover:bg-gray-50 hover:text-gray-700"
                                                            }`}
                                                    >
                                                        <svg
                                                            xmlns="http://www.w3.org/2000/svg"
                                                            className="size-5 opacity-75"
                                                            fill="none"
                                                            viewBox="0 0 24 24"
                                                            stroke="currentColor"
                                                            strokeWidth="2"
                                                        >
                                                            <path
                                                                strokeLinecap="round"
                                                                strokeLinejoin="round"
                                                                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                                                            />
                                                        </svg>

                                                        <span
                                                            className="invisible absolute start-full top-1/2 ms-4 -translate-y-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible whitespace-nowrap"
                                                        >
                                                            {circle.name}
                                                        </span>
                                                    </a>
                                                </button>
                                            </li>
                                        ))}
                                        <li>
                                            <a
                                                className="group relative flex justify-center rounded px-2 py-1.5 text-gray-500 hover:bg-gray-50 hover:text-gray-700"
                                            >
                                                <svg
                                                    xmlns="http://www.w3.org/2000/svg"
                                                    className="size-5 opacity-75"
                                                    fill="none"
                                                    viewBox="0 0 24 24"
                                                    stroke="currentColor"
                                                    strokeWidth="2"
                                                >
                                                    <path
                                                        strokeLinecap="round"
                                                        strokeLinejoin="round"
                                                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                                                    />
                                                </svg>

                                                <span
                                                    className="invisible absolute start-full top-1/2 ms-4 -translate-y-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible"
                                                >
                                                    Account
                                                </span>
                                            </a>
                                        </li>


                                        <li className="flex justify-center">
                                            <button
                                                className="group relative flex justify-center rounded px-2 py-1.5 text-gray-500 hover:bg-gray-50 hover:text-gray-700 w-full"
                                                onClick={createCircle}
                                            >
                                                <svg
                                                    id="Layer_1"
                                                    data-name="Layer 1"
                                                    xmlns="http://www.w3.org/2000/svg"
                                                    viewBox="0 0 24 24"
                                                    strokeWidth="2"
                                                    opacity={0.75}
                                                    width="24"
                                                    height="24"
                                                    stroke="currentColor"
                                                    className="size-5 opacity-75"
                                                >
                                                    <defs>
                                                        <style>
                                                            {`.cls-637642e7c3a86d32eae6f177-1{fill:none;stroke:currentColor;stroke-miterlimit:10;}`}
                                                        </style>
                                                    </defs>
                                                    <path
                                                        className="cls-637642e7c3a86d32eae6f177-1"
                                                        d="M18.68,8.16V15.8a2.86,2.86,0,0,1-2.86,2.86H13.91v2.86L8.18,18.66H4.36A2.86,2.86,0,0,1,1.5,15.8V8.16A2.86,2.86,0,0,1,4.36,5.3H15.82A2.86,2.86,0,0,1,18.68,8.16Z"
                                                    ></path>
                                                    <path
                                                        className="cls-637642e7c3a86d32eae6f177-1"
                                                        d="M18.68,14.84h1A2.86,2.86,0,0,0,22.5,12V4.34a2.86,2.86,0,0,0-2.86-2.86H8.18A2.86,2.86,0,0,0,5.32,4.34v1"
                                                    ></path>
                                                    <line
                                                        className="cls-637642e7c3a86d32eae6f177-1"
                                                        x1="6.27"
                                                        y1="11.98"
                                                        x2="13.91"
                                                        y2="11.98"
                                                    ></line>
                                                    <line
                                                        className="cls-637642e7c3a86d32eae6f177-1"
                                                        x1="10.09"
                                                        y1="8.16"
                                                        x2="10.09"
                                                        y2="15.8"
                                                    ></line>
                                                </svg>
                                                <span
                                                    className="invisible absolute start-full top-1/2 ms-4 -translate-y-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible whitespace-nowrap"
                                                >
                                                    Add Circle
                                                </span>
                                            </button>
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        </div>

                        <div className="sticky inset-x-0 bottom-0 border-t bg-white p-2" style={{ borderColor: 'hsl(0, 0%, 87%)' }}>
                            <button
                                className="group relative flex w-full justify-center rounded-lg px-2 py-1.5 text-sm text-gray-500 hover:bg-gray-50 hover:text-gray-700"
                                onClick={() => {
                                    window.location.href = "./login";
                                }}
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    className="size-5 opacity-75"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                    strokeWidth="2"
                                >
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                                    />
                                </svg>
                                <span
                                    className="invisible absolute start-full top-1/2 ms-4 -translate-y-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible"
                                >
                                    Logout
                                </span>
                            </button>
                        </div>
                        {openModal && <div className="modal-container" onClick={handleClose}>
                            <div className="modal inset-0 bg-black bg-opacity-50 z-50">
                                <CreateCircleModal isOpen={openModal} setOpen={() => setOpenModal(false)} />
                            </div>
                        </div>}
                        {openInviteModal && <div className="modal-container" onClick={handleInviteClose}>
                            <div className="modal inset-0 bg-black bg-opacity-50 z-50">
                                <InviteModal isOpen={openInviteModal} setOpen={() => setOpenInviteModal(false)} circleId={selectedCircleID} />
                            </div>
                        </div>}
                    </div>
                    <div className="chat-container">
                        <div className="chat-title">
                            <h1>{circles.find((circle) => circle.id === selectedCircleID)?.name}&nbsp;</h1>
                            <div className="chat-menu">
                                <button
                                    className="group relative rounded px-2 py-1.5 text-gray-500 hover:bg-gray-50 hover:text-gray-700 h-full disabled:opacity-50 disabled:cursor-not-allowed"
                                    onClick={selectedCircleID !== "" ? handleOpenInviteModal : undefined}
                                    disabled={selectedCircleID === ""}
                                >
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        className="size-5 opacity-75"
                                        fill="currentColor"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                        strokeWidth="1"
                                    >
                                        <g transform="matrix(0.71 0 0 0.71 12 12)" >
                                            <path
                                                transform=" translate(-16, -16)"
                                                d="M 12 2 C 8.144531 2 5 5.144531 5 9 C 5 11.410156 6.230469 13.550781 8.09375 14.8125 C 4.527344 16.34375 2 19.882813 2 24 L 4 24 C 4 19.570313 7.570313 16 12 16 C 13.375 16 14.65625 16.359375 15.78125 16.96875 C 14.671875 18.34375 14 20.101563 14 22 C 14 26.40625 17.59375 30 22 30 C 26.40625 30 30 26.40625 30 22 C 30 17.59375 26.40625 14 22 14 C 20.253906 14 18.628906 14.574219 17.3125 15.53125 C 16.871094 15.253906 16.390625 15.019531 15.90625 14.8125 C 17.769531 13.550781 19 11.410156 19 9 C 19 5.144531 15.855469 2 12 2 Z M 12 4 C 14.773438 4 17 6.226563 17 9 C 17 11.773438 14.773438 14 12 14 C 9.226563 14 7 11.773438 7 9 C 7 6.226563 9.226563 4 12 4 Z M 22 16 C 25.324219 16 28 18.675781 28 22 C 28 25.324219 25.324219 28 22 28 C 18.675781 28 16 25.324219 16 22 C 16 18.675781 18.675781 16 22 16 Z M 21 18 L 21 21 L 18 21 L 18 23 L 21 23 L 21 26 L 23 26 L 23 23 L 26 23 L 26 21 L 23 21 L 23 18 Z"
                                                strokeLinecap="round" />
                                        </g>
                                    </svg>
                                    <span
                                        className="invisible absolute start-1/2 top-full mt-4 -translate-x-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible"
                                    >
                                        Invite Users
                                    </span>
                                </button>
                                <button
                                    className="group relative rounded px-2 py-1.5 text-gray-500 hover:bg-gray-50 hover:text-gray-700 h-full disabled:opacity-50 disabled:cursor-not-allowed"
                                    onClick={selectedCircleID !== "" ? handleDelete : undefined}
                                    disabled={selectedCircleID === ""}
                                >

                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        className="size-5 opacity-75"
                                        fill="currentColor"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                        strokeWidth="0.1"
                                    >
                                        <path d="M 10 2 L 9 3 L 4 3 L 4 5 L 5 5 L 5 20 C 5 20.522222 5.1913289 21.05461 5.5683594 21.431641 C 5.9453899 21.808671 6.4777778 22 7 22 L 17 22 C 17.522222 22 18.05461 21.808671 18.431641 21.431641 C 18.808671 21.05461 19 20.522222 19 20 L 19 5 L 20 5 L 20 3 L 15 3 L 14 2 L 10 2 z M 7 5 L 17 5 L 17 20 L 7 20 L 7 5 z M 9 7 L 9 18 L 11 18 L 11 7 L 9 7 z M 13 7 L 13 18 L 15 18 L 15 7 L 13 7 z"></path>
                                    </svg>
                                    <span
                                        className="invisible absolute start-1/2 top-full mt-4 -translate-x-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible"
                                    >
                                        Delete Circle
                                    </span>
                                </button>
                            </div>
                        </div>
                        <div className="chat-placeholder">
                            {selectedCircleID !== "" && (
                                <div className="chat">
                                    {allMessages[selectedCircleID]?.map((message) => (
                                        <div key={message.created_at} className={`message ${message.author_id === String(userID) ? "owner" : "other"}`}>
                                            <div className="message-content">
                                                <p>{message.content}</p>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                        <div className="text-box">
                            <input type="text" className={`text-input disabled:opacity-50 ${selectedCircleID === "" ? "hover:cursor-not-allowed" : ""}`} placeholder="Type a message" onKeyDown={handleEnter} disabled={selectedCircleID === ""} />
                        </div>
                    </div>
                    <div className="analytics-container">
                        <div className="analytics-placeholder">
                        </div>
                    </div>
                </div>
            }
        </>

    )
};