import React from 'react';

interface CreateCircleModalProps {
    isOpen: boolean;
    setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

function CreateCircleModal({ isOpen, setOpen }: CreateCircleModalProps) {
    interface CircleData {
        name: string;
    }

    const createCircle = (circleName: string): void => {
        if (circleName === '') return;
        const data: CircleData = {
            name: circleName
        };
        fetch('http://localhost:8000/api/create-circle', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
            credentials: 'include',
        });
    };

    if (!isOpen) return null;

    return (
        <div className="mx-auto max-w-screen-xl px-4 py-4 sm:px-6 lg:px-8 relative z-10 focus:outline-none">
            <form action="#" className="mx-auto mb-0 mt-8 max-w-md space-y-4">
                <div>
                    <div className="relative">
                        <input
                            type="text"
                            className="w-full rounded-lg border-gray-200 p-4 text-sm shadow-sm text-black"
                            placeholder="Circle Name"
                        />
                    </div>
                </div>

                <div className="flex items-center justify-between">
                    <button
                        type="submit"
                        className="inline-block rounded-lg bg-blue-500 px-5 py-3 text-sm font-medium text-white"
                        onClick={(e) => {
                            e.preventDefault();
                            setOpen(false);
                            createCircle(document.querySelector('input')?.value || '');
                        }}
                    >
                        Create
                    </button>
                </div>
            </form>
        </div>
    );
}

export default CreateCircleModal;