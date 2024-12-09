package postgres

import (
	"fmt"
	"os"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Leo7Deng/ChatApp/models"
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

func CreateAccount(data user.FormData) {
	// Acquire a connection from the pool
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
		os.Exit(1)
	}
	fmt.Println("Successfully created account")
}

func ClosePSQL() {
	if pool != nil {
		pool.Close()
		fmt.Println("Database connection closed")
	}
}