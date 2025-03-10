package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Review struct {
	Username    string
	Description string
	Estimation  int
	Date        string
}

func CreateReviews(id_route int, pool *pgxpool.Pool) ([]Review, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
		SELECT username, description, estimation, date_review
		FROM reviews
		WHERE id_route = $1`, id_route)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Review

	for rows.Next() {
		var username string
		var description string
		var estimation int
		var date_review time.Time
		if err := rows.Scan(&username, &description, &estimation, &date_review); err != nil {
			return nil, fmt.Errorf("Error scan Reviews: %s", err.Error())
		}
		reviews = append(reviews, Review{
			Username:    username,
			Description: description,
			Estimation:  estimation,
			Date:        date_review.Format("2006-01-02"),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil

}
