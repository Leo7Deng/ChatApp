package postgres

import (
	"context"
	"fmt"
	"github.com/Leo7Deng/ChatApp/models"
	"os"
	"time"
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
		return fmt.Errorf("user is not admin of circle")
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
