package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Leo7Deng/ChatApp/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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
		"INSERT INTO users (first_name, last_name, username, email, password) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		data.FirstName,
		data.LastName,
		data.Username,
		data.Email,
		data.Password,
	).Scan(&userID)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" { // 23505 = unique violation
			if strings.Contains(pgErr.Message, "users_email_key") {
				return "", fmt.Errorf("email")
			} else if strings.Contains(pgErr.Message, "users_username_key") {
				return "", fmt.Errorf("username")
			}
		}
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

