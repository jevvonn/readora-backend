package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/app/book/usecase"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
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

	router.Get("/books", middleware.Authenticated, handler.GetBooks)
	router.Get("/books/:bookId", middleware.Authenticated, handler.GetSpecificBook)
	router.Post("/books", middleware.Authenticated, handler.CreateBook)
}

func (h *BookHandler) CreateBook(ctx *fiber.Ctx) error {
	var req dto.CreateBookRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return errorpkg.ErrBadRequest.WithCustomMessage(err.Error())
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

func (h *BookHandler) GetBooks(ctx *fiber.Ctx) error {
	var req dto.GetBooksQuery
	_ = ctx.QueryParser(&req)

	err := h.validator.Validate(req)
	if err != nil {
		return err
	}

	books, page, limit, err := h.bookUsecase.GetBooks(ctx, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(models.JSONResponseModel{
		Message: "Books fetched successfully",
		Data: map[string]any{
			"books": books,
			"page":  page,
			"limit": limit,
			"total": len(books),
		},
	})
}

func (h *BookHandler) GetSpecificBook(ctx *fiber.Ctx) error {
	book, err := h.bookUsecase.GetSpecificBook(ctx)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(models.JSONResponseModel{
		Message: "Book fetched successfully",
		Data:    book,
	})
}
