package usecase

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/jevvonn/readora-backend/internal/app/book/repository"
	"github.com/jevvonn/readora-backend/internal/constant"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type AIFeatureUsecaseItf interface {
	HighlightText(ctx *fiber.Ctx, req dto.HighlightTextRequest) (dto.HighlightTextResponse, error)
}

type AIFeatureUsecase struct {
	geminiModel *genai.GenerativeModel
	bookRepo    repository.BookPostgreSQLItf
	log         logger.LoggerItf
}

func NewAIFeatureUsecase(geminiModel *genai.GenerativeModel, bookRepo repository.BookPostgreSQLItf, log logger.LoggerItf) AIFeatureUsecaseItf {
	return &AIFeatureUsecase{geminiModel, bookRepo, log}
}

func (u *AIFeatureUsecase) HighlightText(ctx *fiber.Ctx, req dto.HighlightTextRequest) (dto.HighlightTextResponse, error) {
	log := "[AIFeatureUsecase][HighlightText]"

	bookId := ctx.Params("bookId")
	book, err := u.bookRepo.GetSpecificBook(bookId)
	if err != nil {
		u.log.Error(log, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.HighlightTextResponse{}, errorpkg.ErrNotFoundResource.WithCustomMessage("Book not found")
		}
		return dto.HighlightTextResponse{}, errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	if book.FileAIStatus != constant.BookFileAIStatusReady {
		return dto.HighlightTextResponse{}, errorpkg.ErrBadRequest.WithCustomMessage("Book AI file feature is not ready")
	}

	txtURL := constant.GetBookTxtFile(book.ID.String())
	txtResp, err := http.Get(txtURL)
	if err != nil {
		u.log.Error(log, err)
		return dto.HighlightTextResponse{}, errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}
	defer txtResp.Body.Close()

	txtBytes, err := io.ReadAll(txtResp.Body)
	if err != nil {
		u.log.Error(log, err)
		return dto.HighlightTextResponse{}, errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	prompt := fmt.Sprintf(`
		context: Reading a book from a text file

		input : 
		"highlight-text": %s,
		"page/chapter": %s

		what's the meaning of highlight text from the dictionary and from the context of the story by around the page / chapter number?

		note: the response should be in text of a 2 paragraph merged and not more that 10 sentences in each paragraph. dont say about the page number or page chapter in the response. dont use bold or italic in the response, just plain text with no without any other backticks, and markdown in it. response language should be the same as the book language.

		if ther's no meaning, just say "No meaning found" and if the text is not found in the book, just say "Text not found in the book". If can't find the page/chapter, just say "Page/chapter not found in the book". if ask about other things, just say "I can't answer that".
	`, req.HighlightText, req.Page)

	aiReq := []genai.Part{
		genai.Blob{MIMEType: "text/plain", Data: txtBytes},
		genai.Text(prompt),
	}

	// Generate content.
	resp, err := u.geminiModel.GenerateContent(ctx.Context(), aiReq...)
	if err != nil {
		u.log.Error(log, err)
		return dto.HighlightTextResponse{}, errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	// Handle the response of generated text.
	var content genai.Content
	for _, c := range resp.Candidates {
		if c.Content != nil {
			content = *c.Content
		}
	}

	aiResponse := content.Parts[0].(genai.Text)

	return dto.HighlightTextResponse{
		AIResponse: aiResponse,
	}, nil
}
