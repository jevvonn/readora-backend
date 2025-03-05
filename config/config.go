package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string `env:"APP_ENV,required"`
	AppPort    string `env:"APP_PORT,required"`
	AppBaseURL string `env:"APP_BASE_URL,required"`

	FrontendBaseURL string `env:"FRONTEND_BASE_URL,required"`

	JWTSecret string `env:"JWT_SECRET,required"`

	DbHost     string `env:"DB_HOST,required"`
	DbPort     string `env:"DB_PORT,required"`
	DbUser     string `env:"DB_USER,required"`
	DbPassword string `env:"DB_PASSWORD,required"`
	DbName     string `env:"DB_NAME,required"`

	RedisHost     string `env:"REDIS_HOST,required"`
	RedisPort     string `env:"REDIS_PORT,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	RedisDB       int    `env:"REDIS_DB,required"`

	SMTPHost     string `env:"SMTP_HOST,required"`
	SMTPPort     int    `env:"SMTP_PORT,required"`
	SMTPUsername string `env:"SMTP_USERNAME,required"`
	SMTPPassword string `env:"SMTP_PASSWORD,required"`
	SMTPEmail    string `env:"SMTP_EMAIL,required"`

	SupabaseProjectURL   string `env:"SUPABASE_PROJECT_URL,required"`
	SupabaseProjectToken string `env:"SUPABASE_PROJECT_TOKEN,required"`
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
