package queries

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRegistrationResult struct {
	UserID   int
	Username string
}

func NewUserRegistrationResult() UserRegistrationResult {
	var ur UserRegistrationResult
	return ur
}

func (ur *UserRegistrationResult) RegistrationQuery(username string, email string, password string, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	// Проверка наличия почты в базе данных
	var exists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("email exists")
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id_user", username, email, password).Scan(&ur.UserID)
	if err != nil {
		return err
	}
	ur.Username = username

	return nil
}
