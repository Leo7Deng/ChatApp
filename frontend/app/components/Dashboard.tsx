"use client"

import "./Dashboard.css"
import React, { useEffect, useRef, useState } from "react";
import CreateCircleModal from "./CreateCircleModal";

interface Circle {
    id: number;
    name: string;
    created_at: string;
}

export default function Dashboard() {
    const [circles, setCircles] = useState<Circle[]>([]);


    useEffect(() => {
        async function fetchUserData() {
            const headers = {
                'Content-Type': 'application/json',
            };
            fetch('http://localhost:8000/api/dashboard', {
                method: 'POST',
                headers: headers,
                credentials: 'include'
            })
                .then(response => response.json())
                .then(data => {
                    console.log("Data:", data);
                    const mappedCircles = data.map((circle: any) => ({
                        id: circle.id,
                        name: circle.name,
                        created_at: circle.created_at,
                    }));
                    setCircles(mappedCircles);
                })
                .catch(error => {
                    console.log(error);
                })
        }
        fetchUserData();
    }, []);


    const [openModal, setOpenModal] = useState(false);
    function createCircle() {
        setOpenModal(true);
    }

    const [selectedCircleID, setSelectedCircleID] = useState(0);

    interface HandleCloseEvent extends React.MouseEvent<HTMLDivElement> {
        target: EventTarget & HTMLDivElement;
    }

    const handleClose = (event: HandleCloseEvent) => {
        if (event.target.classList.contains('modal-container')) {
            setOpenModal(false);
        }
    };

    interface HandleEnterEvent extends React.KeyboardEvent<HTMLInputElement> {}
    const handleEnter = (event: HandleEnterEvent) => {
        if (event.key === 'Enter') {
            console.log("Enter key pressed");
        }
    };

    const ws = useRef<WebSocket>();
    useEffect(() => {
        ws.current = new WebSocket("ws://localhost:8000/ws");
        ws.current.onopen = () => console.log("ws opened");
        ws.current.onclose = () => console.log("ws closed");
        const wsCurrent = ws.current;
        return () => {
            wsCurrent.close();
        };
    }, []);

    useEffect(() => {
        if (!ws.current) return;
        ws.current.onmessage = (e) => {
            console.log("ws message: ", e.data);
            try {
                const parsedData = JSON.parse(e.data);
                console.log("Parsed data:", parsedData);
                const newCircle = {
                    id: parsedData.id,
                    name: parsedData.name,
                    created_at: parsedData.created_at,
                };
                setCircles((prevCircles) => [...prevCircles, newCircle]);
            } catch (error) {
                console.error("Error parsing WebSocket message:", error);
            }
        };
    }, []);

    return (
        <div className="dashboard">
            <div className="flex h-screen w-16 flex-col justify-between border-e bg-white sidebar">
                <div>
                    <div className="inline-flex size-16 items-center justify-center">
                        <span className="grid size-10 place-content-center rounded-lg bg-gray-100 text-xs text-gray-600">
                            L
                        </span>
                    </div>

                    <div className="border-t border-gray-100">
                        <div className="px-2">
                            <div className="py-4">
                                <a
                                    href="#"
                                    className="t group relative flex justify-center rounded bg-blue-50 px-2 py-1.5 text-blue-700"
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
                                        className="invisible absolute start-full top-1/2 ms-4 -translate-y-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible"
                                    >
                                        General
                                    </span>
                                </a>
                            </div>

                            <ul className="space-y-1 border-t border-gray-100 pt-4">
                                {circles.map((circle) => (
                                    <li key={circle.id}>
                                        <button onClick={() => setSelectedCircleID(circle.id)}>
                                            <a
                                                href="#"
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
                                        href="#"
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
                                                d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"
                                            />
                                        </svg>

                                        <span
                                            className="invisible absolute start-full top-1/2 ms-4 -translate-y-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible"
                                        >
                                            Billing
                                        </span>
                                    </a>
                                </li>

                                <li>
                                    <a
                                        href="#"
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
                                                d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"
                                            />
                                        </svg>

                                        <span
                                            className="invisible absolute start-full top-1/2 ms-4 -translate-y-1/2 rounded bg-gray-900 px-2 py-1.5 text-xs font-medium text-white group-hover:visible"
                                        >
                                            Invoices
                                        </span>
                                    </a>
                                </li>

                                <li>
                                    <a
                                        href="#"
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
                                        className="group relative flex justify-center rounded px-2 py-1.5 text-gray-500 hover:bg-gray-50 hover:text-gray-700"
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

                <div className="sticky inset-x-0 bottom-0 border-t border-gray-100 bg-white p-2">
                    <form action="#">
                        <button
                            type="submit"
                            className="group relative flex w-full justify-center rounded-lg px-2 py-1.5 text-sm text-gray-500 hover:bg-gray-50 hover:text-gray-700"
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
                    </form>
                </div>
                {openModal && <div className="modal-container" onClick={handleClose}>
                    <div className="modal inset-0 bg-black bg-opacity-50 z-50">
                        <CreateCircleModal isOpen={openModal} setOpen={() => setOpenModal(false)} />
                    </div>
                </div>}
            </div>
            <div className="chat-container">
                <div className="chat-placeholder">
                    {selectedCircleID !== 0 && (
                        <div>
                            <h2>Chat for Circle {selectedCircleID}</h2>
                        </div>
                    )}
                </div>
                <div className="text-box">
                    <input type="text" className="text-input" placeholder="Type a message" onKeyDown={handleEnter} />
                </div>
            </div>
            <div className="analytics-container">
                <div className="analytics-placeholder">
                </div>
            </div>
        </div>
    )
};