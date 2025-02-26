package config

import (
	"fmt"
	"os"
)

type ConfigIP struct {
	IP   string
	Port string
}

func NewConfigIP() *ConfigIP {
	return &ConfigIP{
		IP:   os.Getenv("IP"),
		Port: os.Getenv("PORT"),
	}
}

func (c *ConfigIP) AdressStr() string {
	return fmt.Sprintf(
		"%s:%s",
		c.IP, c.Port,
	)
}
