package repository

import (
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type GenreRepositoryItf interface {
	GetAllGenres() ([]string, error)
	CreateGenre(genre entity.Genre) error
	IsGenreExist(genreName string) (bool, error)
}

type GenreRepository struct {
	db  *gorm.DB
	log logger.LoggerItf
}

func NewGenreRepository(db *gorm.DB, log logger.LoggerItf) *GenreRepository {
	return &GenreRepository{db, log}
}

func (r *GenreRepository) GetAllGenres() ([]string, error) {
	log := "[GenreRepository][GetAllGenres]"

	var genres []string
	if err := r.db.Model(&entity.Genre{}).Pluck("name", &genres).Error; err != nil {
		r.log.Error(log, err)
		return nil, err
	}

	return genres, nil
}

func (r *GenreRepository) IsGenreExist(genreName string) (bool, error) {
	log := "[GenreRepository][IsGenreExist]"

	var genre entity.Genre
	if err := r.db.Where("name = ?", genreName).First(&genre).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		r.log.Error(log, err)
		return false, err
	}

	return true, nil
}

func (r *GenreRepository) CreateGenre(genre entity.Genre) error {
	log := "[GenreRepository][CreateGenre]"

	if err := r.db.Create(&genre).Error; err != nil {
		r.log.Error(log, err)
		return err
	}

	return nil
}
