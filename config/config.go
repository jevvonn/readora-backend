package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type config struct {
	AppEnv  string `env:"APP_ENV"`
	AppPort string `env:"APP_PORT"`

	DbHost     string `env:"DB_HOST"`
	DbPort     string `env:"DB_PORT"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbName     string `env:"DB_NAME"`
}

var cfg config

func New() config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// parse
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	return cfg
}
