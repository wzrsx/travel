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

func TakePhotos(routeID int, pool *pgxpool.Pool) ([]Photo, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
		SELECT id_photo, path_to_photo
		FROM photos
		WHERE id_route = $1
		ORDER BY id_photo`,
		routeID)
	if err != nil {
		return nil, fmt.Errorf("query photos: %w", err)
	}
	defer rows.Close()

	var photos []Photo
	for rows.Next() {
		var p Photo
		// Scan только в существующие поля: IdPhoto, PhotoPath
		if err := rows.Scan(&p.IdPhoto, &p.PhotoPath); err != nil {
			return nil, fmt.Errorf("scan photo: %w", err)
		}
		photos = append(photos, p)
		// ❌ Не пытайтесь делать: p.IdRoute = routeID — такого поля нет!
	}
	return photos, rows.Err()
}
