package bootstrap

import (
	"fmt"
	"net/url"
	"os"
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
	"github.com/jevvonn/readora-backend/internal/infra/postgresql"
	"github.com/jevvonn/readora-backend/internal/infra/redis"
	"github.com/jevvonn/readora-backend/internal/infra/validator"
	"github.com/jevvonn/readora-backend/internal/infra/worker"
)

const idleTimeout = 5 * time.Second

func Start() error {
	// Load .env file
	conf := config.New()

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		IdleTimeout:  idleTimeout,
		ErrorHandler: FiberErrorHandler,
		BodyLimit:    50 * 1024 * 1024 * 1024, // 50 MB
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

	// Connect Worker
	workerClient := worker.NewWorkerClient()

	if err != nil {
		return err
	}

	// Command flag check (migration, seeder)
	CommandHandler(db)

	// Check for aditional folder
	requiredFolder := []string{"tmp", "logs"}
	for _, folder := range requiredFolder {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.Mkdir(folder, 0755)
		}
	}

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
	authUsecase := authUsecase.NewAuthUsecase(userRepo, authRepo, workerClient, logger)
	bookUsecase := bookUsecase.NewBookUsecase(bookRepo, workerClient, logger)

	// Handler Instance
	authHandler.NewAuthHandler(apiRouter, authUsecase, vd)
	bookHandler.NewBookHandler(apiRouter, bookUsecase, vd)

	// Swagger Docs
	parsedURL, err := url.Parse(conf.AppBaseURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return err
	}

	host := parsedURL.Host + parsedURL.Path

	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.Title = "Readora Backend Service Documentation"
	swaggerHandler := swagger.New(swagger.Config{
		URL: conf.AppBaseURL + "/docs/doc.json",
	})

	app.Get("/docs/*", swaggerHandler)

	// Start the server
	listenAddr := fmt.Sprintf("127.0.0.1:%s", conf.AppPort)
	if conf.AppEnv == "production" {
		listenAddr = fmt.Sprintf(":%s", conf.AppPort)
	}

	return app.Listen(listenAddr)
}
