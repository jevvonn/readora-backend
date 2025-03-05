package repository

import (
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type CommentPostgreSQLItf interface {
	CreateComment(req entity.Comment) error
	GetSpecificComment(req entity.Comment) (entity.Comment, error)
}

type CommentPostgreSQL struct {
	db     *gorm.DB
	logger logger.LoggerItf
}

func NewCommentPostgreSQL(db *gorm.DB, logger logger.LoggerItf) CommentPostgreSQLItf {
	return &CommentPostgreSQL{db, logger}
}

func (r *CommentPostgreSQL) CreateComment(req entity.Comment) error {
	log := "[CommentPostgreSQL][CreateComment]"

	err := r.db.Create(&req).Error
	if err != nil {
		r.logger.Error(log, err)
		return err
	}

	return nil
}

func (r *CommentPostgreSQL) GetSpecificComment(req entity.Comment) (entity.Comment, error) {
	log := "[CommentPostgreSQL][GetSpecificComment]"
	var comment entity.Comment
	err := r.db.Where(&req).First(&comment).Error

	if err != nil {
		r.logger.Error(log, err)
		return entity.Comment{}, err
	}

	return comment, nil
}
