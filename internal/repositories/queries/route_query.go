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
	var r Route
	return &r
}

func (r *Route) CreateNewRoute(pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(), "INSERT INTO routes (route_name, yandex_route, estimation, route_description, id_user) VALUES ($1, $2, $3, $4, $5)", r.Yandex_route, r.Route_name, r.Route_place, r.Route_description, r.UserID).Scan()
	if err != nil {
		return err
	}
	return nil
}
