"use client"

import "./Dashboard.css"
import React, { useEffect, useRef, useState } from "react";
import CreateCircleModal from "./CreateCircleModal";
import EditModal from "./EditModal";
import InviteModal from "./InviteModal";
import { useFetchDashboard } from "../hooks/useFetchDashboard";
import { useWebSocketDashboard } from "../hooks/useWebsocketDashboard";
import { Circle, Message, WebSocketMessage } from "../types";

export default function Dashboard2() {
    const [isUserDataFetched, setIsUserDataFetched] = useState(false);
    const [selectedCircleID, setSelectedCircleID] = useState("");
    const [allMessages, setAllMessages] = useState<{ [key: string]: Message[] }>({});

    const [circles, setCircles] = useState<Circle[]>([]);
    const { userID, username, fetchUserData, handleDelete } = useFetchDashboard(setCircles);
    const { sendMessage } = useWebSocketDashboard(setCircles, allMessages, setAllMessages);

    useEffect(() => {
        const fetchData = async () => {
            await fetchUserData();
            setIsUserDataFetched(true);
        };
        fetchData();
    }, []);

    useEffect(() => {
        const chatContainer = document.querySelector(".chat-placeholder") as HTMLElement;
        if (chatContainer) {
            const textBoxHeight = window.innerHeight * 0.94;
            chatContainer.scrollTop = chatContainer.scrollHeight - textBoxHeight;
        }
    }, [allMessages[selectedCircleID]]);

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

    // Edit users in circle
    const [openEditModal, setOpenEditModal] = useState(false);
    function handleOpenEditModal() {
        setOpenEditModal(true);
    }
    const handleEditClose = (event: HandleCloseEvent) => {
        if (event.target.classList.contains('modal-container')) {
            setOpenEditModal(false);
        }
    };

    const handleEnter = (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (event.key === "Enter") {
            const content = (event.target as HTMLInputElement).value;
            if (content.trim() === "") return;
            const message: Message = {
                circle_id: selectedCircleID,
                content: content,
                created_at: new Date().toISOString(),
                author_id: userID,
            };
            const webSocketMessage: WebSocketMessage = {
                origin: "client",
                type: "message",
                action: "create",
                message: message,
            };
            sendMessage(webSocketMessage);
            setAllMessages((prevMessages) => {
                const messages = prevMessages[selectedCircleID] || [];
                return {
                    ...prevMessages,
                    [selectedCircleID]: [...messages, message],
                };
            });
            (event.target as HTMLInputElement).value = "";
        }
    }


    return (
        <>
            {!isUserDataFetched ?
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
                                    {username[0]}
                                </span>
                            </div>

                            <div>
                                <div className="px-2">
                                    <ul className="left-column space-y-1 border-t pt-4" style={{ borderColor: 'hsl(0, 0%, 87%)' }}>
                                        {circles.map((circle) => (
                                            <li key={circle.id} className="flex justify-center">
                                                <button onClick={() => setSelectedCircleID(circle.id)} className="w-full">
                                                    <a
                                                        className={`circles button-size group ${selectedCircleID === circle.id ? "selected-circle" : ""}`}
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
                                                            className="invisible absolute start-full top-1/2 ms-5 -translate-y-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible whitespace-nowrap"
                                                        >
                                                            {circle.name}
                                                        </span>
                                                    </a>
                                                </button>
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            </div>
                        </div>

                        <div className="left-column sticky inset-x-0 bottom-0 border-t bg-white ml-2 mr-2 pt-3 pb-3" style={{ borderColor: 'hsl(0, 0%, 87%)' }}>
                            <button
                                className="button-size group"
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
                                    className="invisible absolute start-full top-1/2 ms-5 -translate-y-1/2 rounded bg-gray-900 px-1.5 py-1.5 text-xs font-medium text-white group-hover:visible whitespace-nowrap"
                                >
                                    Add Circle
                                </span>
                            </button>
                            <button
                                className="button-size group"
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
                                    className="invisible absolute start-full top-1/2 ms-5 -translate-y-1/2 rounded bg-gray-900 px-1.5 py-1.5 text-xs font-medium text-white group-hover:visible"
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
                        {openEditModal && <div className="modal-container" onClick={handleEditClose}>
                            <div className="modal inset-0 bg-black bg-opacity-50 z-50">
                                <EditModal isOpen={openEditModal} setOpen={() => setOpenEditModal(false)} circleId={selectedCircleID} />
                            </div>
                        </div>}
                    </div>
                    <div className="chat-container">
                        <div className="chat-title">
                            <h1>{circles.find((circle) => circle.id === selectedCircleID)?.name}&nbsp;</h1>
                            <div className="chat-menu">
                                <button
                                    className="button-size group disabled:opacity-50 disabled:cursor-not-allowed"
                                    onClick={selectedCircleID !== "" ? handleOpenEditModal : undefined}
                                    disabled={selectedCircleID === ""}
                                >
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        className="size-4 opacity-75"
                                        fill="currentColor"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                        strokeWidth="1"
                                    >
                                        <path d="M3.5,24h15A3.51,3.51,0,0,0,22,20.487V12.95a1,1,0,0,0-2,0v7.537A1.508,1.508,0,0,1,18.5,22H3.5A1.508,1.508,0,0,1,2,20.487V5.513A1.508,1.508,0,0,1,3.5,4H11a1,1,0,0,0,0-2H3.5A3.51,3.51,0,0,0,0,5.513V20.487A3.51,3.51,0,0,0,3.5,24Z"></path>
                                        <path d="M9.455,10.544l-.789,3.614a1,1,0,0,0,.271.921,1.038,1.038,0,0,0,.92.269l3.606-.791a1,1,0,0,0,.494-.271l9.114-9.114a3,3,0,0,0,0-4.243,3.07,3.07,0,0,0-4.242,0l-9.1,9.123A1,1,0,0,0,9.455,10.544Zm10.788-8.2a1.022,1.022,0,0,1,1.414,0,1.009,1.009,0,0,1,0,1.413l-.707.707L19.536,3.05Zm-8.9,8.914,6.774-6.791,1.4,1.407-6.777,6.793-1.795.394Z"></path>
                                    </svg>
                                    <span
                                        className="invisible absolute start-1/2 top-full mt-5 -translate-x-1/2 rounded bg-gray-900 px-1.5 py-1.5 text-xs font-medium text-white group-hover:visible"
                                    >
                                        Edit Users
                                    </span>
                                </button>
                                <button
                                    className="button-size group disabled:opacity-50 disabled:cursor-not-allowed"
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
                                        className="invisible absolute start-1/2 top-full mt-5 -translate-x-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible"
                                    >
                                        Invite Users
                                    </span>
                                </button>
                                <button
                                    className="button-size group disabled:opacity-50 disabled:cursor-not-allowed"
                                    onClick={selectedCircleID !== "" ? () => handleDelete(selectedCircleID) : undefined}
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
                                        className="invisible absolute start-1/2 top-full mt-5 -translate-x-1/2 rounded bg-gray-900 px-1.5 py-1.5 text-xs font-medium text-white group-hover:visible"
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
                </div>}
        </>

    )
};

