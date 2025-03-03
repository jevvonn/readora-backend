package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/app/book/usecase"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/validator"
	"github.com/jevvonn/readora-backend/internal/models"
)

type BookHandler struct {
	router    fiber.Router
	authBook  usecase.BookUsecaseItf
	validator validator.ValidationService
	log       logger.LoggerItf
	response  models.ResponseItf
}

func NewBookHandler() {
	handler := BookHandler{}
}
