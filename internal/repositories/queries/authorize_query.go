package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)
type UserAuthResult struct {
    UserID   int
    Username string
}
func NewUserAuthResult() UserAuthResult{
	var ur UserAuthResult
	return ur
}

func (ur *UserAuthResult) AuthorizeQuery(email string, password string, pool *pgxpool.Pool) (error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(), "SELECT id, username FROM users WHERE email = $1 AND password = $2", email, password).Scan(&ur.UserID, &ur.Username)
	if err != nil {
		return err
	}
	return nil
}
