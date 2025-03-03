package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigIP struct {
	IP   string
	Port string
}

func NewConfigIP() (*ConfigIP, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file")
		return nil, err
	}
	return &ConfigIP{
		IP:   os.Getenv("IP"),
		Port: os.Getenv("PORT"),
	}, nil
}

func (c *ConfigIP) AdressStr() string {
	return fmt.Sprintf(
		"%s:%s",
		c.IP, c.Port,
	)
}
