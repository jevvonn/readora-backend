package usecase

import (
	"github.com/jevvonn/readora-backend/internal/app/genre/repository"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
)

type GenreUsecaseItf interface {
	GetAllGenres() ([]string, error)
}

type GenreUsecase struct {
	genreRepo repository.GenreRepositoryItf
	log       logger.LoggerItf
}

func NewGenreUsecase(genreRepo repository.GenreRepositoryItf) *GenreUsecase {
	return &GenreUsecase{genreRepo: genreRepo}
}

// @Summary      Get Genres
// @Description  Get Genres
// @Tags         Genres
// @Accept       json
// @Produce      json
// @Success      200  object   models.JSONResponseModel{data=[]string,errors=nil}
// @Success      400  object   models.JSONResponseModel{data=nil,errors=nil}
// @Success      500  object   models.JSONResponseModel{data=nil,errors=nil}
// @Security     BearerAuth
// @Router       /api/genres [get]
func (u *GenreUsecase) GetAllGenres() ([]string, error) {
	log := "[GenreUsecase][GetAllGenres]"

	genres, err := u.genreRepo.GetAllGenres()
	if err != nil {
		u.log.Error(log, err)
		return nil, errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	return genres, nil
}
