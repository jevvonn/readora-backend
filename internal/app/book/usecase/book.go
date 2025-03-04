package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jevvonn/readora-backend/helper"
	"github.com/jevvonn/readora-backend/internal/app/book/repository"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/storage"
)

type BookUsecaseItf interface {
	CreateBook(ctx *fiber.Ctx, req dto.CreateBookRequest) error
}

type BookUsecase struct {
	bookRepo repository.BookPostgreSQLItf
	storage  storage.StorageItf
	log      logger.LoggerItf
}

func NewBookUsecase(userRepo repository.BookPostgreSQLItf, storage storage.StorageItf, log logger.LoggerItf) BookUsecaseItf {
	return &BookUsecase{userRepo, storage, log}
}

func (u *BookUsecase) CreateBook(ctx *fiber.Ctx, req dto.CreateBookRequest) error {
	log := "[BookUsecase][CreateBook]"

	publishDate, err := helper.StringISOToDateTime(req.PublishDate)
	if err != nil {
		errValidation := errorpkg.ErrValidationTimeFormat("publish_date")
		u.log.Error(log, errValidation)

		return errValidation
	}

	pdfFile, _ := ctx.FormFile("pdf_file")
	if pdfFile == nil {
		errValidation := errorpkg.ErrValidationFileRequired("pdf_file")
		u.log.Error(log, errValidation)

		return errValidation
	}

	fileType, err := helper.GetFileMimeType(pdfFile)
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}
	if fileType != "application/pdf" {
		errValidation := errorpkg.ErrValidationFileMimeType("pdf_file", []string{".pdf"})
		u.log.Error(log, errValidation)
		return errValidation
	}

	// Upload PDF File
	uniqueId := uuid.New().String()
	fileName := uniqueId + ".pdf"
	fileKey := "books/" + fileName

	pdfURL, err := u.storage.UploadFile(pdfFile, "books", fileName, "application/pdf")
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error() + " - failed to upload file")
	}

	userId := ctx.Locals("userId").(string)
	book := entity.Book{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,
		PublishDate: publishDate,
		FileKey:     fileKey,
		FileURL:     pdfURL,
		OwnerID:     uuid.MustParse(userId),
		// COVER IMAGE NOT DONE YET
		CoverImageKey: "-",
		CoverImageURL: "-",
	}

	err = u.bookRepo.Create(book)
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	return nil
}
