package infra

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/reodora-backend/config"
	"github.com/jevvonn/reodora-backend/internal/infra/postgresql"
)

const idleTimeout = 5 * time.Second

func Bootstrap() error {
	// Load .env file
	conf := config.New()

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
	})

	// Connect to PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		conf.DbHost,
		conf.DbPort,
		conf.DbUser,
		conf.DbPassword,
		conf.DbName,
	)
	_, err := postgresql.New(dsn)

	if err != nil {
		return err
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})

	return app.Listen("localhost:8080")
}
