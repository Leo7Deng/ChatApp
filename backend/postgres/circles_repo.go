
package postgres

import (
	"context"
	"fmt"
	"os"
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

	var circleID int
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
		INSERT INTO users_circles (user_id, circle_id, joined_at)
		VALUES ($1, $2, NOW())
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

func DeleteCircle(circleID string) error {
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