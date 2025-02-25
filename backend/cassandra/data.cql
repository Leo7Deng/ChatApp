CREATE KEYSPACE IF NOT EXISTS chat_app WITH REPLICATION = 
{ 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

CREATE TABLE IF NOT EXISTS chat_app.messages (
    circle_id TEXT,         -- Partition key (groups messages by circle)
    created_at TIMESTAMP,   -- Clustering key (orders messages by time)
    message_id UUID,        -- Unique ID for the message
    author_id TEXT,         -- ID of the user who sent it
    content TEXT,           -- message
    PRIMARY KEY (circle_id, created_at, message_id)
) WITH CLUSTERING ORDER BY (created_at DESC);
