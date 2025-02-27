package pool_conections

import (
	"context"
	"log"
	"travel/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool_conections struct {
	PoolConns *pgxpool.Pool
}

func Create_pool() (Pool_conections, error) {
	conf_DB := config.NewConfigDB() // Tafing data from configDB.go
	var p Pool_conections

	config, err := pgxpool.ParseConfig(conf_DB.ConnStr())
	if err != nil {
		return Pool_conections{}, err
	}
	config.MaxConns = 10 // Устанавливаем максимальное число соединений в пуле

	p.PoolConns, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return Pool_conections{}, err
	}
	defer p.PoolConns.Close()

	log.Println("pool conections sucsessfuly setted")
	return p, nil
}
