package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/app/ai-feature/usecase"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/validator"
	"github.com/jevvonn/readora-backend/internal/middleware"
	"github.com/jevvonn/readora-backend/internal/models"
)

type AIFeatureHandler struct {
	AIFeatureUsecase usecase.AIFeatureUsecaseItf
	validator        validator.ValidationService
}

func NewAIFeatureHandler(
	router fiber.Router,
	AIFeatureUsecase usecase.AIFeatureUsecaseItf,
	validator validator.ValidationService,
) {
	handler := &AIFeatureHandler{
		AIFeatureUsecase, validator,
	}

	router.Post("/books/:bookId/highlight", middleware.Authenticated, handler.HighlightText)
}

// @Summary      Highlight Text in Book
// @Description  Highlight Text in Book Response by AI
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        bookId path string true "Book ID"
// @Param 		 req body dto.HighlightTextRequest true "Highlight Text Request"
// @Success      200  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/books/{bookId}/highlight [post]
func (h *AIFeatureHandler) HighlightText(ctx *fiber.Ctx) error {
	var req dto.HighlightTextRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return errorpkg.ErrBadRequest.WithCustomMessage(err.Error())
	}

	err = h.validator.Validate(req)
	if err != nil {
		return err
	}

	res, err := h.AIFeatureUsecase.HighlightText(ctx, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(models.JSONResponseModel{
		Message: "Successfully highlighted text",
		Data:    res,
	})
}
