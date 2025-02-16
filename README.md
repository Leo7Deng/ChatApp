# Introduction

Chat App is a real time chat application designed for performance and scalability. 

### Motivation

The motivation behind this project was to learn distributed system design and build a scalable architecture that can handle a large number of users and messages.
I wanted to explore new technologies like Go, WebSockets, and Apache Kafka, and develop a real-time data pipeline that supports machine learning (ML) training.

# Tech Stack

| Component         | Technology Used            |
|------------------|--------------------------|
| **Frontend**     | (React/Next.js, TypeScript) |
| **Backend**      | Go    |
| **Database**     | PostgreSQL (relational) + Redis (NoSQL) + xxx (NoSQL) |
| **Caching**      | Redis                      |
| **Real-Time**    | WebSockets (gorilla/websocket for Go) |    
| **Message Queue** |Apache Kafka                |
| **Deployment**   | Docker Compose             |

# Design 

### Authentication

I chose to not use third-party authentication libraries. Instead, I designed a refresh token and JWT-based access token system:
- The refresh token is securely stored in an httpOnly cookie and managed as a key in Redis for fast retrieval and easy revocation.
- The access token is a JWT that includes the user ID as a claim. It is stored in memory and has a short lifespan for improved security.
- Using authContext in React, I share the access token across all child components and centralize the token refresh logic.

### Websockets

I used the gorilla/websocket library for Go to manage WebSocket connections:
- Created a hub that manages all WebSocket connections and broadcasts messages.
- Mapped user IDs to each WebSocket client to support direct messaging.
- To scale horizontally, I plan to use Kafka consumers to broadcast messages across multiple WebSocket servers.

### Kafka

I chose Apache Kafka because of its high throughput and fault-tolerant design. It also supports multiple consumers, which is ideal for adding new features in the future.

<!-- # Difficulties
- Designing how to handle errors in Kafka
- Designing how to metigate eventual consistency in xxx -->