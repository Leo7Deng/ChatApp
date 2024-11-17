import client from "../db_config/redisConfig";

interface RedisService {
    createSession(sessionId: string, userId: string, ttl: number): Promise<void>;
    authenticateSession(sessionId: string): Promise<boolean>;
    deleteSession(sessionId: string): Promise<void>;
}

interface RedisObject {
    session_id: string;
    user_id: string;
    time_expire: Date;
}

const redisService: RedisService = {
    async createSession(sessionId, userId, ttl) {
        const currentTime = new Date();
        const redisObject: RedisObject = {
            session_id: sessionId,
            user_id: userId,
            time_expire: new Date(currentTime.getTime() + ttl),
        }
        
    }
};