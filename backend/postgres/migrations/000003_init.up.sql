CREATE TABLE circles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_circles (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    circle_id INT NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role VARCHAR(50) DEFAULT 'member',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (circle_id) REFERENCES circles(id) ON DELETE CASCADE,
    UNIQUE (user_id, circle_id)
);

CREATE INDEX idx_users_circles_user_id
  ON users_circles(user_id);

CREATE INDEX idx_users_circles_circle_id
  ON users_circles(circle_id);
