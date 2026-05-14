package queries

import (
	"context"
	"fmt"
	"log"
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

// CheckReviewExists — проверяет, оставлял ли пользователь отзыв на этот маршрут
// Возвращает true, если отзыв уже есть (для предотвращения дубликатов)
func CheckReviewExists(routeID, userID int, pool *pgxpool.Pool) (bool, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return false, fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	var exists bool
	err = conn.QueryRow(context.Background(), `
		SELECT EXISTS(
			SELECT 1 FROM reviews 
			WHERE id_route = $1 AND id_user = $2
		)`, routeID, userID).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check review exists: %w", err)
	}
	return exists, nil
}

func UpdateRouteEstimation(routeID int, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	result, err := conn.Exec(context.Background(), `
		UPDATE routes 
		SET estimation = COALESCE((
			SELECT ROUND(AVG(estimation)::numeric, 1)
			FROM reviews
			WHERE id_route = $1
		), 0)
		WHERE id_route = $1`,
		routeID)

	if err != nil {
		return fmt.Errorf("update route estimation: %w", err)
	}

	// Логируем, сколько строк было обновлено
	rowsAffected := result.RowsAffected()
	log.Printf("📊 Route #%d: estimation updated, rows affected: %d", routeID, rowsAffected)

	if rowsAffected == 0 {
		log.Printf("⚠️ Warning: no rows updated for route #%d", routeID)
	}

	return nil
}

// GetRouteAverageEstimation — вспомогательная функция: получает текущий средний рейтинг маршрута
// Может пригодиться для отладки или дополнительной логики
func GetRouteAverageEstimation(routeID int, pool *pgxpool.Pool) (float64, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return 0, fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	var avgEstimation float64
	err = conn.QueryRow(context.Background(), `
		SELECT COALESCE(ROUND(AVG(estimation)::numeric, 1), 0)
		FROM reviews 
		WHERE id_route = $1`,
		routeID).Scan(&avgEstimation)

	if err != nil {
		return 0, fmt.Errorf("get average estimation: %w", err)
	}
	return avgEstimation, nil
}
