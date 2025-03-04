package queries

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoutesInfoResults struct {
	IdRoute      int
	IdUser       int
	Name_route   string
	Yandex_route string
	Estimation   int8
	Reviews      []Reviews
	PathToPhoto  string
}

func TakeRoutesInfoQueryPopular(pool *pgxpool.Pool) ([]RoutesInfoResults, error) {
	var routsInfo []RoutesInfoResults
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), "SELECT id_route, route_name, yandex_route, estimation, path_to_photo from routes ORDER BY estimation DESC")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id_route int
		var name_route string
		var yandex_route string
		var estimation int
		var path_to_photo string

		if err := rows.Scan(&id_route, &name_route, &yandex_route, &estimation, &path_to_photo); err != nil {
			return nil, fmt.Errorf("Error scan route info: %s", err.Error())
		}
		reviews, err := CreateReviews(id_route, pool)
		if err != nil {
			return nil, fmt.Errorf("Error append reviews to Route: %s", err.Error())
		}
		routsInfo = append(routsInfo, RoutesInfoResults{
			IdRoute:      id_route,
			Name_route:   name_route,
			Yandex_route: yandex_route,
			Estimation:   int8(estimation),
			Reviews:      reviews,
			PathToPhoto:  path_to_photo,
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
													route_name, yandex_route, 
													estimation, path_to_photo 
													from routes 
													WHERE id_user = $1 
													ORDER BY estimation DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("Error query: %s", err.Error())
	}
	for rows.Next() {
		var id_route int
		var name_route string
		var yandex_route string
		var estimation int
		var path_to_photo string

		if err := rows.Scan(&id_route, &name_route, &yandex_route, &estimation, &path_to_photo); err != nil {
			return nil, fmt.Errorf("Error scan route info: %s", err.Error())
		}
		reviews, err := CreateReviews(id_route, pool)
		if err != nil {
			return nil, fmt.Errorf("Error append reviews to Route: %s", err.Error())
		}
		routsInfo = append(routsInfo, RoutesInfoResults{
			IdRoute:      id_route,
			Name_route:   name_route,
			Yandex_route: yandex_route,
			Estimation:   int8(estimation),
			Reviews:      reviews,
			PathToPhoto:  path_to_photo,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return routsInfo, nil
}
