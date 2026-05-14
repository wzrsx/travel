package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Review — структура отзыва (оставляем как есть)
type Review struct {
	Id_Review   int
	Username    string
	Description string
	Estimation  float64
	Date        string
	// Photos []Photo — ❌ удалите это поле, отзывы не содержат фото маршрута
}

func ConstructReview(username, description string, estimation float64, date string) *Review {
	return &Review{
		Username:    username,
		Description: description,
		Estimation:  estimation,
		Date:        date,
	}
}

// TakeReviews — получает отзывы маршрута
func TakeReviews(routeID int, pool *pgxpool.Pool) ([]Review, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
		SELECT r.id_review, u.username, r.description, r.estimation, r.date_review
		FROM reviews r
		INNER JOIN users u ON r.id_user = u.id_user
		WHERE r.id_route = $1
		ORDER BY r.date_review DESC`,
		routeID)
	if err != nil {
		return nil, fmt.Errorf("query reviews: %w", err)
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var rev Review
		var dateRaw time.Time
		// Scan в порядке полей вашей структуры:
		// Id_Review (с подчёркиванием!), Username, Description, Estimation, Date
		if err := rows.Scan(&rev.Id_Review, &rev.Username, &rev.Description, &rev.Estimation, &dateRaw); err != nil {
			return nil, fmt.Errorf("scan review: %w", err)
		}
		rev.Date = dateRaw.Format("02.01.2006")
		// ❌ Не добавляем rev.Photos = ... — отзывы не содержат фото маршрута
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}

// TakeReviewById — получает отзыв по ID (исправлено: фильтр по id_review)
func TakeReviewById(reviewID int, pool *pgxpool.Pool) (*Review, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	var rev Review
	var dateRaw time.Time
	err = conn.QueryRow(context.Background(), `
		SELECT r.id_review, u.username, r.description, r.estimation, r.date_review
		FROM reviews r
		INNER JOIN users u ON r.id_user = u.id_user
		WHERE r.id_review = $1`,
		reviewID).Scan(&rev.Id_Review, &rev.Username, &rev.Description, &rev.Estimation, &dateRaw)

	if err != nil {
		return nil, fmt.Errorf("scan review: %w", err)
	}
	rev.Date = dateRaw.Format("02.01.2006")
	return &rev, nil
}

// CreateReview — создаёт отзыв (исправлено: Exec вместо Query)
func (rev *Review) CreateReview(routeID, userID int, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	// ❌ Было: conn.Query(...) — неправильно для INSERT
	// ✅ Стало:
	_, err = conn.Exec(context.Background(), `
		INSERT INTO reviews (id_user, id_route, description, estimation, date_review)
		VALUES ($1, $2, $3, $4, $5)`,
		userID, routeID, rev.Description, rev.Estimation, time.Now())

	if err != nil {
		return fmt.Errorf("insert review: %w", err)
	}
	return nil
}
