"use client"

import React from 'react';
import Dashboard from "../components/Dashboard";
import { AuthProvider } from "../context/authContext";

export default function Home() {
  return (
      <AuthProvider>
        <Dashboard />
      </AuthProvider>
  );
}
