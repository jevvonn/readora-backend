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
	router.Delete("/comments/:commentId/replies/:replyId", middleware.Authenticated, handler.DeleteReply)
}

// @Summary      Get Reply of a Comment
// @Description  Get Reply of a Comment
// @Tags         Replies
// @Accept       json
// @Produce      json
// @Param        commentId path string true "Comment ID"
// @Param        req body dto.CreateReplyRequest true "Request Reply"
// @Success      200  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/comments/{commentId}/replies [post]
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

// @Summary      Get Replies
// @Description  Get Replies
// @Tags         Replies
// @Accept       json
// @Produce      json
// @Param        commentId path string true "Comment ID"
// @Param        limit query int false "Limit" default:10
// @Param        page query int false "Page" default:1
// @Param        sort_order query string false "Sort Order"
// @Param        parent_id query string false "Reply Parent ID"
// @Success      200  object   models.JSONResponseModel{data=[]dto.GetRepliesResponse{},errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/comments/{commentId}/replies [get]
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

// @Summary      Delete Reply
// @Description  Delete Reply
// @Tags         Replies
// @Accept       json
// @Produce      json
// @Param        commentId path string true "Comment ID"
// @Param        replyId path string true "Reply ID"
// @Success      200  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/comments/{commentId}/replies/{replyId} [delete]
func (h *ReplyHandler) DeleteReply(ctx *fiber.Ctx) error {
	err := h.replyUsecase.DeleteReply(ctx)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(models.JSONResponseModel{
		Message: "Reply deleted successfully",
	})
}
