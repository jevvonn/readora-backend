package usecase

import (
	"github.com/jevvonn/readora-backend/internal/app/genre/repository"
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
)

type GenreUsecaseItf interface {
	GetAllGenres() ([]string, error)
	CreateGenre(req dto.CreateGenreRequest) error
}

type GenreUsecase struct {
	genreRepo repository.GenreRepositoryItf
	log       logger.LoggerItf
}

func NewGenreUsecase(genreRepo repository.GenreRepositoryItf) *GenreUsecase {
	return &GenreUsecase{genreRepo: genreRepo}
}

func (u *GenreUsecase) GetAllGenres() ([]string, error) {
	log := "[GenreUsecase][GetAllGenres]"

	genres, err := u.genreRepo.GetAllGenres()
	if err != nil {
		u.log.Error(log, err)
		return nil, errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	return genres, nil
}

func (u *GenreUsecase) CreateGenre(req dto.CreateGenreRequest) error {
	log := "[GenreUsecase][CreateGenre]"

	isGenreExist, err := u.genreRepo.IsGenreExist(req.Name)
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	if isGenreExist {
		return errorpkg.ErrBadRequest.WithCustomMessage("Genre already exist")
	}

	err = u.genreRepo.CreateGenre(entity.Genre{
		Name: req.Name,
	})
	if err != nil {
		u.log.Error(log, err)
		return errorpkg.ErrInternalServerError.WithCustomMessage(err.Error())
	}

	return nil
}
