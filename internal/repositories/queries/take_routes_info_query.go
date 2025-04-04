package queries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoutesInfoResults struct {
	IdRoute           int
	IdUser            int // if we request routes by user
	Name_route        string
	Place_route       string
	Description_route string
	Yandex_route      string
	Estimation        float64
	// Reviews            []Reviews
	// Photos             []Photo
	PathToPhotoPreview string
}

func TakeRoutesInfoQueryPopular(pool *pgxpool.Pool) ([]RoutesInfoResults, error) {
	var routsInfo []RoutesInfoResults
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `SELECT id_route, route_name,
													 route_place, yandex_route, 
													 route_description, 
													 estimation, 
													 path_to_photo_preview 
													 FROM routes
													 ORDER BY estimation DESC`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id_route int
		var name_route string
		var place_route string
		var yandex_route string
		var route_description string
		var estimation float64
		var path_to_photo_preview sql.NullString

		if err := rows.Scan(&id_route, &name_route, &place_route, &yandex_route, &route_description, &estimation, &path_to_photo_preview); err != nil {
			return nil, fmt.Errorf("Error scan route info: %s", err.Error())
		}
		var pathPreview string
		if path_to_photo_preview.Valid {
			pathPreview = path_to_photo_preview.String
		} else {
			pathPreview = "" // Или любое другое значение по умолчанию
		}

		routsInfo = append(routsInfo, RoutesInfoResults{
			IdRoute:           id_route,
			Name_route:        name_route,
			Place_route:       place_route,
			Yandex_route:      yandex_route,
			Description_route: route_description,
			Estimation:        estimation,
			// Reviews:            reviews,
			// Photos:             photos,
			PathToPhotoPreview: pathPreview,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return routsInfo, nil
}

func TakeUsersRoutesInfoQuery(pool *pgxpool.Pool, userID int) ([]RoutesInfoResults, error) {
	var routsInfo []RoutesInfoResults
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `SELECT id_route, 
													route_name, route_place, yandex_route, 
													route_description, 
													estimation, path_to_photo_preview 
													from routes 
													WHERE id_user = $1 
													ORDER BY estimation DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("Error query: %s", err.Error())
	}
	for rows.Next() {
		var id_route int
		var name_route string
		var place_route string
		var yandex_route string
		var route_description string
		var estimation float64
		var path_to_photo_preview sql.NullString

		if err := rows.Scan(&id_route, &name_route, &place_route, &yandex_route, &route_description, &estimation, &path_to_photo_preview); err != nil {
			return nil, fmt.Errorf("Error scan route info: %s", err.Error())
		}
		var pathPreview string
		if path_to_photo_preview.Valid {
			pathPreview = path_to_photo_preview.String
		} else {
			pathPreview = "" // Или любое другое значение по умолчанию
		}

		routsInfo = append(routsInfo, RoutesInfoResults{
			IdRoute:           id_route,
			Name_route:        name_route,
			Place_route:       place_route,
			Yandex_route:      yandex_route,
			Description_route: route_description,
			Estimation:        estimation,
			// Reviews:            reviews,
			// Photos:             photos,
			PathToPhotoPreview: pathPreview,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return routsInfo, nil
}
