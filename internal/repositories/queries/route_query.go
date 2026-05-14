package queries

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Route — данные маршрута
type Route struct {
	IdRoute            int
	Yandex_route       string
	Route_name         string
	Route_place        string
	Route_description  string
	Route_estimation   float64
	PathToPhotoPreview string
	UserID             int
}

// ConstructRoute — конструктор маршрута
func ConstructRoute(yandex_route, route_name, route_place, route_description, photoPreview string, userID int) *Route {
	return &Route{
		Yandex_route:       yandex_route,
		Route_name:         route_name,
		Route_place:        route_place,
		Route_description:  route_description,
		PathToPhotoPreview: photoPreview,
		UserID:             userID,
	}
}
func GetRouteWithPhotos(routeID int, pool *pgxpool.Pool) (*Route, []Photo, error) {
	// 1. Получаем данные маршрута
	route, err := TakeRouteInfo(routeID, pool)
	if err != nil {
		return nil, nil, fmt.Errorf("get route: %w", err)
	}

	// 2. Получаем фото отдельно
	photos, err := TakePhotos(routeID, pool)
	if err != nil {
		return nil, nil, fmt.Errorf("get photos: %w", err)
	}

	return route, photos, nil
}

// ✅ GetRouteWithReviews — получает маршрут + отзывы (без фото в отзывах)
func GetRouteWithReviews(routeID int, pool *pgxpool.Pool) (*Route, []Review, error) {
	route, err := TakeRouteInfo(routeID, pool)
	if err != nil {
		return nil, nil, fmt.Errorf("get route: %w", err)
	}

	reviews, err := TakeReviews(routeID, pool)
	if err != nil {
		return nil, nil, fmt.Errorf("get reviews: %w", err)
	}

	return route, reviews, nil
}

// ✅ InsertPhoto — добавляет фото в БД (если ещё нет такой функции)
func InsertPhoto(pool *pgxpool.Pool, routeID int, photoPath string) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `
		INSERT INTO photos (id_route, path_to_photo)
		VALUES ($1, $2)`,
		routeID, photoPath)

	if err != nil {
		return fmt.Errorf("insert photo: %w", err)
	}
	return nil
}

// CreateNewRouteWithID — создаёт маршрут и возвращает его ID (использует RETURNING)
func (r *Route) CreateNewRouteWithID(pool *pgxpool.Pool) (int, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return 0, fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	var newID int
	err = conn.QueryRow(context.Background(), `
		INSERT INTO routes (
			route_name, yandex_route, estimation, 
			route_place, route_description, id_user, path_to_photo_preview
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_route`,
		r.Route_name,
		r.Yandex_route,
		0, // начальная оценка
		r.Route_place,
		r.Route_description,
		r.UserID,
		r.PathToPhotoPreview,
	).Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("insert route: %w", err)
	}
	return newID, nil
}

// TakeRouteInfo — получает данные маршрута по ID
func TakeRouteInfo(routeID int, pool *pgxpool.Pool) (*Route, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(), `
		SELECT id_route, route_name, route_place, route_description, 
		       yandex_route, estimation, path_to_photo_preview, id_user
		FROM routes
		WHERE id_route = $1`,
		routeID)

	var route Route
	err = row.Scan(&route.IdRoute, &route.Route_name, &route.Route_place,
		&route.Route_description, &route.Yandex_route, &route.Route_estimation,
		&route.PathToPhotoPreview, &route.UserID)
	if err != nil {
		return nil, fmt.Errorf("scan route: %w", err)
	}
	return &route, nil
}
