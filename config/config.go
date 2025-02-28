package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string `env:"APP_ENV"`
	AppPort string `env:"APP_PORT"`

	JWTSecret string `env:"JWT_SECRET"`

	DbHost     string `env:"DB_HOST"`
	DbPort     string `env:"DB_PORT"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbName     string `env:"DB_NAME"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     string `env:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`

	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     int    `env:"SMTP_PORT"`
	SMTPUsername string `env:"SMTP_USERNAME"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
	SMTPEmail    string `env:"SMTP_EMAIL"`
}

var cfg Config

func Load() Config {
	if cfg.AppEnv == "" {
		New()
	}

	return cfg
}

func New() Config {
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
