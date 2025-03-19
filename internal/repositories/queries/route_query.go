package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Route struct {
	Yandex_route      string
	Route_name        string
	Route_place       string
	Route_description string
	Route_estimation  int
	UserID            int
}

func ConstructRoute(yandex_route string, route_name string, route_place string, route_description string, userID int) *Route {
	return &Route{
		Yandex_route:      yandex_route,
		Route_name:        route_name,
		Route_place:       route_place,
		Route_description: route_description,
		UserID:            userID,
	}
}

func (r *Route) CreateNewRoute(pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `
	INSERT INTO routes (route_name, yandex_route, 
	estimation, route_place, route_description, id_user)
	VALUES ($1, $2, $3, $4, $5, $6)`, r.Route_name, r.Yandex_route, 0, r.Route_place, r.Route_description, r.UserID)
	if err != nil {
		return err
	}
	return nil
}

func TakeRouteInfo(id_route int, pool *pgxpool.Pool) (*Route, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(), `
		SELECT route_name, route_place, 
		route_description, yandex_route, estimation, id_user 
		FROM routes 
		WHERE id_route = $1`, id_route)

	var route Route
	err = row.Scan(&route.Route_name, &route.Route_place, &route.Route_description, &route.Yandex_route, &route.Route_estimation, &route.UserID)
	if err != nil {
		return nil, err
	}

	return &route, nil
}
