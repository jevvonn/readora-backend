package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/app/comment/usecase"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/validator"
	"github.com/jevvonn/readora-backend/internal/middleware"
	"github.com/jevvonn/readora-backend/internal/models"
)

type CommentHandler struct {
	commentUsecase usecase.CommentUsecaseItf
	validator      validator.ValidationService
}

func NewCommentHandler(router fiber.Router, commentUsecase usecase.CommentUsecaseItf, validator validator.ValidationService) {
	handler := &CommentHandler{
		commentUsecase, validator,
	}

	router.Post("/books/:bookId/comments", middleware.Authenticated, handler.CreateComment)
	router.Get("/books/:bookId/comments", middleware.Authenticated, handler.GetComments)
	router.Get("/comments", middleware.Authenticated, handler.GetComments)
	router.Delete("/books/:bookId/comments/:commentId", middleware.Authenticated, handler.DeleteComment)
}

func (h *CommentHandler) CreateComment(ctx *fiber.Ctx) error {
	var req dto.CreateCommentRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return errorpkg.ErrBadRequest.WithCustomMessage(err.Error())
	}

	err = h.validator.Validate(req)
	if err != nil {
		return err
	}

	err = h.commentUsecase.CreateComment(ctx, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(models.JSONResponseModel{
		Message: "Comment created successfully",
	})
}

func (h *CommentHandler) GetComments(ctx *fiber.Ctx) error {
	var query dto.GetCommentsQuery
	err := ctx.QueryParser(&query)
	if err != nil {
		return errorpkg.ErrBadRequest.WithCustomMessage(err.Error())
	}

	comments, page, limit, err := h.commentUsecase.GetComments(ctx, query)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(models.JSONResponseModel{
		Data: map[string]any{
			"comments": comments,
			"page":     page,
			"limit":    limit,
			"total":    len(comments),
		},
	})
}

func (h *CommentHandler) DeleteComment(ctx *fiber.Ctx) error {
	err := h.commentUsecase.DeleteComment(ctx)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(models.JSONResponseModel{
		Message: "Comment deleted",
	})
}
