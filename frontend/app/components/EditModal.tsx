import React, { useEffect, useState } from 'react';
import { useAuth } from "../context/authContext";
import { EditUser } from "../types";
import "./EditModal.css"

interface InviteModalProps {
    isOpen: boolean;
    setOpen: React.Dispatch<React.SetStateAction<boolean>>;
    circleId: string;
}

function EditModal({ isOpen, setOpen, circleId }: InviteModalProps) {
    const authContext = useAuth();
    if (!authContext) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    const { getAccessToken } = authContext;
    const [users, setUsers] = useState<EditUser[]>([]);
    
    useEffect(() => {
        async function fetchEditUsers() {
            const token = await getAccessToken();
            const headers = {
                'Authorization': `Bearer ${token}`,
            };
            const body = {
                circle_id: circleId,
            };
            fetch('http://localhost:8000/api/circles/edit', {
                method: 'POST',
                headers: headers,
                body: JSON.stringify(body),
            })
                .then(async (response) => {
                    const data = await response.json();
                    if (!response.ok) {
                        console.log("Error:", data);
                    } else {
                        console.log("Data:", data);
                        const mappedUsers = data.map((user: any) => ({
                            id: user.id,
                            username: user.username,
                            role: user.role,
                        }));
                        setUsers(mappedUsers);
                    }
                })
                .catch(error => {
                    console.log(error);
                });
        }
        fetchEditUsers();
    }, [circleId, getAccessToken]);
    
    const handleRoleChange = (userId: string, newRole: string) => {
        setUsers(prevUsers => 
            prevUsers.map(user => 
                user.id === userId ? { ...user, role: newRole } : user
            )
        );
    };

    async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        const token = await getAccessToken();
        const headers = {
            'Authorization': `Bearer ${token}`,
        };
        const body = {
            circle_id: circleId,
            users: users.map(user => ({ id: user.id, role: user.role })),
        };
        console.log(body);
        fetch('http://localhost:8000/api/circles/edit', {
            method: 'PUT',
            headers: headers,
            body: JSON.stringify(body),
        })
            .then(async (response) => {
                const data = await response.json();
                if (!response.ok) {
                    console.log("Error:", data);
                } else {
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
        <div className="mx-auto max-w-screen-xl relative z-10 focus:outline-none">
            <form action="#" className="mx-auto mb-4 mt-6 max-w-md space-y-4" onSubmit={handleSubmit}>
                <div>
                    <div className="search w-full relative rounded-md">
                        <input
                            className="search-input w-full border-gray-200 p-2 text-sm shadow-sm text-black"
                            placeholder="Search username"
                            id="circleName"
                        />
                        <svg xmlns="http://www.w3.org/2000/svg"
                            className="size-11 search-icon hover:cursor-pointer"
                            fill="hsl(0, 0%, 100%)"
                            viewBox="0 0 24 24"
                            stroke="hsl(0, 0%, 100%)"
                            strokeWidth="1">
                            <path d="M 9 2 C 5.1458514 2 2 5.1458514 2 9 C 2 12.854149 5.1458514 16 9 16 C 10.747998 16 12.345009 15.348024 13.574219 14.28125 L 14 14.707031 L 14 16 L 20 22 L 22 20 L 16 14 L 14.707031 14 L 14.28125 13.574219 C 15.348024 12.345009 16 10.747998 16 9 C 16 5.1458514 12.854149 2 9 2 z M 9 4 C 11.773268 4 14 6.2267316 14 9 C 14 11.773268 11.773268 14 9 14 C 6.2267316 14 4 11.773268 4 9 C 4 6.2267316 6.2267316 4 9 4 z"></path>
                        </svg>
                    </div>
                    <div className="flex users my-4 rounded-md gap-3 shadow-sm">
                        <ul>
                            {users.length === 0 && (
                                <li className="user-item flex items-center gap-2">
                                    <p className="font-medium text-sm text-gray-500">You don't have admin access to this circle</p>
                                </li>
                            )}
                            {users.length != 0 && users.map(user => (
                                <li key={user.id} className="user-item flex items-center gap-2">
                                    <p className="font-medium text-sm text-gray-500 username">{user.username}</p>
                                    <select
                                        value={user.role}
                                        onChange={(e) => handleRoleChange(user.id, e.target.value)}
                                        className="border border-gray-300 p-1 rounded select-role"
                                    >
                                        <option value="member">Member</option>
                                        <option value="admin">Admin</option>
                                    </select>
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
    
                <div className="flex items-center justify-between">
                    <button
                        type="submit"
                        className="inline-block rounded-lg bg-blue-500 px-5 py-3 text-sm font-medium text-white"
                    >
                        Submit
                    </button>
                </div>
            </form>
        </div>
    );
}

export default EditModal;