package pool_conections

import (
	"context"
	"log"
	"travel/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool_conections struct {
	Pool *pgxpool.Pool
}

func (p *Pool_conections) Create_pool() error {
	conf_DB := config.NewConfigDB() // Tafing data from configDB.go

	config, err := pgxpool.ParseConfig(conf_DB.ConnStr())
	if err != nil {
		return err
	}
	config.MaxConns = 10 // Устанавливаем максимальное число соединений в пуле

	p.Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}
	defer p.Pool.Close()

	log.Println("pool conections sucsessfuly setted")
	return nil
}
