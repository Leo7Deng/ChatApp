package postgres

import (
	"fmt"
	"os"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
  )

func CreateDBConnection() {
	db_url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
	os.Getenv("PSQL_USER"),
	os.Getenv("PSQL_PASSWORD"),
	os.Getenv("PSQL_HOST"),
	os.Getenv("PSQL_PORT"),
	os.Getenv("PSQL_DBNAME"),
)
	dbpool, err := pgxpool.New(context.Background(), db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}
