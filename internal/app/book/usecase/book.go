package usecase

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jevvonn/readora-backend/helper"
	"github.com/jevvonn/readora-backend/internal/app/book/repository"
	"github.com/jevvonn/readora-backend/internal/constant"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/storage"
)

type BookUsecaseItf interface {
	CreateBook(ctx *fiber.Ctx, req dto.CreateBookRequest) error
	GetBooks(ctx *fiber.Ctx, query dto.GetBooksQuery) ([]dto.GetBooksResponse, int, int, error)
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

	var genresReq []string
	err := json.Unmarshal([]byte(req.Genres), &genresReq)
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrValidationGenresArray
	}

	var genres []entity.Genre
	for _, genre := range genresReq {
		genres = append(genres, entity.Genre{Name: genre})
	}

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
		Genres:      genres,

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

func (u *BookUsecase) GetBooks(ctx *fiber.Ctx, query dto.GetBooksQuery) (res []dto.GetBooksResponse, page, limit int, err error) {
	log := "[BookUsecase][GetBooks]"

	if query.Limit == 0 {
		query.Limit = 10
	}

	if query.Page == 0 {
		query.Page = 1
	}

	if query.SortOrder == "" {
		query.SortOrder = "asc"
	}

	userId := ctx.Locals("userId").(string)
	filter := repository.GetBooksFilter{
		Search:    query.Search,
		Genre:     query.Genre,
		Limit:     query.Limit,
		Page:      query.Page,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
		Role:      constant.RoleAdmin,
	}

	if query.OwnerID != "" {
		if userId == query.OwnerID {
			filter.Role = constant.RoleUser
			filter.OwnerID = uuid.MustParse(userId)
		}
	}

	books, err := u.bookRepo.GetBooks(filter)
	if err != nil {
		u.log.Error(log, err)
		return nil, filter.Page, filter.Limit, errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	var booksRes []dto.GetBooksResponse
	for _, book := range books {
		booksRes = append(booksRes, dto.GetBooksResponse{
			ID:            book.ID,
			Title:         book.Title,
			Description:   book.Description,
			Author:        book.Author,
			PublishDate:   book.PublishDate,
			CoverImageKey: book.CoverImageKey,
			CoverImageURL: book.CoverImageURL,
			FileKey:       book.FileKey,
			FileURL:       book.FileURL,
			OwnerID:       book.OwnerID,
			Owner: entity.User{
				ID:       book.Owner.ID,
				Username: book.Owner.Username,
			},
			Genres: book.Genres,
		})
	}

	return booksRes, filter.Page, filter.Limit, nil
}
