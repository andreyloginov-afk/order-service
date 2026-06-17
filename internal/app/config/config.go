package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/andreyloginov-afk/order-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository `split_words:"true"`
	Processor  section.Processor  `split_words:"true"`
	Monitor    section.Monitor    `split_words:"true"`
}

var Root Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	if err := envconfig.Process("App", &Root); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}
