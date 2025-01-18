"use client"

import { cookies } from "next/headers";
import "./Dashboard.css"
import React, { useEffect, useState } from "react";

export default function Dashboard() {
    function getCookie(name : string) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        if (parts.length === 2) {
            const part = parts.pop();
            if (part) {
                console.log(part.split(';').shift());
                return part.split(';').shift();
            }
        }
        return
      }
    
    useEffect(() => {
        async function fetchUserData() {
            // const access_token = getCookie("access-token");
            // if (!access_token) {
            //     console.log("Failed to get cookie");
            //     return;
            // }
            const headers = {
                'Content-Type': 'application/json',
                // 'Authorization': access_token
            };
            fetch('http://localhost:8000/api/dashboard', {
                method: 'POST',
                headers: headers,
                credentials: 'include'
            })
        }
        fetchUserData();
    }, []);

    return (
        <p>Currently logged in user: </p>
    );
}