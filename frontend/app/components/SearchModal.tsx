import React, { use, useEffect, useState } from 'react';
import { useAuth } from "../context/authContext";
import "./SearchModal.css"
import { text } from 'stream/consumers';

interface SearchModalProps {
    isOpen: boolean;
    setOpen: React.Dispatch<React.SetStateAction<boolean>>;
    circleId: string;
}

function SearchModal({ isOpen, setOpen, circleId }: SearchModalProps) {
    const authContext = useAuth();
    if (!authContext) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    const { getAccessToken } = authContext;

    async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        const token = await getAccessToken();
        const headers = {
            'Authorization': `Bearer ${token}`,
        };
        const body = {
            circle_id: circleId,
            content: (e.currentTarget.elements.namedItem('circleName') as HTMLInputElement).value,
        };
        console.log(body);
        fetch('http://localhost:8000/api/circles/search', {
            method: 'POST',
            headers: headers,
            body: JSON.stringify(body),
        })
            .then(async (response) => {
                const data = await response.json();
                if (!response.ok) {
                    console.log("Error:", data);
                }
                else {
                    console.log("Data:", data);
                    setOpen(false);
                }
            })
            .catch(error => {
                console.log(error);
            });
    };

    if (!isOpen) return null;

    return (
        <div className="search-modal mx-auto max-w-screen-xl relative z-10 focus:outline-none">
            <form action="#" className="mx-auto mb-4 mt-6 max-w-md space-y-4" onSubmit={handleSubmit}>
                <div className="search w-full relative rounded-md">
                    <input
                        className="search-input w-full border-gray-200 p-2 text-sm shadow-sm text-black"
                        placeholder="Search text"
                        id="circleName"
                    />
                    <button type="submit">
                        <svg xmlns="http://www.w3.org/2000/svg"
                            className="size-11 search-icon hover:cursor-pointer"
                            fill="hsl(0, 0%, 100%)"
                            viewBox="0 0 24 24"
                            stroke="hsl(0, 0%, 100%)"
                            strokeWidth="1">
                            <path d="M 9 2 C 5.1458514 2 2 5.1458514 2 9 C 2 12.854149 5.1458514 16 9 16 C 10.747998 16 12.345009 15.348024 13.574219 14.28125 L 14 14.707031 L 14 16 L 20 22 L 22 20 L 16 14 L 14.707031 14 L 14.28125 13.574219 C 15.348024 12.345009 16 10.747998 16 9 C 16 5.1458514 12.854149 2 9 2 z M 9 4 C 11.773268 4 14 6.2267316 14 9 C 14 11.773268 11.773268 14 9 14 C 6.2267316 14 4 11.773268 4 9 C 4 6.2267316 6.2267316 4 9 4 z"></path>
                        </svg>
                    </button>
                </div>
            </form>
        </div>
    );
}
export default SearchModal;