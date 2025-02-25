package bootstrap

import (
	"fmt"
	"time"

	"github.com/gofiber/swagger"
	docs "github.com/jevvonn/readora-backend/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/config"
	authHandler "github.com/jevvonn/readora-backend/internal/app/auth/interface/rest"
	authUsecase "github.com/jevvonn/readora-backend/internal/app/auth/usecase"
	userRepository "github.com/jevvonn/readora-backend/internal/app/user/repository"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/postgresql"
	"github.com/jevvonn/readora-backend/internal/infra/validator"
	"github.com/jevvonn/readora-backend/internal/models"
)

const idleTimeout = 5 * time.Second

func Start() error {
	// Load .env file
	conf := config.New()

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
	})

	// Logger
	logger := logger.New()

	// Validator
	vd := validator.NewValidator()

	// Response
	res := models.NewResponseModel()

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

	// Command flag check (migration, seeder)
	CommandHandler(db)

	// Routes Group
	apiRouter := app.Group("/api")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})

	// User Instance
	userRepo := userRepository.NewUserPostgreSQL(db, logger)

	// Auth Instance
	authUsecase := authUsecase.NewAuthUsecase(userRepo, logger)
	authHandler.NewAuthHandler(apiRouter, authUsecase, vd, logger, res)

	// Swagger Docs
	docs.SwaggerInfo.Title = "Readora Backend Service Documentation"
	app.Get("/docs/*", swagger.HandlerDefault)

	// Start the server
	listenAddr := fmt.Sprintf("localhost:%s", conf.AppPort)
	if conf.AppEnv == "production" {
		listenAddr = fmt.Sprintf(":%s", conf.AppPort)
	}

	return app.Listen(listenAddr)
}
