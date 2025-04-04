package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Review struct {
	Id_Review   int
	Username    string
	Description string
	Estimation  float64
	Date        string
	Photos      []Photo
}

func ConstructReview(username string, description string, estimation float64, date string) *Review {
	return &Review{
		Username:    username,
		Description: description,
		Estimation:  estimation,
		Date:        date,
	}
}

func TakeReviews(id_route int, pool *pgxpool.Pool) ([]Review, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
		SELECT id_review, username, description, estimation, date_review
		FROM reviews
		INNER JOIN users ON reviews.id_user = users.id_user
		WHERE id_route = $1`, id_route)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Review

	for rows.Next() {
		var review Review
		var date time.Time
		if err := rows.Scan(
			&review.Id_Review,
			&review.Username,
			&review.Description,
			&review.Estimation,
			&date,
		); err != nil {
			return nil, fmt.Errorf("Error scan Reviews: %s", err.Error())
		}
		photos, err := TakePhotos(id_route, pool)
		if err != nil {
			return nil, fmt.Errorf("Error getting Photos: %s", err.Error())
		}
		review.Photos = photos
		review.Date = date.Format("02.01.2006")

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func TakeReviewById(id_review int, pool *pgxpool.Pool) (*Review, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	var review Review

	err = conn.QueryRow(context.Background(), `
		SELECT id_review, username, description, estimation, date_review
		FROM reviews
		INNER JOIN users ON reviews.id_user = users.id_user
		WHERE id_route = $1`, id_review).Scan(&review.Id_Review, &review.Username, &review.Description, &review.Estimation, &review.Date)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (rev *Review) CreateReview(id_route int, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
		INSERT INTO reviews 
		(username, description, estimation, date_review, id_route) 
		VALUES ($1, $2, $3, $4, $5);`, rev.Username, rev.Description, rev.Estimation, rev.Date, id_route)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
