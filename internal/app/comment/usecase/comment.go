package usecase

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	bookRepo "github.com/jevvonn/readora-backend/internal/app/book/repository"
	"github.com/jevvonn/readora-backend/internal/app/comment/repository"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type CommentUsecaseItf interface {
	CreateComment(ctx *fiber.Ctx, req dto.CreateCommentRequest) error
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
		return errorpkg.ErrBadRequest.WithCustomMessage("You have already commented on this book")
	}

	comment := entity.Comment{
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
