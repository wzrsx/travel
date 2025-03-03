package queries

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reviews struct{
	Description string
	Estimation int
}
func CreateReviews (id_route int, pool *pgxpool.Pool) ([]Reviews, error){
	conn, err := pool.Acquire(context.Background())
	if err != nil{
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
		SELECT description, estimation
		FROM reviews
		WHERE (SELECT id_review from routes WHERE id_route = $1)`, id_route)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Reviews

	for rows.Next() {
		var description string
		var estimation int
		if err := rows.Scan(&description, &estimation); err != nil {
			return nil, fmt.Errorf("Error scan Reviews: ", err)
		}
		reviews = append(reviews, Reviews{
			Description: description,
			Estimation:  estimation,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil

}