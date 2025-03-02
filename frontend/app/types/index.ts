export interface Circle {
    id: string;
    name: string;
    created_at: string;
}

export interface WebSocketMessage {
    origin: "server" | "client";
    type: "message" | "circle";
    action: "create" | "delete";
    message?: Message;
    circle?: Circle;
}

export interface Message {
    circle_id: string;
    content: string;
    created_at: string;
    author_id: string;
}

export interface User {
    id: string;
    username: string;
    checked: boolean;
}

export interface EditUser {
    id: string;
    username: string;
    role: string;
}