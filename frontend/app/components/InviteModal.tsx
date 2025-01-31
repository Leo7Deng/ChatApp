import React, { useState } from 'react';
import "./InviteModal.css"

interface InviteModalProps {
    isOpen: boolean;
    setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

function InviteModal({ isOpen, setOpen }: InviteModalProps) {
    const [inviteAll, setInviteAll] = useState(false);

    if (!isOpen) return null;

    return (
        <div className="mx-auto max-w-screen-xl px-4 py-2 sm:px-6 lg:px-8 relative z-10 focus:outline-none">
            <form action="#" className="mx-auto mb-4 mt-6 max-w-md space-y-4">
                <div>
                    <div className="relative">
                        <input
                            type="text"
                            className="w-full rounded-md border-gray-200 p-2 text-sm shadow-sm text-black"
                            placeholder="Search username"
                            id="circleName"
                        />
                    </div>
                    <div className="flex users my-4 rounded-md p-2 gap-3">
                        <ul>
                            <li>user 1</li>
                            <li>user 2</li>
                            <li>user 3</li>
                            <li>user 4</li>
                            <li>user 5</li>
                            <li>user 1</li>
                            <li>user 2</li>
                            <li>user 3</li>
                            <li>user 4</li>
                            <li>user 5</li>
                            <li>user 1</li>
                            <li>user 2</li>
                            <li>user 3</li>
                            <li>user 4</li>
                            <li>user 5</li>
                        </ul>
                    </div>
                </div>

                <div className="flex items-center justify-between">
                    <button
                        type="submit"
                        className="inline-block rounded-lg bg-blue-500 px-5 py-3 text-sm font-medium text-white"
                        onClick={(e) => {
                            e.preventDefault();
                            setOpen(false);
                        }}
                    >
                        Submit
                    </button>

                    <label htmlFor="Option1" className="flex cursor-pointer items-center gap-2">
                        <div className="flex items-center">
                            &#8203;
                            <input
                                type="checkbox"
                                className="size-5 rounded border-gray-300 text-blue-500 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50"
                                checked={inviteAll} 
                                onChange={() => setInviteAll((prev) => !prev)} 
                            />
                        </div>
                        <div className="flex items-center">
                            <p className="font-medium text-sm text-gray-500">Invite all users</p>
                        </div>
                    </label>
                </div>
            </form>
        </div>
    );
}
export default InviteModal;