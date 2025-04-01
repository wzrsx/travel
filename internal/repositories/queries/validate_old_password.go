package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ValidateOldPassword(pool *pgxpool.Pool, user_id int, oldPassword string) (string, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return "", err
	}
	defer conn.Release()

	var email string
	err = conn.QueryRow(context.Background(), "SELECT email FROM users WHERE id_user = $1 AND password = $2", user_id, oldPassword).Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}
