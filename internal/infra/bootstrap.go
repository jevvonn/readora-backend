package infra

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/reodora-backend/config"
	authHandler "github.com/jevvonn/reodora-backend/internal/app/auth/interface/rest"
	authUsecase "github.com/jevvonn/reodora-backend/internal/app/auth/usecase"
	userRepository "github.com/jevvonn/reodora-backend/internal/app/user/repository"
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
	flag.BoolVar(&seederCmd, "s", false, "Seed database")
	flag.Parse()

	if migrationCmd != "" {
		postgresql.Migrate(db, migrationCmd)
		os.Exit(0)
	}

	if seederCmd {
		postgresql.Seed(db)
		os.Exit(0)
	}

	// Routes
	apiRouter := app.Group("/api")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})

	// User Instance
	userRepo := userRepository.NewUserPostgreSQL(db)

	// Auth Instance
	authUsecase := authUsecase.NewAuthUsecase(userRepo)
	authHandler.NewAuthHandler(apiRouter, authUsecase)

	// Start the server
	listenAddr := fmt.Sprintf("localhost:%s", conf.AppPort)
	if conf.AppEnv == "production" {
		listenAddr = fmt.Sprintf(":%s", conf.AppPort)
	}

	return app.Listen(listenAddr)
}
