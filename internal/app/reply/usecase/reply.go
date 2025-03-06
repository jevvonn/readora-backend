package usecase

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	commentRepo "github.com/jevvonn/readora-backend/internal/app/comment/repository"
	replyRepo "github.com/jevvonn/readora-backend/internal/app/reply/repository"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
)

type ReplyUsecaseItf interface {
	CreateComment(ctx *fiber.Ctx, req dto.CreateReplyRequest) error
	GetRepliesByCommentId(ctx *fiber.Ctx, filter dto.GetRepliesQuery) ([]dto.GetRepliesResponse, int, int, error)
}

type ReplyUsecase struct {
	replyRepo   replyRepo.ReplyPostgreSQLItf
	commentRepo commentRepo.CommentPostgreSQLItf
	log         logger.LoggerItf
}

func NewReplyUsecase(replyRepo replyRepo.ReplyPostgreSQLItf, commentRepo commentRepo.CommentPostgreSQLItf, log logger.LoggerItf) ReplyUsecaseItf {
	return &ReplyUsecase{replyRepo, commentRepo, log}
}

func (u *ReplyUsecase) CreateComment(ctx *fiber.Ctx, req dto.CreateReplyRequest) error {
	commentParam := ctx.Params("commentId")
	userId := ctx.Locals("userId").(string)

	commentId, err := uuid.Parse(commentParam)
	if err != nil {
		return errorpkg.ErrBadRequest.WithCustomMessage("Invalid comment ID")
	}

	_, err = u.commentRepo.GetSpecificComment(entity.Comment{
		ID: commentId,
	})

	if err != nil {
		if errors.Is(err, errorpkg.ErrNotFoundResource) {
			return errorpkg.ErrNotFoundResource.WithCustomMessage("Comment not found")
		}
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	reply := entity.Reply{
		ID:        uuid.New(),
		Content:   req.Content,
		CommentId: commentId,
		UserId:    uuid.MustParse(userId),
	}

	if req.ParentId != "" {
		parentId, err := uuid.Parse(req.ParentId)
		if err != nil {
			return errorpkg.ErrBadRequest.WithCustomMessage("Invalid reply parent ID")
		}

		_, err = u.replyRepo.GetSpecificReply(entity.Reply{
			ID: parentId,
		})

		if err != nil {
			if errors.Is(err, errorpkg.ErrNotFoundResource) {
				return errorpkg.ErrNotFoundResource.WithCustomMessage("Reply not found")
			}
			return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
		}

		reply.ParentId = parentId
	}

	err = u.replyRepo.Create(reply)
	if err != nil {
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	return nil
}

func (u *ReplyUsecase) GetRepliesByCommentId(ctx *fiber.Ctx, filter dto.GetRepliesQuery) ([]dto.GetRepliesResponse, int, int, error) {
	commentParam := ctx.Params("commentId")

	if filter.Limit <= 0 {
		filter.Limit = 15
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	replies, err := u.replyRepo.GetRepliesByCommentId(commentParam, filter)
	if err != nil {
		return nil, 0, 0, errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	var res []dto.GetRepliesResponse
	for _, reply := range replies {
		rep := dto.GetRepliesResponse{
			ID:        reply.ID,
			Content:   reply.Content,
			CommentId: reply.CommentId,
			User: entity.User{
				ID:       reply.User.ID,
				Username: reply.User.Username,
			},
			CreatedAt: reply.CreatedAt,
		}

		if reply.ParentId != uuid.Nil {
			rep.ParentId = reply.ParentId
		}

		res = append(res, rep)
	}

	return res, filter.Page, filter.Limit, nil
}
