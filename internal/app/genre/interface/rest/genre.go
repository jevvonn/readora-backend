package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/readora-backend/internal/app/genre/usecase"
	"github.com/jevvonn/readora-backend/internal/constant"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/middleware"
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
	router.Post("/genres", middleware.Authenticated, middleware.RequireRoles(constant.RoleAdmin), handler.CreateGenre)
}

// @Summary      Get Genres
// @Description  Get Genres
// @Tags         Genres
// @Accept       json
// @Produce      json
// @Success      200  object   models.JSONResponseModel{data=[]string,errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Router       /api/genres [get]
func (h *GenreHandler) GetAllGenres(ctx *fiber.Ctx) error {
	genres, err := h.genreUsecase.GetAllGenres()
	if err != nil {
		return err
	}

	return ctx.JSON(models.JSONResponseModel{
		Message: "Success get all genres",
		Data:    genres,
	})
}

// @Summary      Create Genres
// @Description  Create Genres (Require Admin Account)
// @Tags         Genres
// @Accept       json
// @Produce      json
// @Param 		 req body dto.CreateGenreRequest true "Create Genre Request"
// @Success      200  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/genres [post]
func (h *GenreHandler) CreateGenre(ctx *fiber.Ctx) error {
	var req dto.CreateGenreRequest

	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	if err := h.genreUsecase.CreateGenre(req); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(models.JSONResponseModel{
		Message: "Success create genre",
	})
}
