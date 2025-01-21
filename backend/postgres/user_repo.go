package postgres

import (
	"context"
	"fmt"
	"github.com/Leo7Deng/ChatApp/models"
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
