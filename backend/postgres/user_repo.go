package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Leo7Deng/ChatApp/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateAccount(data models.RegisterData) (string, error) {
	var userID string
	err := r.pool.QueryRow(
		context.Background(),
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id",
		data.FirstName,
		data.LastName,
		data.Email,
		data.Password,
	).Scan(&userID)
	if err != nil {
		fmt.Printf("Unable to insert into PSQL: %v\n", err)
		return "", err
	}
	fmt.Println("Successfully created user")
	return userID, nil
}

func (r *UserRepository) FindAccount(email string) (*models.User, error) {
	var user models.User
	err := r.pool.QueryRow(
		context.Background(),
		"SELECT id, first_name, last_name, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		fmt.Printf("Unable to find user: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func InsertRefreshToken(userID string, token uuid.UUID) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
		os.Exit(1)
	}
	defer conn.Release()

	expiryDate := time.Now().AddDate(0, 0, 30).UTC()
	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO refresh_tokens (user_id, token, expires) VALUES ($1, $2, $3)",
		userID,
		token,
		expiryDate,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert into PSQL: %v\n", err)
	}
}

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

func CreateCircle(userID string, name string) error {
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

	var circleID int
	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO circles (name, created_at)
		VALUES ($1, NOW())
		RETURNING id;		
		`,
		name).Scan(&circleID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting into circles: %v\n", err)
		return err
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
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to commit transaction: %v\n", err)
		return err
	}

	fmt.Printf("Circle '%s' created successfully with ID %d and associated with user ID %s\n", name, circleID, userID)
	return nil
}
