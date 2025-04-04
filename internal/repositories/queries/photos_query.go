package queries

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Photo struct {
	IdPhoto   int
	PhotoPath string
}

func TakePhotos(id_route int, pool *pgxpool.Pool) ([]Photo, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
		SELECT id_photo, path_to_photo
		FROM photos
		WHERE id_route = $1`, id_route)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []Photo

	for rows.Next() {
		var id_photo int
		var path_to_photo string
		if err := rows.Scan(&id_photo, &path_to_photo); err != nil {
			return nil, fmt.Errorf("Error scan Reviews: %s", err.Error())
		}
		photos = append(photos, Photo{
			IdPhoto:   id_photo,
			PhotoPath: path_to_photo,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil

}
