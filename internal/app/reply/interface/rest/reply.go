package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/app/reply/usecase"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/validator"
	"github.com/jevvonn/readora-backend/internal/middleware"
	"github.com/jevvonn/readora-backend/internal/models"
)

type ReplyHandler struct {
	validator    validator.ValidationService
	replyUsecase usecase.ReplyUsecaseItf
}

func NewReplyHandler(router fiber.Router, replyUsecase usecase.ReplyUsecaseItf, validator validator.ValidationService) {
	handler := &ReplyHandler{
		validator, replyUsecase,
	}

	router.Post("/comments/:commentId/replies", middleware.Authenticated, handler.CreateComment)
	router.Get("/comments/:commentId/replies", middleware.Authenticated, handler.GetReplies)
}

func (h *ReplyHandler) CreateComment(ctx *fiber.Ctx) error {
	var req dto.CreateReplyRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return errorpkg.ErrBadRequest.WithCustomMessage(err.Error())
	}

	err = h.validator.Validate(req)
	if err != nil {
		return err
	}

	err = h.replyUsecase.CreateComment(ctx, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(models.JSONResponseModel{
		Message: "Reply created successfully",
	})
}

func (h *ReplyHandler) GetReplies(ctx *fiber.Ctx) error {
	var req dto.GetRepliesQuery

	err := ctx.QueryParser(&req)
	if err != nil {
		return errorpkg.ErrBadRequest.WithCustomMessage(err.Error())
	}

	replies, page, limit, err := h.replyUsecase.GetRepliesByCommentId(ctx, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(models.JSONResponseModel{
		Message: "Replies retrieved successfully",
		Data: map[string]any{
			"replies": replies,
			"page":    page,
			"limit":   limit,
			"total":   len(replies),
		},
	})
}
