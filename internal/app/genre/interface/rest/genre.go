package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/app/genre/usecase"
	"github.com/jevvonn/readora-backend/internal/models"
)

type GenreHandler struct {
	genreUsecase usecase.GenreUsecaseItf
}

func NewGenreHandler(router fiber.Router, genreUsecase usecase.GenreUsecaseItf) {
	handler := GenreHandler{
		genreUsecase: genreUsecase,
	}

	router.Get("/genres", handler.GetAllGenres)
}

func (h *GenreHandler) GetAllGenres(c *fiber.Ctx) error {
	genres, err := h.genreUsecase.GetAllGenres()
	if err != nil {
		return err
	}

	return c.JSON(models.JSONResponseModel{
		Message: "Success get all genres",
		Data:    genres,
	})
}
