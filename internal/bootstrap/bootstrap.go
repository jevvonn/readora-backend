package bootstrap

import (
	"fmt"
	"time"

	"github.com/gofiber/swagger"
	docs "github.com/jevvonn/readora-backend/docs"

	"github.com/gofiber/fiber/v2"
	cors "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jevvonn/readora-backend/config"
	authHandler "github.com/jevvonn/readora-backend/internal/app/auth/interface/rest"
	authRepository "github.com/jevvonn/readora-backend/internal/app/auth/repository"
	authUsecase "github.com/jevvonn/readora-backend/internal/app/auth/usecase"
	bookHandler "github.com/jevvonn/readora-backend/internal/app/book/interface/rest"
	bookRepository "github.com/jevvonn/readora-backend/internal/app/book/repository"
	bookUsecase "github.com/jevvonn/readora-backend/internal/app/book/usecase"
	userRepository "github.com/jevvonn/readora-backend/internal/app/user/repository"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/mailer"
	"github.com/jevvonn/readora-backend/internal/infra/postgresql"
	"github.com/jevvonn/readora-backend/internal/infra/redis"
	"github.com/jevvonn/readora-backend/internal/infra/storage"
	"github.com/jevvonn/readora-backend/internal/infra/validator"
)

const idleTimeout = 5 * time.Second

func Start() error {
	// Load .env file
	conf := config.New()

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		IdleTimeout:  idleTimeout,
		ErrorHandler: FiberErrorHandler,
	})

	// Logger
	logger := logger.New()

	// Validator
	vd := validator.NewValidator()

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

	// Connect to Redis
	rdb := redis.New()

	// Connect to Mailer
	mailer := mailer.New()

	if err != nil {
		return err
	}

	// Supabase Storage
	storage := storage.New()

	// Command flag check (migration, seeder)
	CommandHandler(db)

	// Routes Group
	apiRouter := app.Group("/api")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})

	// Cors Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Repo Instance
	authRepo := authRepository.NewAuthRepository(rdb, logger)
	userRepo := userRepository.NewUserPostgreSQL(db, logger)
	bookRepo := bookRepository.NewBookPostgreSQL(db, logger)

	// Usecase Instance
	authUsecase := authUsecase.NewAuthUsecase(userRepo, authRepo, logger, mailer)
	bookUsecase := bookUsecase.NewBookUsecase(bookRepo, storage, logger)

	// Handler Instance
	authHandler.NewAuthHandler(apiRouter, authUsecase, vd)
	bookHandler.NewBookHandler(apiRouter, bookUsecase, vd)

	// Swagger Docs
	httpProtocol := "http"
	if conf.AppEnv == "production" {
		httpProtocol = "https"
	}

	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = conf.AppBaseURL
	docs.SwaggerInfo.Title = "Readora Backend Service Documentation"
	swaggerHandler := swagger.New(swagger.Config{
		URL: fmt.Sprintf("%s://%s/docs/doc.json", httpProtocol, conf.AppBaseURL),
	})

	app.Get("/docs/*", swaggerHandler)

	// Start the server
	listenAddr := fmt.Sprintf("127.0.0.1:%s", conf.AppPort)
	if conf.AppEnv == "production" {
		listenAddr = fmt.Sprintf(":%s", conf.AppPort)
	}

	return app.Listen(listenAddr)
}
