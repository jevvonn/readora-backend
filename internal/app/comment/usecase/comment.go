package usecase

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	bookRepo "github.com/jevvonn/readora-backend/internal/app/book/repository"
	"github.com/jevvonn/readora-backend/internal/app/comment/repository"
	"github.com/jevvonn/readora-backend/internal/constant"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type CommentUsecaseItf interface {
	CreateComment(ctx *fiber.Ctx, req dto.CreateCommentRequest) error
	DeleteComment(ctx *fiber.Ctx) error
}

type CommentUsecase struct {
	commentRepo repository.CommentPostgreSQLItf
	bookRepo    bookRepo.BookPostgreSQLItf
	log         logger.LoggerItf
}

func NewCommentUsecase(
	commentRepo repository.CommentPostgreSQLItf,
	bookRepo bookRepo.BookPostgreSQLItf,
	log logger.LoggerItf,
) CommentUsecaseItf {
	return &CommentUsecase{commentRepo, bookRepo, log}
}

func (u *CommentUsecase) CreateComment(ctx *fiber.Ctx, req dto.CreateCommentRequest) error {
	log := "[CommentUsecase][CreateComment]"
	bookId := ctx.Params("bookId")
	userId := ctx.Locals("userId").(string)

	book, err := u.bookRepo.GetSpecificBook(bookId)
	if err != nil {
		u.log.Error(log, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorpkg.ErrNotFoundResource.WithCustomMessage("Book not found")
		}
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	if !book.IsPublic {
		return errorpkg.ErrForbiddenResource.WithCustomMessage("Book is private")
	}

	_, err = u.commentRepo.GetSpecificComment(entity.Comment{
		BookId: book.ID,
		UserId: uuid.MustParse(userId),
	})

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			u.log.Error(log, err)
			return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
		}
	} else {
		return errorpkg.ErrBadRequest.WithCustomMessage("You have already commented and reviews on this book")
	}

	comment := entity.Comment{
		ID:      uuid.New(),
		BookId:  book.ID,
		UserId:  uuid.MustParse(userId),
		Content: req.Content,
		Rating:  req.Rating,
	}

	err = u.commentRepo.CreateComment(comment)
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	return nil
}

func (u *CommentUsecase) DeleteComment(ctx *fiber.Ctx) error {
	log := "[CommentUsecase][GetComments]"

	bookId := ctx.Params("bookId")
	commentId := ctx.Params("commentId")
	userId := ctx.Locals("userId").(string)
	role := ctx.Locals("role").(string)

	commentUUID, err := uuid.Parse(commentId)
	if err != nil {
		return errorpkg.ErrBadRequest.WithCustomMessage("Invalid comment id")
	}

	book, err := u.bookRepo.GetSpecificBook(bookId)
	if err != nil {
		u.log.Error(log, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorpkg.ErrNotFoundResource.WithCustomMessage("Book not found")
		}
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	if !book.IsPublic {
		return errorpkg.ErrForbiddenResource.WithCustomMessage("Book is private")
	}

	comments, err := u.commentRepo.GetSpecificComment(entity.Comment{
		ID: commentUUID,
	})

	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	if comments.UserId.String() != userId && role != constant.RoleAdmin {
		return errorpkg.ErrForbiddenResource.WithCustomMessage("You are not allowed to delete this comment")
	}

	err = u.commentRepo.DeleteComment(entity.Comment{
		ID: commentUUID,
	})
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	return nil
}
