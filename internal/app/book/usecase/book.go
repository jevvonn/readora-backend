package usecase

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jevvonn/readora-backend/helper"
	"github.com/jevvonn/readora-backend/internal/app/book/repository"
	"github.com/jevvonn/readora-backend/internal/constant"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/worker"
	"gorm.io/gorm"
)

type BookUsecaseItf interface {
	CreateBook(ctx *fiber.Ctx, req dto.CreateBookRequest) error
	GetBooks(ctx *fiber.Ctx, query dto.GetBooksQuery) ([]dto.GetBooksResponse, int, int, error)
	GetReadBook(ctx *fiber.Ctx) (res dto.GetBooksResponse, err error)
	GetSpecificBook(ctx *fiber.Ctx) (res dto.GetBooksResponse, err error)
	DeleteBook(ctx *fiber.Ctx) error
}

type BookUsecase struct {
	bookRepo repository.BookPostgreSQLItf
	worker   worker.WorkerItf
	log      logger.LoggerItf
}

func NewBookUsecase(userRepo repository.BookPostgreSQLItf, worker worker.WorkerItf, log logger.LoggerItf) BookUsecaseItf {
	return &BookUsecase{userRepo, worker, log}
}

func (u *BookUsecase) CreateBook(ctx *fiber.Ctx, req dto.CreateBookRequest) error {
	log := "[BookUsecase][CreateBook]"

	userId := ctx.Locals("userId").(string)
	role := ctx.Locals("role").(string)
	// Validate Genres
	genres := []entity.Genre{}
	if req.Genres != "" && role == constant.RoleAdmin {
		genresReq := strings.Split(req.Genres, ",")

		for _, genre := range genresReq {
			genres = append(genres, entity.Genre{Name: strings.Trim(genre, " ")})
		}
	}

	// Validate Publish Date
	publishDate := time.Now()
	if req.PublishDate != "" {
		date, err := helper.StringISOToDateTime(req.PublishDate)
		if err != nil {
			errValidation := errorpkg.ErrValidationTimeFormat("publish_date")
			u.log.Error(log, errValidation)

			return errValidation
		}

		publishDate = date
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

	if fileType != "application/pdf" &&
		fileType != "application/zip" &&
		strings.Contains(pdfFile.Filename, ".epub") &&
		fileType != "application/epub+zip" &&
		strings.Contains(pdfFile.Filename, ".pdf") {
		errValidation := errorpkg.ErrValidationFileMimeType("pdf_file", []string{".pdf", ".epub"})
		u.log.Error(log, errValidation)
		return errValidation
	}

	// Upload PDF File
	extension := ".pdf"
	if strings.Contains(pdfFile.Filename, ".epub") {
		extension = ".epub"
		fileType = "application/epub+zip"
	}

	bookId := uuid.New()
	fileName := bookId.String() + extension
	fileKey := "books/" + fileName
	tempFile := "./tmp/" + fileName

	err = ctx.SaveFile(pdfFile, tempFile)
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	err = u.worker.NewBooksFileUpload(tempFile, fileName, bookId.String(), fileType)
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	book := entity.Book{
		ID:          bookId,
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,
		PublishDate: publishDate,
		FileKey:     fileKey,
		FileURL:     "-",
		OwnerID:     uuid.MustParse(userId),
		Genres:      genres,
		IsPublic:    false,
		FileType:    fileType,

		CoverImageURL: constant.GetBookDefultCoverImage(),
	}

	if role == constant.RoleAdmin {
		book.IsPublic = true
	} else {
		book.IsPublic = false
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

	sortFields := []string{"title", "author", "publish_date"}
	sortOrder := []string{"asc", "desc"}

	if query.SortBy == "" || !slices.Contains(sortFields, query.SortBy) {
		query.SortBy = ""
	}

	if query.SortOrder == "" || !slices.Contains(sortOrder, query.SortOrder) {
		query.SortOrder = "asc"
	}

	if query.Limit <= 0 {
		query.Limit = 10
	}

	if query.Page <= 0 {
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
			ID:               book.ID,
			Title:            book.Title,
			Description:      book.Description,
			Author:           book.Author,
			PublishDate:      book.PublishDate,
			CoverImageURL:    book.CoverImageURL,
			OwnerID:          book.OwnerID,
			IsPublic:         book.IsPublic,
			FileUploadStatus: book.FileUploadStatus,
			FileAIStatus:     book.FileAIStatus,
			Owner: entity.User{
				ID:       book.Owner.ID,
				Username: book.Owner.Username,
			},
			Genres: book.Genres,
		})
	}

	return booksRes, filter.Page, filter.Limit, nil
}

func (u *BookUsecase) GetSpecificBook(ctx *fiber.Ctx) (res dto.GetBooksResponse, err error) {
	log := "[BookUsecase][GetSpecificBook]"

	bookId := ctx.Params("bookId")
	book, err := u.bookRepo.GetSpecificBook(bookId)
	if err != nil {
		u.log.Error(log, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, errorpkg.ErrNotFoundResource.WithCustomMessage("Book not found")
		}
		return res, errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	userId := ctx.Locals("userId").(string)

	if !book.IsPublic {
		if userId != book.OwnerID.String() {
			return res, errorpkg.ErrForbiddenResource
		}
	}

	booksRes := dto.GetBooksResponse{
		ID:               book.ID,
		Title:            book.Title,
		Description:      book.Description,
		Author:           book.Author,
		PublishDate:      book.PublishDate,
		CoverImageURL:    book.CoverImageURL,
		OwnerID:          book.OwnerID,
		IsPublic:         book.IsPublic,
		FileUploadStatus: book.FileUploadStatus,
		FileAIStatus:     book.FileAIStatus,
		Rating:           strconv.FormatFloat(book.Rating, 'f', 1, 64),
		Owner: entity.User{
			ID:       book.Owner.ID,
			Username: book.Owner.Username,
		},
		Genres: book.Genres,
	}

	return booksRes, nil
}

func (u *BookUsecase) GetReadBook(ctx *fiber.Ctx) (res dto.GetBooksResponse, err error) {
	log := "[BookUsecase][GetReadBook]"

	bookId := ctx.Params("bookId")
	book, err := u.bookRepo.GetSpecificBook(bookId)
	if err != nil {
		u.log.Error(log, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, errorpkg.ErrNotFoundResource.WithCustomMessage("Book not found")
		}
		return res, errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	userId := ctx.Locals("userId").(string)

	if !book.IsPublic {
		if userId != book.OwnerID.String() {
			return res, errorpkg.ErrForbiddenResource
		}
	}

	booksRes := dto.GetBooksResponse{
		ID:               book.ID,
		Title:            book.Title,
		Description:      book.Description,
		Author:           book.Author,
		PublishDate:      book.PublishDate,
		CoverImageURL:    book.CoverImageURL,
		OwnerID:          book.OwnerID,
		FileKey:          book.FileKey,
		FileURL:          book.FileURL,
		FileType:         book.FileType,
		IsPublic:         book.IsPublic,
		FileUploadStatus: book.FileUploadStatus,
		FileAIStatus:     book.FileAIStatus,
	}

	return booksRes, nil
}

func (u *BookUsecase) DeleteBook(ctx *fiber.Ctx) error {
	log := "[BookUsecase][DeleteBook]"

	bookId := ctx.Params("bookId")
	userId := ctx.Locals("userId").(string)
	role := ctx.Locals("role").(string)

	book, err := u.bookRepo.GetSpecificBook(bookId)
	if err != nil {
		u.log.Error(log, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorpkg.ErrNotFoundResource.WithCustomMessage("Book not found")
		}
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	if book.OwnerID.String() != userId && role != constant.RoleAdmin {
		return errorpkg.ErrForbiddenResource
	}

	err = u.bookRepo.DeleteBook(bookId)
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	fileName := strings.Split(book.FileKey, "/")[1]
	err = u.worker.NewBooksFileDelete(fileName)
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	return nil
}
