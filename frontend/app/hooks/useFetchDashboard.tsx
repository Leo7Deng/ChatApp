import { useState } from "react";
import { Circle } from "../types";
import { useAuth } from "../context/authContext";

export function useFetchDashboard(setCircles: React.Dispatch<React.SetStateAction<Circle[]>>) {
    const [userID, setUserID] = useState("");
    const [username, setUsername] = useState("");
    const authContext = useAuth();
    if (!authContext) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    const { getAccessToken } = authContext;

    // Fetch User Data
    async function fetchUserData() {
        const token = await getAccessToken();
        fetchCircleData();
        const headers = {
            "Authorization": `Bearer ${token}`,
        };
        fetch('http://localhost:8000/api/user', {
            method: 'GET',
            headers: headers,
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

    // Fetch Circle Data
    async function fetchCircleData() {
        const token = await getAccessToken();
        const headers = {
            "Authorization": `Bearer ${token}`,
        };
        fetch('http://localhost:8000/api/circles', {
            method: 'GET',
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
                }
            })
            .catch(error => {
                console.log(error);
            });
    }



    // Delete Circle
    async function handleDelete(circleID: string) {
        const token = await getAccessToken();
        try {
            const response = await fetch(`http://localhost:8000/api/circles/delete/${circleID}`, {
                method: 'DELETE',
                headers: {
                    "Authorization": `Bearer ${token}`,
                },
            });

            const data = await response.json();
            if (!response.ok) {
                if (data === "user is not admin of circle") {
                    alert("You are not an admin of the circle");
                }
                throw new Error("Failed to delete circle");
            }

            setCircles((prevCircles) => prevCircles.filter((circle) => circle.id !== circleID));
        } catch (error) {
            console.error("Error deleting circle:", error);
        }
    }

    return { userID, username, fetchUserData, handleDelete };
}

