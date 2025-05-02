CREATE OR REPLACE FUNCTION search_circle_messages(circle_id_input INT, content_input TEXT)
RETURNS TABLE (
    content TEXT,
    created_at TIMESTAMP,
    author_username VARCHAR(100)
) AS $$
BEGIN
    RETURN QUERY
    SELECT m.content, m.created_at, u.username
    FROM messages m
    INNER JOIN users u ON m.author_id = u.id
    WHERE m.circle_id = circle_id_input AND m.content ILIKE '%' || content_input || '%';
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_users_in_circle(request_user_id INT, circle_id_input INT)
RETURNS TABLE (
    user_id INT,
    username VARCHAR(100),
    role VARCHAR(50)
) AS $$
BEGIN
    IF NOT EXISTS (
        SELECT *
        FROM users_circles uc
        WHERE uc.user_id = request_user_id
          AND uc.circle_id = circle_id_input
          AND uc.role = 'admin'
    ) THEN
        RAISE EXCEPTION 'permission denied';
    END IF;

    RETURN QUERY
    SELECT u.id, u.username, uc.role
    FROM users u
    INNER JOIN users_circles uc ON u.id = uc.user_id
    WHERE u.id != request_user_id AND uc.circle_id = circle_id_input
    ORDER BY u.username ASC;
END;
$$ LANGUAGE plpgsql;
