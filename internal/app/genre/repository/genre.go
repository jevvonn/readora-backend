package repository

import (
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type GenreRepositoryItf interface {
	// GetGenres returns all genres
	GetAllGenres() ([]string, error)
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
