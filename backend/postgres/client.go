package postgres

import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var pool *pgxpool.Pool

func ConnectPSQL() {
	var err error
	err = godotenv.Load()
    if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
    }

	db_url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("PSQL_USER"),
		os.Getenv("PSQL_PASSWORD"),
		os.Getenv("PSQL_HOST"),
		os.Getenv("PSQL_PORT"),
		os.Getenv("PSQL_DBNAME"),
		os.Getenv("PSQL_SSLMODE"),
	)
	println(db_url)

	config, err := pgxpool.ParseConfig(db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse database URL: %v\n", err)
		os.Exit(1)
	}

	config.MaxConns = 10 

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Successfully connected to database")
}

func CreateAccount(data user.RegisterData) bool {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
		os.Exit(1)
	}
	defer conn.Release()

	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)",
		data.FirstName,
		data.LastName,
		data.Email,
		data.Password,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert into database: %v\n", err)
		return false
	}
	fmt.Println("Successfully created account")
	return true
}

func FindAccount(data user.LoginData) bool {
	var email string
	var password string
	err := pool.QueryRow(
		context.Background(),
		"SELECT email, password FROM users WHERE email = $1",
		data.Email,
	).Scan(&email, &password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find account: %v\n", err)
		return false
	}
	if email == data.Email && password == data.Password {
		return true
	}
	return false
}

func InsertRefreshToken(email string, token uuid.UUID) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
		os.Exit(1)
	}
	defer conn.Release()

	var userID int
	err = conn.QueryRow(
		context.Background(),
		"SELECT id FROM users WHERE email = $1",
		email,
	).Scan(&userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to select user ID: %v\n", err)
	}
	
	expiryDate := time.Now().AddDate(0, 0, 30).UTC()
	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO refresh_tokens (user_id, token, expires) VALUES ($1, $2, $3)",
		userID,
		token,
		expiryDate,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert into database: %v\n", err)
	}
}

func ClosePSQL() {
	if pool != nil {
		pool.Close()
		fmt.Println("Database connection closed")
	}
}