package infra

import (
	"flag"
	"fmt"
	"os"
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
	db, err := postgresql.New(dsn)

	if err != nil {
		return err
	}

	// Migration flag check
	var migrationCmd string
	var seederCmd bool

	flag.StringVar(&migrationCmd, "m", "", "Migrate database 'up' or 'down'")
	flag.BoolVar(&seederCmd, "s", true, "Seed database")
	flag.Parse()

	if migrationCmd != "" {
		postgresql.Migrate(db, migrationCmd)
		os.Exit(0)
	}

	if seederCmd {
		postgresql.Seed(db)
		os.Exit(0)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})

	// Start the server
	listenAddr := fmt.Sprintf("localhost:%s", conf.AppPort)
	if conf.AppEnv == "production" {
		listenAddr = fmt.Sprintf(":%s", conf.AppPort)
	}

	return app.Listen(listenAddr)
}
