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
	UserID            int
}

func CreateRouteStruct(yandex_route string, route_name string, route_place string, route_description string, userID int) *Route {
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
	INSERT INTO routes (route_name, yandex_route, estimation, route_place, route_description, id_user)
	VALUES ($1, $2, $3, $4, $5, $6)`, r.Route_name, r.Yandex_route, 0, r.Route_place, r.Route_description, r.UserID)
	if err != nil {
		return err
	}
	return nil
}
