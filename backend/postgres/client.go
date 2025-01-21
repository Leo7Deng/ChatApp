package postgres

import (
	"context"
	"fmt"
	"os"
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

func GetPool() *pgxpool.Pool {
	return pool
}

func ClosePSQL() {
	if pool != nil {
		pool.Close()
		fmt.Println("PSQL connection closed")
	}
}