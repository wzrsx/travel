package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigEmail struct {
	Sender   string
	Password string
	SmtpHost string
	SmtpPort string
}

func NewConfigEmail() (*ConfigEmail, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file")
		return nil, err
	}
	return &ConfigEmail{
		Sender:   os.Getenv("GMAIL"),
		Password: os.Getenv("GMAIL_PASSWD"),
		SmtpHost: os.Getenv("SMTP_HOST"),
		SmtpPort: os.Getenv("SMTP_PORT"),
	}, nil
}
