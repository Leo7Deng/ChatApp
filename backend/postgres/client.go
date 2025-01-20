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

var pool *pgxpool.Pool

func ConnectPSQL() {
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
		fmt.Fprintf(os.Stderr, "Unable to parse PSQL URL: %v\n", err)
		os.Exit(1)
	}

	config.MaxConns = 10 

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	
	
	fmt.Println("Successfully connected to PSQL")
}

func CreateAccount(data models.RegisterData) (string, error) {
	// conn, err := pool.Acquire(context.Background())
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer conn.Release()
	// _, err := conn.Exec(

	var userID string
	err := pool.QueryRow(
		context.Background(),
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id",
		data.FirstName,
		data.LastName,
		data.Email,
		data.Password,
	).Scan(&userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert into PSQL: %v\n", err)
		return "", err
	}
	fmt.Println("Successfully created account")
	return userID, nil
}

func FindAccount(data models.LoginData) (string, error) {
	var id string
	var email string
	var password string
	err := pool.QueryRow(
		context.Background(),
		"SELECT id, email, password FROM users WHERE email = $1",
		data.Email,
	).Scan(&id, &email, &password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find account: %v\n", err)
		return "", err
	}
	if email == data.Email && password == data.Password {
		return id, nil
	}
	return "", err
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


func ClosePSQL() {
	if pool != nil {
		pool.Close()
		fmt.Println("PSQL connection closed")
	}
}