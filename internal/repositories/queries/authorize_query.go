package queries

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAuthResult struct {
	UserID   int
	Username string
}

func NewUserAuthResult() UserAuthResult {
	var ur UserAuthResult
	return ur
}

func (ur *UserAuthResult) AuthorizeQuery(email string, password string, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	//проверка есть ли email в бд
	var count int
    errCheckEmail := conn.QueryRow(context.Background(), "SELECT count(*) FROM users WHERE email = $1", email).Scan(&count)
    if errCheckEmail != nil {
        return errCheckEmail
    }

    if count == 0 {
        return errors.New("email not found") 
    }

	err = conn.QueryRow(context.Background(), "SELECT id_user, username FROM users WHERE email = $1 AND password = $2", email, password).Scan(&ur.UserID, &ur.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("pass invalid")
		}
		return err
	}
	return nil
}
