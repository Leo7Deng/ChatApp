import { createClient, RedisClientType } from 'redis';

const client: RedisClientType = createClient({
    url: process.env.REDIS_URL || 'redis://localhost:6379',
    socket: {
        reconnectStrategy: (retries: number): number => {
            const jitter = Math.floor(Math.random() * 200); // Random jitter (0-200ms)
            const delay = Math.min(Math.pow(2, retries) * 50, 2000); // Exponential backoff
            return delay + jitter;
        },
    },
});

// Event listeners for connection handling
client.on('connect', () => {
    console.log('Connected to Redis');
});

client.on('error', (err: Error) => {
    console.error('Redis connection error:', err.message);
});

// Export the client for reuse across the application
export default client;
