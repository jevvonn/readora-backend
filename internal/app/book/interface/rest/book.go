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
	router.Delete("/books/:bookId", middleware.Authenticated, handler.DeleteBook)
	router.Post("/books", middleware.Authenticated, handler.CreateBook)
}

// @Summary      Create Book
// @Description  Create Book
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param 	  	 pdf_file formData file true "PDF File"
// @Param        title formData string true "Title"
// @Param        description formData string false "Description"
// @Param        author formData string true "Author"
// @Param        publish_date formData string true "Publish Date" example:"2021-01-02T15:04:05Z"
// @Param        genres formData string false "Genres" default:"[]" example:"[\"Fiction\",\"Non-Fiction\"]"
// @Success      200  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/books [post]
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

// @Summary      Get Books
// @Description  Get Books
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        search query string false "Search"
// @Param        genre query string false "Genre"
// @Param        limit query int false "Limit" default:10
// @Param        page query int false "Page" default:1
// @Param        sort_by query string false "Sort By"
// @Param        sort_order query string false "Sort Order"
// @Param        owner_id query string false "Owner ID"
// @Success      200  object   models.JSONResponseModel{data=[]dto.GetBooksResponse{genres=[]entity.Genre{books=nil},owner=nil},errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/books [get]
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

// @Summary      Get Specific Book
// @Description  Get Specific Book
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        bookId path string true "Book ID"
// @Success      200  object   models.JSONResponseModel{data=dto.GetBooksResponse{genres=[]entity.Genre{books=nil},owner=nil},errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/books/{bookId} [get]
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

// @Summary      Delete Book
// @Description  Delete Book
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        bookId path string true "Book ID"
// @Success      200  object   models.JSONResponseModel{data=dto.GetBooksResponse{genres=[]entity.Genre{books=nil},owner=nil},errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/books/{bookId} [delete]
func (h *BookHandler) DeleteBook(ctx *fiber.Ctx) error {
	err := h.bookUsecase.DeleteBook(ctx)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(models.JSONResponseModel{
		Message: "Book deleted successfully",
	})
}
