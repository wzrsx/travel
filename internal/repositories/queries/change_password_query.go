package queries

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ChangePasswordQuery(p *pgxpool.Pool, email string, password string) error {
	conn, err := p.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	log.Printf("psw: %s, email: %s", password, email)
	_, err = conn.Exec(context.Background(), "UPDATE users SET password = $1 WHERE email = $2", password, email)
	if err != nil {
		return fmt.Errorf("Error changing password user: %v", err)
	}

	return nil
}
