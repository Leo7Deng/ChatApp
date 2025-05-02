package postgres

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/Leo7Deng/ChatApp/models"
)

func GetUserCircles(userID string) ([]models.Circle, error) {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`
		SELECT c.id, c.name, c.created_at
		FROM circles c
		INNER JOIN users_circles uc ON c.id = uc.circle_id
		WHERE uc.user_id = $1;
		`,
		userID,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query circles: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var circles []models.Circle
	for rows.Next() {
		var circle models.Circle
		err := rows.Scan(&circle.ID, &circle.Name, &circle.CreatedAt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan circle row: %v\n", err)
			return nil, err
		}
		circles = append(circles, circle)
	}

	if rows.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error during rows iteration: %v\n", rows.Err())
		return nil, rows.Err()
	}
	return circles, nil
}

func CreateCircle(userID string, name string) (models.Circle, error) {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
		return models.Circle{}, err
	}
	defer conn.Release()

	// Create a transaction because two operations
	tx, err := conn.Begin(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to begin transaction: %v\n", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var circleID string
	var currentTime = time.Now()
	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO circles (name, created_at)
		VALUES ($1, $2)
		RETURNING id;		
		`,
		name,
		currentTime).Scan(&circleID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting into circles: %v\n", err)
		return models.Circle{}, err
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO users_circles (user_id, circle_id, joined_at, role)
		VALUES ($1, $2, NOW(), 'admin');
		`,
		userID,
		circleID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting into users_circles: %v\n", err)
		return models.Circle{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to commit transaction: %v\n", err)
		return models.Circle{}, err
	}

	circle := models.Circle{ID: circleID, Name: name, CreatedAt: currentTime}
	fmt.Printf("Circle created: %v\n", circle)
	return circle, nil
}

func DeleteCircle(userID string, circleID string) error {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
		return err
	}
	defer conn.Release()

	// Create a transaction because two operations
	tx, err := conn.Begin(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to begin transaction: %v\n", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Check if user is admin of circle
	var role string
	err = tx.QueryRow(
		ctx,
		`
		SELECT role
		FROM users_circles
		WHERE user_id = $1 AND circle_id = $2;
		`,
		userID,
		circleID).Scan(&role)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking user role: %v\n", err)
	}
	if role != "admin" {
		fmt.Fprintf(os.Stderr, "User is not admin of circle\n")
		return fmt.Errorf("permission error")
	}

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM circles
		WHERE id = $1
		`,
		circleID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting from circles: %v\n",
			err)
	}

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM users_circles
		WHERE circle_id = $1
		`,
		circleID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting from users_circles: %v\n",
			err)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to commit transaction: %v\n", err)
		return err
	}
	fmt.Printf("Circle deleted: %v\n", circleID)
	return nil
}

func GetInviteUsersInCircle(userID string, circleID string) ([]models.User, error) {
	var users []models.User
	rows, err := pool.Query(
		context.Background(),
		`
		SELECT DISTINCT u.id, u.username 
		FROM users u
		WHERE u.id != $1 
		AND NOT EXISTS (
			SELECT 1 
			FROM users_circles uc 
			WHERE u.id = uc.user_id 
			AND uc.circle_id = $2 
		)
		ORDER BY u.username ASC
		LIMIT 10;
		`,
		userID,
		circleID,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query PSQL: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan row: %v\n", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func AddUsersToCircle(circleID string, userIDs []string) error {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to begin transaction: %v\n", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	for _, id := range userIDs {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO users_circles (user_id, circle_id, joined_at, role)
			VALUES ($1, $2, NOW(), 'member');
			`,
			id,
			circleID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error inserting into users_circles: %v\n", err)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to commit transaction: %v\n", err)
		return err
	}

	return nil
}

func EditRoleInCircle(circleID string, targetUserID string, role string) error {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
	}
	defer conn.Release()

	if role == "remove" {
		_, err = conn.Exec(
			ctx,
			`
			DELETE FROM users_circles
			WHERE user_id = $1 AND circle_id = $2;
			`,
			targetUserID,
			circleID,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to query PSQL: %v\n", err)
			return err
		}
		fmt.Printf("Role removed\n")
		return nil
	} else {
		err = conn.QueryRow(
			ctx,
			`
		UPDATE users_circles
		SET role = $1
		WHERE user_id = $2 AND circle_id = $3
		RETURNING role;
		`,
			role,
			targetUserID,
			circleID,
		).Scan(&role)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to query PSQL: %v\n", err)
			return err
		}
		fmt.Printf("Role updated: %v\n", role)
		return nil
	}
}

func GetExistingUsersInCircle(userID string, circleID string) ([]models.UserRole, error) {
	rows, err := pool.Query(
		context.Background(),
		"SELECT * FROM get_users_in_circle($1, $2);",
		userID,
		circleID,
	)
	if err != nil {
		if err.Error() == "permission error" {
			fmt.Fprintf(os.Stderr, "User is not admin of circle\n")
			return nil, fmt.Errorf("permission error")
		}
		fmt.Fprintf(os.Stderr, "Unable to query get_existing_users_in_circle: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.UserRole
	for rows.Next() {
		var user models.UserRole
		if err := rows.Scan(&user.UserID, &user.Username, &user.Role); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan row: %v\n", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func LoadCircleUserMap() (map[string]map[string]bool, error) {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
		return nil, err
	}
	fmt.Println("LoadCircleUserMap")
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`
		SELECT user_id, circle_id
		FROM users_circles;
		`,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query PSQL: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	circleUsers := make(map[string]map[string]bool)
	for rows.Next() {
		var userID, circleID string
		err = rows.Scan(&userID, &circleID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan row: %v\n", err)
			return nil, err
		}
		if _, ok := circleUsers[circleID]; !ok {
			circleUsers[circleID] = make(map[string]bool)
		}
		circleUsers[circleID][userID] = true
	}
	return circleUsers, nil
}

func SearchCircle(circleID string, content string) ([]models.SearchMessage, error) {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		"SELECT * FROM search_circle_messages($1, $2)",
		circleID,
		content,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query search_circle_messages: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.SearchMessage
	var time time.Time
	for rows.Next() {
		var message models.SearchMessage
		err = rows.Scan(&message.Content, &time, &message.AuthorUsername)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan row: %v\n", err)
			return nil, err
		}
		message.CreatedAt = time.Format("2006-01-02 15:04:05 EST")
		messages = append(messages, message)
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt > messages[j].CreatedAt
	})

	return messages, nil
}
