package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/app/book/usecase"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/infra/validator"
	"github.com/jevvonn/readora-backend/internal/middleware"
	"github.com/jevvonn/readora-backend/internal/models"
)

type BookHandler struct {
	bookUsecase usecase.BookUsecaseItf
	validator   validator.ValidationService
}

func NewBookHandler(
	router fiber.Router,
	bookUsecase usecase.BookUsecaseItf,
	validator validator.ValidationService,
) {
	handler := &BookHandler{
		bookUsecase, validator,
	}

	router.Post("/books", middleware.Authenticated, handler.CreateBook)
}

func (h *BookHandler) CreateBook(ctx *fiber.Ctx) error {
	var req dto.CreateBookRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = h.validator.Validate(req)
	if err != nil {
		return err
	}

	err = h.bookUsecase.CreateBook(ctx, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(models.JSONResponseModel{
		Message: "Book created successfully",
	})
}
